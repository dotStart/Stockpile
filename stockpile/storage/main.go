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

  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/google/uuid"
)

// provides an abstraction layer between the application and a storage backend
type StorageBackend interface {
  Close() error

  // Profile Data
  GetProfileId(name string, at time.Time) (*mojang.ProfileId, error)
  PutProfileId(profileId *mojang.ProfileId) error
  PurgeProfileId(name string, at time.Time) error
  GetNameHistory(id uuid.UUID) (*mojang.NameChangeHistory, error)
  PutNameHistory(id uuid.UUID, history *mojang.NameChangeHistory) error
  PurgeNameHistory(id uuid.UUID) error
  GetProfile(id uuid.UUID) (*mojang.Profile, error)
  PutProfile(profile *mojang.Profile) error
  PurgeProfile(id uuid.UUID) error

  // Server Data
  GetBlacklist() (*mojang.Blacklist, error)
  PutBlacklist(blacklist *mojang.Blacklist) error
  PurgeBlacklist() error
}
