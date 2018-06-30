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
  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/dotStart/Stockpile/stockpile/plugin"
  "github.com/dotStart/Stockpile/stockpile/server"
  "github.com/dotStart/Stockpile/stockpile/server/rpc"
  "github.com/op/go-logging"
  "golang.org/x/net/context"
)

type ServerServiceImpl struct {
  logger  *logging.Logger
  api     *mojang.MojangAPI
  cfg     *server.Config
  storage plugin.StorageBackend
}

func NewServerService(api *mojang.MojangAPI, cfg *server.Config, backend plugin.StorageBackend) *ServerServiceImpl {
  return &ServerServiceImpl{
    logger:  logging.MustGetLogger("server-srv"),
    api:     api,
    cfg:     cfg,
    storage: backend,
  }
}

func (s *ServerServiceImpl) getBlacklist() (*mojang.Blacklist, error) {
  blacklist, err := s.storage.GetBlacklist()
  if err != nil || blacklist == nil {
    if err != nil {
      s.logger.Errorf("Failed to retrieve blacklist from storage backend: %s", err)
    }

    s.logger.Debugf("Cache miss - Requesting update from upstream")
    blacklist, err = s.api.GetBlacklist()
    if err != nil {
      s.logger.Errorf("Failed to retrieve blacklist: %s", err)
      return nil, err
    }

    s.logger.Debugf("Updating cached version")
    err = s.storage.PutBlacklist(blacklist)
    if err != nil {
      s.logger.Errorf("Failed to push blacklist to storage backend: %s", err)
    }
  }

  return blacklist, err
}

func (s *ServerServiceImpl) GetBlacklist(context.Context, *rpc.EmptyRequest) (*rpc.Blacklist, error) {
  s.logger.Debugf("Processing request for complete server blacklist")
  blacklist, err := s.getBlacklist()
  if err != nil {
    return nil, err
  }

  return rpc.BlacklistToRpc(blacklist), nil
}

func (s *ServerServiceImpl) CheckBlacklist(_ context.Context, req *rpc.CheckBlacklistRequest) (*rpc.CheckBlacklistResponse, error) {
  s.logger.Debugf("Processing request to check blacklist for %d addresses", len(req.Addresses))
  blacklist, err := s.getBlacklist()
  if err != nil {
    return nil, err
  }

  matches := make([]string, 0)
  for _, addr := range req.Addresses {
    match, err := blacklist.IsBlacklisted(addr)
    if err != nil {
      return nil, err
    }

    if match {
      matches = append(matches, addr)
    }
  }
  return &rpc.CheckBlacklistResponse{
    MatchedAddresses: matches,
  }, nil
}

func (s *ServerServiceImpl) Login(_ context.Context, req *rpc.LoginRequest) (*rpc.Profile, error) {
  s.logger.Debugf("Processing request to login user \"%s\" with serverId %s and ip %s", req.DisplayName, req.ServerId, req.Ip)
  profile, err := s.api.Login(req.DisplayName, req.ServerId, req.Ip)
  if err != nil {
    s.logger.Errorf("Failed to perform login: %s", err)
    return nil, err
  }

  err = s.storage.PutProfile(profile)
  if err != nil {
    s.logger.Errorf("Failed to push profile to storage backend: %s", err)
  }
  // TODO: Update profile mappings

  return rpc.ProfileToRpc(profile), nil
}
