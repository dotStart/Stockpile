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
package plugin

import (
  "errors"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
  "time"

  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/dotStart/Stockpile/stockpile/server"
  "github.com/google/uuid"
  "github.com/hashicorp/hcl2/gohcl"
  "github.com/op/go-logging"
)

const lockExpiration = time.Minute * 5 // keep-alive occurs every minute, expiration after 5 minutes
const lockKeepalive = time.Minute

const dirPerms = 664  // rw-rw-r--
const filePerms = 664 // rw-rw-r--

type FileStorageBackend struct {
  cfg        *server.Config
  fileCfg    *FileStorageBackendCfg
  logger     *logging.Logger
  lockPath   string
  lockTicker *time.Ticker
}

type FileStorageBackendCfg struct {
  Path string `hcl:"path,attr"`
}

// creates a new file storage backend
func NewFileStorageBackend(cfg *server.Config) (StorageBackend, error) {
  fileCfg := &FileStorageBackendCfg{}
  gohcl.DecodeBody(cfg.Storage.Parameters, nil, fileCfg)

  lockPath := filepath.Join(fileCfg.Path, "storage.lock")
  _, err := os.Stat(fileCfg.Path)
  if err != nil {
    if !os.IsNotExist(err) {
      return nil, err
    }

    err = os.MkdirAll(fileCfg.Path, dirPerms)
    if err != nil {
      return nil, err
    }
  } else {
    stat, err := os.Stat(lockPath)
    if err == nil {
      if !stat.ModTime().Add(lockExpiration).Before(time.Now()) {
        return nil, errors.New("file storage directory is locked by another instance - verify whether another instance is running and delete 'storage.lock' if this issue persists")
      }
    }
    if !os.IsNotExist(err) {
      return nil, err
    }
  }

  err = ioutil.WriteFile(lockPath, []byte{}, filePerms)
  if err != nil {
    return nil, err
  }

  impl := &FileStorageBackend{
    cfg:        cfg,
    fileCfg:    fileCfg,
    logger:     logging.MustGetLogger("file"),
    lockPath:   lockPath,
    lockTicker: time.NewTicker(lockKeepalive),
  }
  go impl.updateLock()
  return impl, nil
}

// updates the modification time of the lock file periodically to prevent its automatic expiration
func (f *FileStorageBackend) updateLock() {
  for range f.lockTicker.C {
    f.logger.Debugf("Updating database lock")
    ioutil.WriteFile(f.lockPath, []byte{}, filePerms)
  }
}

// retrieves the data of a previously stored cache entry (given that it exists and is still
// considered valid in accordance with its ttl)
func (f *FileStorageBackend) getCacheEntry(category string, key string, ttl time.Duration) ([]byte, error) {
  path := filepath.Join(f.fileCfg.Path, category, strings.ToLower(key))

  stat, err := os.Stat(path)
  if os.IsNotExist(err) {
    return nil, nil
  }
  if err != nil {
    return nil, err
  }
  if ttl != -1 && stat.ModTime().Add(ttl).Before(time.Now()) {
    return nil, nil
  }

  return ioutil.ReadFile(path)
}

// creates or updates a cache entry
func (f *FileStorageBackend) putCacheEntry(category string, key string, data []byte) error {
  dir := filepath.Join(f.fileCfg.Path, category)
  path := filepath.Join(dir, strings.ToLower(key))

  _, err := os.Stat(dir)
  if err != nil {
    if !os.IsNotExist(err) {
      return err
    }

    os.MkdirAll(dir, dirPerms)
  }

  return ioutil.WriteFile(path, data, filePerms)
}

// purges a cache entry (if it exists)
func (f *FileStorageBackend) purgeCacheEntry(category string, key string) error {
  path := filepath.Join(f.fileCfg.Path, category, strings.ToLower(key))

  _, err := os.Stat(path)
  if err != nil {
    if os.IsNotExist(err) {
      return nil
    }

    return err
  }

  return os.Remove(path)
}

