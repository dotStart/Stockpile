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
  "github.com/dotStart/Stockpile/stockpile/metadata"
  "github.com/dotStart/Stockpile/stockpile/plugin"
  "github.com/dotStart/Stockpile/rpc"
  "github.com/golang/protobuf/ptypes/empty"
  "github.com/op/go-logging"
  "golang.org/x/net/context"
)

type SystemServiceImpl struct {
  logger *logging.Logger
  plugin *plugin.Manager
}

func NewSystemService(plugin *plugin.Manager) (*SystemServiceImpl) {
  return &SystemServiceImpl{
    logger: logging.MustGetLogger("system-srv"),
    plugin: plugin,
  }
}

func (s *SystemServiceImpl) GetStatus(context.Context, *empty.Empty) (*rpc.Status, error) {
  return &rpc.Status{
    Brand:          metadata.Brand(),
    Version:        metadata.Version(),
    VersionFull:    metadata.VersionFull(),
    CommitHash:     metadata.CommitHash(),
    BuildTimestamp: metadata.Timestamp().Unix(),
  }, nil
}

func (s *SystemServiceImpl) GetPlugins(context.Context, *empty.Empty) (*rpc.PluginList, error) {
  list := make([]*plugin.Metadata, len(s.plugin.Plugins))
  for i, p := range s.plugin.Plugins {
    list[i] = &p.Metadata
  }
  return rpc.PluginMetadataListToRpc(list), nil
}
