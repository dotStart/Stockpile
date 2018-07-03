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
  "errors"
  "io/ioutil"
  "os"
  "path/filepath"
  "strings"
  "time"

  "github.com/dotStart/Stockpile/stockpile/server"
  "github.com/hashicorp/hcl2/gohcl"
  "github.com/op/go-logging"
)

const lockExpiration = time.Minute * 5 // keep-alive occurs every minute, expiration after 5 minutes
const lockKeepalive = time.Minute

const dirPerms = 664  // rw-rw-r--
const filePerms = 664 // rw-rw-r--

type fileStorageBackendInterface struct {
  logger     *logging.Logger
  cfg        *FileStorageBackendCfg
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

  impl := &fileStorageBackendInterface{
    logger:     logging.MustGetLogger("file"),
    cfg:        fileCfg,
    lockPath:   lockPath,
    lockTicker: time.NewTicker(lockKeepalive),
  }
  go impl.updateLock()
  return NewEncodedStorageBackend(cfg, impl), nil
}

// updates the modification time of the lock file periodically to prevent its automatic expiration
func (f *fileStorageBackendInterface) updateLock() {
  for range f.lockTicker.C {
    f.logger.Debugf("Updating database lock")
    ioutil.WriteFile(f.lockPath, []byte{}, filePerms)
  }
}

func (f *fileStorageBackendInterface) GetCacheEntry(category string, key string, ttl time.Duration) ([]byte, error) {
  path := filepath.Join(f.cfg.Path, category, strings.ToLower(key))

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

func (f *fileStorageBackendInterface) PutCacheEntry(category string, key string, data []byte, ttl time.Duration) error {
  dir := filepath.Join(f.cfg.Path, category)
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

func (f *fileStorageBackendInterface) PurgeCacheEntry(category string, key string) error {
  path := filepath.Join(f.cfg.Path, category, strings.ToLower(key))

  _, err := os.Stat(path)
  if err != nil {
    if os.IsNotExist(err) {
      return nil
    }

    return err
  }

  return os.Remove(path)
}

func (f *fileStorageBackendInterface) Close() error {
  return nil
}
