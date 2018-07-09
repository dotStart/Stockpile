/*
 * Copyright 2018 Johannes Donath <johannesd@torchmind.com>
 * and other copyright owners as documented in the project's IP log.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package cache

import (
  "fmt"
  "strings"
  "time"

  "github.com/dotStart/Stockpile/entity"
  "github.com/google/uuid"
)

// retrieves the profile to which a given display name has been assigned at a specific time
func (c *Cache) GetProfileId(name string, at time.Time) (*entity.ProfileId, error) {
  c.logger.Debugf("processing query for profile Id associated with name \"%s\" at time %s", name, at)

  id, err := c.storage.GetProfileId(name, at)
  if err != nil {
    c.logger.Errorf("storage backend responded with error: %s", err)
    id = nil
  }
  if id == nil {
    c.logger.Debugf("cache miss - requesting update from upstream")

    c.incrementRequestCounter()
    id, err = c.upstream.GetId(name, at)
    if err != nil {
      return nil, fmt.Errorf("upstream responded with error: %s", err)
    }

    if id != nil {
      err := c.storage.PutProfileId(id)
      if err != nil {
        return nil, fmt.Errorf("storage backend responded with error: %s", err)
      }

      c.logger.Debugf("wrote new data to storage backend")

      c.events <- &entity.Event{
        Type: entity.ProfileIdEvent,
        Key: &entity.ProfileIdKey{
          Name: name,
          At:   at,
        },
        Object: id,
      }
      c.logger.Debugf("notified event channel")
    } else {
      c.logger.Debugf("cannot find resource on upstream")
    }
  } else {
    c.logger.Debugf("query fulfilled using cached data")
  }
  return id, nil
}

// resolves multiple profile associations at the current time
func (c *Cache) BulkGetProfileId(names []string) ([]*entity.ProfileId, error) {
  c.logger.Debugf("processing query for profile Ids associated with names %s", strings.Join(names, ", "))

  ids := make([]*entity.ProfileId, 0)
  at := time.Now()

  for i := 0; i < len(names); {
    name := names[i]
    id, err := c.storage.GetProfileId(name, at) // TODO: bulk lookup support in storage backend?
    if err != nil {
      c.logger.Errorf("storage backend responded with error: %s", err)
      continue
    }

    if id != nil {
      ids = append(ids, id)
      names = append(names[:i], names[i+1:]...)
      continue
    }

    i++
  }
  c.logger.Debugf("resolved %d profile Ids from cache, %d will be resolved from upstream", len(ids), len(names))

  if len(names) == 0 {
    c.logger.Debugf("query fulfilled using cached data")
    return ids, nil
  }

  c.incrementRequestCounter()
  newIds, err := c.upstream.BulkGetId(names)
  if err != nil {
    return nil, fmt.Errorf("upstream responded with error: %s", err)
  }

  for _, id := range newIds {
    err := c.storage.PutProfileId(id) // TODO: bulk upload support in storage backend?
    if err != nil {
      return nil, fmt.Errorf("storage backend responded with error: %s", err)
    }

    c.events <- &entity.Event{
      Type: entity.ProfileIdEvent,
      Key: &entity.ProfileIdKey{
        Name: id.Name,
        At:   at,
      },
      Object: id,
    }
  }

  c.logger.Debugf("wrote new data to storage backend")
  c.logger.Debugf("notified event channel")
  return append(ids, newIds...), nil
}

// purges the profile association of a given name at a given time
func (c *Cache) PurgeProfileId(name string, at time.Time) error {
  c.logger.Debugf("purging name association for name \"%s\" at time %s", name, at)
  return c.storage.PurgeProfileId(name, at)
}

// retrieves the name history of a given profile
func (c *Cache) GetNameHistory(id uuid.UUID) (*entity.NameChangeHistory, error) {
  c.logger.Debugf("processing query for name history of profile %s", id)

  history, err := c.storage.GetNameHistory(id)
  if err != nil {
    c.logger.Errorf("storage backend responded with an error: %s", err)
    history = nil
  }
  if history == nil {
    c.logger.Debugf("cache miss - requesting update from upstream")

    c.incrementRequestCounter()
    history, err = c.upstream.GetHistory(id)
    if err != nil {
      return nil, fmt.Errorf("upstream responded with error: %s", err)
    }

    if history != nil {
      err := c.storage.PutNameHistory(id, history)
      if err != nil {
        return nil, fmt.Errorf("storage backend responded with error: %s", err)
      }
      c.logger.Debugf("wrote new data to storage backend")

      c.events <- &entity.Event{
        Type:   entity.NameHistoryEvent,
        Key:    &id,
        Object: history,
      }
      c.logger.Debugf("notified event channel")
    } else {
      c.logger.Debugf("cannot find resource on upstream")
    }
  } else {
    c.logger.Debugf("query fulfilled using cached data")
  }
  return history, nil
}

// purges a name history from the cache
func (c *Cache) PurgeNameHistory(id uuid.UUID) error {
  c.logger.Debugf("purging name history for profile %s", id)
  return c.storage.PurgeNameHistory(id)
}

// retrieves a single profile
func (c *Cache) GetProfile(id uuid.UUID) (*entity.Profile, error) {
  c.logger.Debugf("processing query for profile %s", id)

  profile, err := c.storage.GetProfile(id)
  if err != nil {
    c.logger.Errorf("storage backend responded with an error: %s", err)
    profile = nil
  }
  if profile == nil {
    c.logger.Debugf("cache miss - requesting update from upstream")

    profile, err = c.upstream.GetProfile(id)
    if err != nil {
      return nil, fmt.Errorf("upstream responded with error: %s", err)
    }

    if profile != nil {
      err := c.storage.PutProfile(profile)
      if err != nil {
        return nil, fmt.Errorf("storage backend responded with error: %s", err)
      }

      err = c.updateNameMapping(profile)
      if err != nil {
        return nil, fmt.Errorf("storage backend responded with error: %s", err)
      }
      c.logger.Debugf("wrote new data to storage backend")

      c.events <- &entity.Event{
        Type:   entity.ProfileEvent,
        Key:    &id,
        Object: profile,
      }
      c.logger.Debugf("notified event channel")
    } else {
      c.logger.Debugf("cannot find resource on upstream")
    }
  } else {
    c.logger.Debugf("query fulfilled using cached data")
  }
  return profile, nil
}

// purges a specific profile from the cache
func (c *Cache) PurgeProfile(id uuid.UUID) error {
  c.logger.Debugf("purging profile with id %s", id)
  return c.storage.PurgeProfile(id)
}
