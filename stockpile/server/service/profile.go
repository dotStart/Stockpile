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

  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/dotStart/Stockpile/stockpile/server"
  "github.com/dotStart/Stockpile/stockpile/server/rpc"
  "github.com/dotStart/Stockpile/stockpile/storage"
  "github.com/op/go-logging"
  "golang.org/x/net/context"
)

type ProfileServiceImpl struct {
  logger  *logging.Logger
  api     *mojang.MojangAPI
  cfg     *server.Config
  storage storage.StorageBackend
}

func NewProfileService(api *mojang.MojangAPI, cfg *server.Config, backend storage.StorageBackend) (*ProfileServiceImpl) {
  return &ProfileServiceImpl{
    logger:  logging.MustGetLogger("profile-srv"),
    api:     api,
    cfg:     cfg,
    storage: backend,
  }
}

func (s *ProfileServiceImpl) GetId(_ context.Context, req *rpc.GetIdRequest) (*rpc.ProfileId, error) {
  at := time.Unix(req.Timestamp, 0)
  s.logger.Debugf("Processing request for profile id for name \"%s\" at time %s", req.Name, at)
  association, err := s.storage.GetProfileId(req.Name, at)
  if err != nil {
    s.logger.Errorf("Database responded with error during lookup of name \"%s\" at time %s: %s", req.Name, at.String(), err)
  } else if association != nil {
    s.logger.Debugf("Returning cached result of database: %+v", association)
  }
  if err != nil || association == nil {
    s.logger.Debugf("Cache miss - Requesting information from upstream")
    association, err = s.api.GetId(req.Name, time.Unix(req.Timestamp, 0))
    if err != nil {
      s.logger.Errorf("Failed to retrieve association of name \"%s\" at time %s: %s", req.Name, at.String(), err)
      return nil, err
    }

    if association != nil {
      s.logger.Debugf("Pushing updated information to cache backend")
      err = s.storage.PutProfileId(association)
      if err != nil {
        s.logger.Errorf("Failed to push profile association to storage backend: %s", err)
      }
    }
  }

  if association == nil {
    s.logger.Debugf("No profile with name \"%s\" at time %s", req.Name, at)
    return &rpc.ProfileId{}, nil
  }

  s.logger.Debugf("Display name \"%s\" resolved to profile %s at time %s", req.Name, association.Id, at)
  return rpc.ProfileIdToRpc(association), nil
}

func (s *ProfileServiceImpl) GetNameHistory(_ context.Context, req *rpc.IdRequest) (*rpc.NameHistory, error) {
  id, err := mojang.ParseId(req.Id)
  if err != nil {
    return nil, err
  }

  s.logger.Debugf("Processing request for name history for profile %s", id)
  history, err := s.storage.GetNameHistory(id)
  if err != nil {
    s.logger.Errorf("Database responded with error during lookup of history of profile \"%s\": %s", id, err)
  }
  if err != nil || history == nil {
    history, err = s.api.GetHistory(id)
    if err != nil {
      s.logger.Errorf("Failed to retrieve history of profile \"%s\": %s", id, err)
      return nil, err
    }

    if history == nil {
      s.logger.Debugf("No profile with id %s for name history request", id)
      return &rpc.NameHistory{}, nil
    }

    s.logger.Debugf("Updated history with %d elements from upstream", len(history.History))
    err = s.storage.PutNameHistory(id, history)
    if err != nil {
      s.logger.Errorf("Failed to push name history to storage backend: %s", err)
    }
  }

  s.logger.Debugf("Name history for profile %s consists of %d elements", id, len(history.History))
  return rpc.NameHistoryToRpc(history), nil
}

func (s *ProfileServiceImpl) BulkGetId(_ context.Context, req *rpc.BulkIdRequest) (*rpc.BulkIdResponse, error) {
  if len(req.Names) > 100 {
    return nil, errors.New("cannot process more than 100 names at once")
  }

  s.logger.Debugf("Processing request for ids of %d profiles", len(req.Names))

  names := make([]string, 0)
  results := make([]*mojang.ProfileId, 0)
  for _, name := range req.Names {
    profileId, err := s.storage.GetProfileId(name, time.Now())
    if profileId == nil || err != nil {
      names = append(names, name)

      if err != nil {
        s.logger.Errorf("Failed to resolve profileId for name \"%s\" from storage backend: %s", name, err)
      }
    } else {
      s.logger.Debugf("Resolved name \"%s\" to profile %s via cache", name, profileId.Id)
      results = append(results, profileId)
    }
  }
  s.logger.Debugf("Resolved %d out of %d names from cache (%d remain)", len(results), len(req.Names), len(names))

  if len(names) != 0 {
    upstreamResults, err := s.api.BulkGetId(names)
    if err != nil {
      s.logger.Errorf("Failed to resolve %d names from upstream: %s", len(names), err)
      return nil, err
    }

    for _, profileId := range upstreamResults {
      err := s.storage.PutProfileId(profileId)
      if err != nil {
        s.logger.Errorf("Failed to push profile association to storage backend: %s", err)
      }
    }

    results = append(results, upstreamResults...)
  }

  return rpc.BulkIdsToRpc(results), nil
}

func (s *ProfileServiceImpl) GetProfile(_ context.Context, req *rpc.IdRequest) (*rpc.Profile, error) {
  id, err := mojang.ParseId(req.Id)
  if err != nil {
    return nil, err
  }

  s.logger.Debugf("Processing request for profile %s", id)
  profile, err := s.storage.GetProfile(id)
  if err != nil || profile == nil {
    if err != nil {
      s.logger.Errorf("Failed to resolve profile %s from storage backend: %s", id, err)
    }

    s.logger.Debugf("Cache miss - Requesting update from upstream server")
    profile, err = s.api.GetProfile(id)
    if err != nil {
      s.logger.Errorf("Failed to resolve profile %s from upstream server: %s", id, err)
      return nil, err
    }
    if profile == nil {
      return &rpc.Profile{}, nil
    }

    err = s.storage.PutProfile(profile)
    if err != nil {
      s.logger.Errorf("Failed to push profile to storage backend: %s", err)
    }
    // TODO: Update profile mappings
  }

  return rpc.ProfileToRpc(profile), nil
}
