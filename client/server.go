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

  "github.com/dotStart/Stockpile/entity"
  "github.com/dotStart/Stockpile/rpc"
  "github.com/golang/protobuf/ptypes/empty"
)

// queries a server for the full server blacklist
func (s *Stockpile) GetBlacklist() (*entity.Blacklist, error) {
  blacklist, err := s.serverService.GetBlacklist(context.Background(), &empty.Empty{})
  if err != nil {
    return nil, err
  }

  return rpc.BlacklistFromRpc(blacklist)
}

// asks a server to verify whether one or more of the given addresses have been blacklisted
// only matching addresses will be returned, when the list is empty, none matched
func (s *Stockpile) CheckBlacklist(addresses []string) ([]string, error) {
  result, err := s.serverService.CheckBlacklist(context.Background(), &rpc.CheckBlacklistRequest{
    Addresses: addresses,
  })
  if err != nil {
    return nil, err
  }

  return result.MatchedAddresses, nil
}

// performs a cache assisted login
// an empty string may be passed to the ip field to omit the ip check
func (s *Stockpile) Login(displayName string, serverId string, ip string) (*entity.Profile, error) {
  profile, err := s.serverService.Login(context.Background(), &rpc.LoginRequest{
    DisplayName: displayName,
    ServerId:    serverId,
    Ip:          ip,
  })
  if err != nil {
    return nil, err
  }

  return rpc.ProfileFromRpc(profile)
}
