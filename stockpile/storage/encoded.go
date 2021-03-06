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
package storage

import (
  "time"

  "github.com/dotStart/Stockpile/entity"
  "github.com/dotStart/Stockpile/stockpile/server"
  "github.com/google/uuid"
)

// provides a storage backend which provides format agnostic storage of cache entries for all
// entry types
type EncodedStorageBackend struct {
  cfg  *server.Config
  impl EncodedStorageBackendInterface
}

type EncodedStorageBackendInterface interface {
  // retrieves the data of a previously stored cache entry (given that it exists and is still
  // considered valid in accordance with its ttl)
  GetCacheEntry(category string, name string, ttl time.Duration) ([]byte, error)
  // creates or updates a cache entry
  PutCacheEntry(category string, name string, encoded []byte, ttl time.Duration) error
  // purges a cache entry (if it exists)
  PurgeCacheEntry(category string, name string) error

  // clears all allocated resources
  Close() error
}

func NewEncodedStorageBackend(cfg *server.Config, impl EncodedStorageBackendInterface) *EncodedStorageBackend {
  return &EncodedStorageBackend{
    cfg:  cfg,
    impl: impl,
  }
}

func (f *EncodedStorageBackend) GetProfileId(name string, at time.Time) (*entity.ProfileId, error) {
  enc, err := f.impl.GetCacheEntry("name", calculateHash(name), f.cfg.Ttl.Name)
  if err != nil {
    return nil, err
  }
  if enc == nil {
    return nil, nil
  }

  ids, err := entity.DeserializeProfileIdArray(enc)
  if err != nil {
    return nil, err
  }

  for _, id := range ids {
    if id.IsValid(at) {
      return id, nil
    }
  }

  return nil, nil
}

func (f *EncodedStorageBackend) PutProfileId(profileId *entity.ProfileId) error {
  key := calculateHash(profileId.Name)
  enc, err := f.impl.GetCacheEntry("name", key, f.cfg.Ttl.Name)
  if err != nil {
    return err
  }
  var entries []*entity.ProfileId
  found := false
  if enc != nil {
    entries, err = entity.DeserializeProfileIdArray(enc)
    if err != nil {
      return err
    }

    // if there is an overlap (e.g. the passed profile was encountered while during the local
    // validity period or was assigned to the same profile within the 30 day limitation period)
    for _, e := range entries {
      if e.IsOverlappingWith(profileId) {
        e.UpdateExpiration(profileId.LastSeenAt)
        found = true
        break
      }
    }
  } else {
    entries = make([]*entity.ProfileId, 0)
    found = false
  }

  if !found {
    entries = append(entries, profileId)
  }

  enc, err = entity.SerializeProfileIdArray(entries)
  if err != nil {
    return err
  }
  return f.impl.PutCacheEntry("name", key, enc, f.cfg.Ttl.Name)
}

func (f *EncodedStorageBackend) PurgeProfileId(name string, at time.Time) error {
  key := calculateHash(name)
  enc, err := f.impl.GetCacheEntry("name", key, f.cfg.Ttl.Name)
  if err != nil {
    return err
  }
  if enc == nil {
    return nil
  }

  entries, err := entity.DeserializeProfileIdArray(enc)
  if err != nil {
    return err
  }

  for i := 0; i < len(entries); {
    e := entries[i]
    if e.IsValid(at) {
      entries = append(entries[:i], entries[i+1:]...)
      break
    }
    i++
  }

  if len(entries) == 0 {
    f.impl.PurgeCacheEntry("name", key)
    return nil
  }

  enc, err = entity.SerializeProfileIdArray(entries)
  if err != nil {
    return err
  }
  return f.impl.PutCacheEntry("name", key, enc, f.cfg.Ttl.Name)
}

func (f *EncodedStorageBackend) GetNameHistory(id uuid.UUID) (*entity.NameChangeHistory, error) {
  enc, err := f.impl.GetCacheEntry("history", id.String(), f.cfg.Ttl.NameHistory)
  if err != nil {
    return nil, err
  }
  if enc == nil {
    return nil, nil
  }

  history := &entity.NameChangeHistory{}
  err = history.Deserialize(enc)
  return history, err
}

func (f *EncodedStorageBackend) PutNameHistory(id uuid.UUID, history *entity.NameChangeHistory) error {
  enc, err := history.Serialize()
  if err != nil {
    return err
  }

  return f.impl.PutCacheEntry("history", id.String(), enc, f.cfg.Ttl.NameHistory)
}

func (f *EncodedStorageBackend) PurgeNameHistory(id uuid.UUID) error {
  return f.impl.PurgeCacheEntry("history", id.String())
}

func (f *EncodedStorageBackend) GetProfile(id uuid.UUID) (*entity.Profile, error) {
  enc, err := f.impl.GetCacheEntry("profile", id.String(), f.cfg.Ttl.Profile)
  if err != nil {
    return nil, err
  }
  if enc == nil {
    return nil, nil
  }

  profile := &entity.Profile{}
  err = profile.Deserialize(enc)
  return profile, err
}

func (f *EncodedStorageBackend) PutProfile(profile *entity.Profile) error {
  enc, err := profile.Serialize()
  if err != nil {
    return err
  }

  return f.impl.PutCacheEntry("profile", profile.Id.String(), enc, f.cfg.Ttl.Profile)
}

func (f *EncodedStorageBackend) PurgeProfile(id uuid.UUID) error {
  return f.impl.PurgeCacheEntry("profile", id.String())
}

// Server Data
func (f *EncodedStorageBackend) GetBlacklist() (*entity.Blacklist, error) {
  enc, err := f.impl.GetCacheEntry("misc", "blacklist", f.cfg.Ttl.Blacklist)
  if err != nil {
    return nil, err
  }
  if enc == nil {
    return nil, nil
  }

  blacklist := &entity.Blacklist{}
  err = blacklist.Deserialize(enc)
  return blacklist, err
}

func (f *EncodedStorageBackend) PutBlacklist(blacklist *entity.Blacklist) error {
  enc, err := blacklist.Serialize()
  if err != nil {
    return err
  }

  return f.impl.PutCacheEntry("misc", "blacklist", enc, f.cfg.Ttl.Blacklist)
}

func (f *EncodedStorageBackend) PurgeBlacklist() error {
  return f.impl.PurgeCacheEntry("misc", "blacklist")
}

func (f *EncodedStorageBackend) Close() error {
  return f.impl.Close()
}
