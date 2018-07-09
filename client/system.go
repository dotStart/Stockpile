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
  "github.com/dotStart/Stockpile/stockpile/plugin"
  "github.com/golang/protobuf/ptypes/empty"
)

// queries a server for its status
func (s *Stockpile) GetStatus() (*entity.Status, error) {
  status, err := s.systemService.GetStatus(context.Background(), &empty.Empty{})
  if err != nil {
    return nil, err
  }

  return rpc.StatusFromRpc(status), nil
}

// queries a server for its loaded plugins
func (s *Stockpile) GetPluginList() ([]*plugin.Metadata, error) {
  result, err := s.systemService.GetPlugins(context.Background(), &empty.Empty{})
  if err != nil {
    return nil, err
  }

  return rpc.PluginMetadataListFromRpc(result), nil
}
