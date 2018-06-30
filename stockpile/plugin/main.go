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
package plugin

import (
  "errors"
  "plugin"
  "runtime"
  "strings"

  "github.com/dotStart/Stockpile/stockpile/server"
)

// represents the metadata associated with a plugin implementation
type Metadata struct {
  Name    string
  Authors []string
  Website string
}

// represents a loaded plugin
type StockpilePlugin struct {
  handle         *plugin.Plugin
  storageFactory func(backend *server.Config) (StorageBackend, error)

  Meta Metadata
}

// opens an arbitrary plugin
func Open(path string) (*StockpilePlugin, error) {
  if runtime.GOOS == "windows" {
    return nil, errors.New("plugins are not supported on windows")
  }

  if !strings.HasSuffix(path, ".so") {
    path += ".so"
  }

  handle, err := plugin.Open(path)
  if err != nil {
    return nil, err
  }

  metaSymbol, err := handle.Lookup("GetMetadata")
  if err != nil {
    return nil, err
  }

  var storageFactory func(backend *server.Config) (StorageBackend, error)
  storageFactorySymbol, err := handle.Lookup("CreateStorageBackend")
  if err != nil {
    storageFactory = nil
  } else {
    storageFactory = storageFactorySymbol.(func(backend *server.Config) (StorageBackend, error))
  }

  return &StockpilePlugin{
    handle:         handle,
    storageFactory: storageFactory,

    Meta: metaSymbol.(func() (Metadata))(),
  }, nil
}

// evaluates whether the plugin provides a storage backend implementation
func (p *StockpilePlugin) HasStorageBackendImplementation() bool {
  return p.storageFactory != nil
}

// creates a new storage backend using the plugin's registered provider
func (p *StockpilePlugin) CreateStorageBackend(cfg *server.Config) (StorageBackend, error) {
  return p.storageFactory(cfg)
}
