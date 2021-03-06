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
  "strings"
  "time"

  "github.com/dotStart/Stockpile/entity"
  "github.com/dotStart/Stockpile/stockpile/server"
  "github.com/google/uuid"
  "github.com/op/go-logging"
)

type MemoryStorageBackend struct {
  cfg    *server.Config
  logger *logging.Logger

  profileId   map[string][]expirationWrapper
  nameHistory map[uuid.UUID]*expirationWrapper
  profile     map[uuid.UUID]*expirationWrapper

  blacklist *expirationWrapper
}

// creates a new memory based storage backend
func NewMemoryStorageBackend(cfg *server.Config) (StorageBackend, error) {
  return &MemoryStorageBackend{
    cfg:    cfg,
    logger: logging.MustGetLogger("memdb"),

    profileId:   make(map[string][]expirationWrapper),
    nameHistory: make(map[uuid.UUID]*expirationWrapper),
    profile:     make(map[uuid.UUID]*expirationWrapper),
  }, nil
}

func (m *MemoryStorageBackend) Close() error {
  return nil
}

func (m *MemoryStorageBackend) GetProfileId(name string, at time.Time) (*entity.ProfileId, error) {
  m.clearExpiredEntries()

  name = strings.ToLower(name)
  m.logger.Debugf("checking profile associations for \"%s\" at time %s", name, at)
  mappings := m.profileId[name]
  if mappings == nil {
    m.logger.Debugf("no associations for \"%s\"", name)
    return nil, nil
  }

  for _, exp := range mappings {
    association := exp.content.(*entity.ProfileId)
    if association.IsValid(at) {
      m.logger.Debugf("association to profile %s matches", association.Id)
      return association, nil
    } else {
      m.logger.Debugf("association to profile %s is invalid", association.Id)
    }
  }

  return nil, nil
}

func (m *MemoryStorageBackend) PutProfileId(profileId *entity.ProfileId) error {
  m.clearExpiredEntries()

  name := strings.ToLower(profileId.Name)
  m.logger.Debugf("updating association for name \"%s\" to profile %s at time %s (valid until %s)", profileId.Name, profileId.Id, profileId.LastSeenAt, profileId.ValidUntil)
  mappings := m.profileId[name]
  found := false
  if mappings != nil {
    for _, e := range mappings {
      entry := e.content.(*entity.ProfileId)
      if entry.IsOverlappingWith(profileId) {
        entry.UpdateExpiration(profileId.LastSeenAt)
        found = true
      }
    }
  } else {
    mappings = make([]expirationWrapper, 0)
  }

  if !found {
    mappings = append(mappings, expirationWrapper{
      content:   profileId,
      createdAt: time.Now(),
    })
  }

  m.profileId[name] = mappings
  return nil
}

func (m *MemoryStorageBackend) PurgeProfileId(name string, at time.Time) error {
  m.clearExpiredEntries()

  m.logger.Debugf("purging profile associations for \"%s\" at time %s", name, at)
  mappings := m.profileId[name]
  if mappings == nil {
    m.logger.Debugf("No associations for \"%s\"", name)
    return nil
  }

  if at.Unix() == -1 {
    delete(m.profileId, name)
    return nil
  }

  for i := 0; i < len(mappings); {
    exp := mappings[i]
    association := exp.content.(*entity.ProfileId)

    if association.IsValid(at) {
      m.logger.Debugf("purging association to profile %s", association.Id)
      mappings = append(mappings[:i], mappings[:i+1]...)
      continue
    }

    i++
  }

  m.profileId[name] = mappings
  return nil
}

func (m *MemoryStorageBackend) GetNameHistory(id uuid.UUID) (*entity.NameChangeHistory, error) {
  m.clearExpiredEntries()

  exp := m.nameHistory[id]
  if exp == nil {
    return nil, nil
  }

  return exp.content.(*entity.NameChangeHistory), nil
}

func (m *MemoryStorageBackend) PutNameHistory(id uuid.UUID, history *entity.NameChangeHistory) error {
  m.clearExpiredEntries()

  m.logger.Debugf("storing history for profile %s (consisting of %d elements)", id, len(history.History))
  m.nameHistory[id] = &expirationWrapper{
    content:   history,
    createdAt: time.Now(),
  }

  return nil
}

func (m *MemoryStorageBackend) PurgeNameHistory(id uuid.UUID) error {
  m.clearExpiredEntries()

  m.logger.Debugf("purging history for profile %s", id)
  delete(m.nameHistory, id)
  return nil
}

func (m *MemoryStorageBackend) GetProfile(id uuid.UUID) (*entity.Profile, error) {
  m.clearExpiredEntries()

  exp := m.profile[id]
  if exp == nil {
    return nil, nil
  }

  return exp.content.(*entity.Profile), nil
}

func (m *MemoryStorageBackend) PutProfile(profile *entity.Profile) error {
  m.clearExpiredEntries()

  m.logger.Debugf("storing profile %s", profile.Id)
  m.profile[profile.Id] = &expirationWrapper{
    content:   profile,
    createdAt: time.Now(),
  }

  return nil
}

func (m *MemoryStorageBackend) PurgeProfile(id uuid.UUID) error {
  m.clearExpiredEntries()

  m.logger.Debugf("purging profile %s", id)
  delete(m.profile, id)
  return nil
}

func (m *MemoryStorageBackend) GetBlacklist() (*entity.Blacklist, error) {
  if m.blacklist == nil {
    return nil, nil
  }

  return m.blacklist.content.(*entity.Blacklist), nil
}

func (m *MemoryStorageBackend) PutBlacklist(blacklist *entity.Blacklist) error {
  m.blacklist = &expirationWrapper{
    content:   blacklist,
    createdAt: time.Now(),
  }
  return nil
}

func (m *MemoryStorageBackend) PurgeBlacklist() error {
  m.logger.Debugf("purging blacklist")
  m.blacklist = nil
  return nil
}

// clears all expired entries from the database
func (m *MemoryStorageBackend) clearExpiredEntries() { // TODO: run on a timer instead?
  m.logger.Debug("purging expired data")

  deletedProfileIds := 0
  deletedNameHistories := 0
  deletedProfiles := 0
  deletedBlacklists := 0

  for profileId, mappings := range m.profileId {
    for i := 0; i < len(mappings); {
      exp := mappings[i]

      if !exp.isValid(m.cfg.Ttl.Name) {
        deletedProfileIds++

        if len(mappings) == 1 {
          delete(m.profileId, profileId)
          continue
        }

        mappings = append(mappings[:i], mappings[i+1:]...)
        continue
      }

      i++
    }
  }

  for key, history := range m.nameHistory {
    if !history.isValid(m.cfg.Ttl.NameHistory) {
      deletedNameHistories++
      delete(m.nameHistory, key)
    }
  }

  for key, profile := range m.profile {
    if !profile.isValid(m.cfg.Ttl.Profile) {
      deletedProfiles++
      delete(m.profile, key)
    }
  }

  if m.blacklist != nil && m.blacklist.isValid(m.cfg.Ttl.Blacklist) {
    deletedBlacklists = 1
    m.blacklist = nil
  }

  m.logger.Debugf("removed %d profile Ids, %d name histories, %d profiles and %d blacklists from memory", deletedProfileIds, deletedNameHistories, deletedProfiles, deletedBlacklists)
}
