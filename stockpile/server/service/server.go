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
  "github.com/dotStart/Stockpile/stockpile/cache"
  "github.com/dotStart/Stockpile/stockpile/server/rpc"
  "github.com/op/go-logging"
  "golang.org/x/net/context"
  empty "github.com/golang/protobuf/ptypes/empty"
)

type ServerServiceImpl struct {
  logger *logging.Logger
  cache  *cache.Cache
}

func NewServerService(cache *cache.Cache) *ServerServiceImpl {
  return &ServerServiceImpl{
    logger: logging.MustGetLogger("server-srv"),
    cache:  cache,
  }
}

func (s *ServerServiceImpl) GetBlacklist(context.Context, *empty.Empty) (*rpc.Blacklist, error) {
  blacklist, err := s.cache.GetBlacklist()
  if err != nil {
    return nil, err
  }

  return rpc.BlacklistToRpc(blacklist), nil
}

func (s *ServerServiceImpl) CheckBlacklist(_ context.Context, req *rpc.CheckBlacklistRequest) (*rpc.CheckBlacklistResponse, error) {
  blacklist, err := s.cache.GetBlacklist()
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
  profile, err := s.cache.Login(req.DisplayName, req.ServerId, req.Ip)
  if err != nil {
    return nil, err
  }

  return rpc.ProfileToRpc(profile), nil
}
