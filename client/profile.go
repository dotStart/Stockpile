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
package client

import (
  "context"
  "time"

  "github.com/dotStart/Stockpile/entity"
  "github.com/dotStart/Stockpile/rpc"
  "github.com/google/uuid"
)

// queries a server for the player profile which has been associated to a given name at a specific
// time
func (s *Stockpile) GetProfileId(name string, at time.Time) (*entity.ProfileId, error) {
  profileId, err := s.profileService.GetId(context.Background(), &rpc.GetIdRequest{
    Name:      name,
    Timestamp: at.Unix(),
  })
  if err != nil {
    return nil, err
  }

  if !profileId.IsPopulated() {
    return nil, nil
  }

  return rpc.ProfileIdFromRpc(profileId)
}

// queries the server for the player profiles which are currently associated to the given names
func (s *Stockpile) BulkGetProfileId(names []string) ([]*entity.ProfileId, error) {
  response, err := s.profileService.BulkGetId(context.Background(), &rpc.BulkIdRequest{
    Names: names,
  })
  if err != nil {
    return nil, err
  }

  return rpc.ProfileIdsFromRpcArray(response.Ids)
}

// queries the server for the complete name history of a given profile
func (s *Stockpile) GetNameHistory(id uuid.UUID) (*entity.NameChangeHistory, error) {
  history, err := s.profileService.GetNameHistory(context.Background(), &rpc.IdRequest{
    Id: id.String(),
  })
  if err != nil {
    return nil, err
  }

  if !history.IsPopulated() {
    return nil, nil
  }

  return rpc.NameHistoryFromRpc(history), nil
}

// queries the server for a given profile
func (s *Stockpile) GetProfile(id uuid.UUID) (*entity.Profile, error) {
  profile, err := s.profileService.GetProfile(context.Background(), &rpc.IdRequest{
    Id: id.String(),
  })
  if err != nil {
    return nil, err
  }

  if !profile.IsPopulated() {
    return nil, nil
  }

  return rpc.ProfileFromRpc(profile)
}