func (f *FileStorageBackend) GetProfileId(name string, at time.Time) (*mojang.ProfileId, error) {
  enc, err := f.getCacheEntry("name", calculateHash(name), f.cfg.Ttl.Name)
  if err != nil {
    return nil, err
  }

  ids, err := mojang.DeserializeProfileIdArray(enc)
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

func (f *FileStorageBackend) PutProfileId(profileId *mojang.ProfileId) error {
  key := calculateHash(profileId.Name)
  enc, err := f.getCacheEntry("name", key, f.cfg.Ttl.Name)
  if err != nil {
    return err
  }
  var entries []*mojang.ProfileId
  found := false
  if enc != nil {
    entries, err = mojang.DeserializeProfileIdArray(enc)
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
    entries = make([]*mojang.ProfileId, 0)
    found = false
  }

  if !found {
    entries = append(entries, profileId)
  }

  enc, err = mojang.SerializeProfileIdArray(entries)
  if err != nil {
    return err
  }
  return f.putCacheEntry("name", key, enc)
}

func (f *FileStorageBackend) PurgeProfileId(name string, at time.Time) error {
  key := calculateHash(name)
  enc, err := f.getCacheEntry("name", key, f.cfg.Ttl.NameHistory)
  if err != nil {
    return err
  }

  entries, err := mojang.DeserializeProfileIdArray(enc)
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
    f.purgeCacheEntry("name", key)
    return nil
  }

  enc, err = mojang.SerializeProfileIdArray(entries)
  if err != nil {
    return err
  }
  return f.putCacheEntry("name", key, enc)
}

func (f *FileStorageBackend) GetNameHistory(id uuid.UUID) (*mojang.NameChangeHistory, error) {
  enc, err := f.getCacheEntry("history", id.String(), f.cfg.Ttl.NameHistory)
  if err != nil {
    return nil, err
  }

  history := &mojang.NameChangeHistory{}
  err = history.Deserialize(enc)
  return history, err
}

func (f *FileStorageBackend) PutNameHistory(id uuid.UUID, history *mojang.NameChangeHistory) error {
  enc, err := history.Serialize()
  if err != nil {
    return err
  }

  return f.putCacheEntry("history", id.String(), enc)
}

func (f *FileStorageBackend) PurgeNameHistory(id uuid.UUID) error {
  return f.purgeCacheEntry("history", id.String())
}

func (f *FileStorageBackend) GetProfile(id uuid.UUID) (*mojang.Profile, error) {
  enc, err := f.getCacheEntry("profile", id.String(), f.cfg.Ttl.Profile)
  if err != nil {
    return nil, err
  }

  profile := &mojang.Profile{}
  err = profile.Deserialize(enc)
  return profile, err
}

func (f *FileStorageBackend) PutProfile(profile *mojang.Profile) error {
  enc, err := profile.Serialize()
  if err != nil {
    return err
  }

  return f.putCacheEntry("profile", profile.Id.String(), enc)
}

func (f *FileStorageBackend) PurgeProfile(id uuid.UUID) error {
  return f.purgeCacheEntry("profile", id.String())
}

// Server Data
func (f *FileStorageBackend) GetBlacklist() (*mojang.Blacklist, error) {
  enc, err := f.getCacheEntry("misc", "blacklist", f.cfg.Ttl.Blacklist)
  if err != nil {
    return nil, err
  }

  blacklist := &mojang.Blacklist{}
  err = blacklist.Deserialize(enc)
  return blacklist, err
}

func (f *FileStorageBackend) PutBlacklist(blacklist *mojang.Blacklist) error {
  enc, err := blacklist.Serialize()
  if err != nil {
    return err
  }

  return f.putCacheEntry("misc", "blacklist", enc)
}

func (f *FileStorageBackend) PurgeBlacklist() error {
  return f.purgeCacheEntry("misc", "blacklist")
}

func (f *FileStorageBackend) Close() error {
  f.lockTicker.Stop()
  return os.Remove(f.lockPath)
}
