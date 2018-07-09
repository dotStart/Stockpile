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

  "github.com/dotStart/Stockpile/entity"
)

// retrieves the current server blacklist
func (c *Cache) GetBlacklist() (*entity.Blacklist, error) {
  c.logger.Debugf("processing query for server blacklist")

  blacklist, err := c.storage.GetBlacklist()
  if err != nil {
    c.logger.Errorf("storage backend responded with an error: %s", err)
    blacklist = nil
  }
  if blacklist == nil {
    c.logger.Debugf("cache miss - requesting update from upstream")

    c.incrementRequestCounter()
    blacklist, err = c.upstream.GetBlacklist()
    if err != nil {
      return nil, fmt.Errorf("upstream responded with error: %s", err)
    }

    if blacklist != nil {
      err = c.storage.PutBlacklist(blacklist)
      if err != nil {
        return nil, fmt.Errorf("storage backend responded with error: %s", err)
      }
      c.logger.Debugf("wrote new data to storage backend")

      c.events <- &entity.Event{
        Type:   entity.BlacklistEvent,
        Key:    nil,
        Object: blacklist,
      }
      c.logger.Debugf("notified event channel")
    } else {
      c.logger.Debugf("cannot find resource on upstream")
    }
  } else {
    c.logger.Debugf("query fulfilled using cached data")
  }
  return blacklist, nil
}

func (c *Cache) PurgeBlacklist() error {
  c.logger.Debugf("purging blacklist")
  return c.storage.PurgeBlacklist()
}

// performs a cache assisted server login
func (c *Cache) Login(displayName string, serverId string, ip string) (*entity.Profile, error) {
  c.logger.Debugf("processing login for user \"%s\" on server \"%s\" (with address \"%s\")", displayName, serverId, ip)

  profile, err := c.upstream.Login(displayName, serverId, ip)
  if err != nil {
    return nil, fmt.Errorf("upstream responded with error: %s", err)
  }

  err = c.storage.PutProfile(profile)
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
    Key:    profile.Id,
    Object: profile,
  }
  c.logger.Debugf("notified event channel")

  return profile, nil
}
