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
package service

import (
  "errors"
  "time"

  "github.com/dotStart/Stockpile/entity"
  "github.com/dotStart/Stockpile/rpc"
  "github.com/dotStart/Stockpile/stockpile/cache"
  "github.com/op/go-logging"
  "golang.org/x/net/context"
)

type ProfileServiceImpl struct {
  logger *logging.Logger
  cache  *cache.Cache
}

func NewProfileService(cache *cache.Cache) (*ProfileServiceImpl) {
  return &ProfileServiceImpl{
    logger: logging.MustGetLogger("profile-srv"),
    cache:  cache,
  }
}

func (s *ProfileServiceImpl) GetId(_ context.Context, req *rpc.GetIdRequest) (*rpc.ProfileId, error) {
  at := time.Unix(req.Timestamp, 0)

  profile, err := s.cache.GetProfileId(req.Name, at)
  if err != nil {
    return nil, err
  }
  if profile == nil {
    return &rpc.ProfileId{}, nil
  }
  return rpc.ProfileIdToRpc(profile), nil
}

func (s *ProfileServiceImpl) BulkGetId(_ context.Context, req *rpc.BulkIdRequest) (*rpc.BulkIdResponse, error) {
  if len(req.Names) > 100 {
    return nil, errors.New("cannot process more than 100 names at once")
  }

  ids, err := s.cache.BulkGetProfileId(req.Names)
  if err != nil {
    return nil, err
  }

  return rpc.BulkIdsToRpc(ids), nil
}

func (s *ProfileServiceImpl) GetNameHistory(_ context.Context, req *rpc.IdRequest) (*rpc.NameHistory, error) {
  id, err := entity.ParseId(req.Id)
  if err != nil {
    return nil, err
  }

  history, err := s.cache.GetNameHistory(id)
  if err != nil {
    return nil, err
  }

  return rpc.NameHistoryToRpc(history), nil
}

func (s *ProfileServiceImpl) GetProfile(_ context.Context, req *rpc.IdRequest) (*rpc.Profile, error) {
  id, err := entity.ParseId(req.Id)
  if err != nil {
    return nil, err
  }

  profile, err := s.cache.GetProfile(id)
  if err != nil {
    return nil, err
  }
  if profile == nil {
    return &rpc.Profile{}, nil
  }

  return rpc.ProfileToRpc(profile), nil
}
