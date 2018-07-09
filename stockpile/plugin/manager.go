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
  "io/ioutil"
  "os"
  "path/filepath"
  "runtime"
  "strings"

  "github.com/dotStart/Stockpile/stockpile/storage"
  "github.com/op/go-logging"
)

const pluginExt = ".so"

type Manager struct {
  logger  *logging.Logger
  path    string
  Context *Context
  Plugins []*Plugin
}

// creates a new empty plugin manager with the given base path
func NewManager(path string) *Manager {
  ctx := &Context{
    storage: make(map[string]StorageBackendFactory),
  }
  ctx.RegisterStorageBackend("mem", storage.NewMemoryStorageBackend)
  ctx.RegisterStorageBackend("file", storage.NewFileStorageBackend)

  return &Manager{
    logger:  logging.MustGetLogger("plugin"),
    path:    path,
    Context: ctx,
  }
}

// loads all plugins in the plugin directory
func (m *Manager) LoadAll() error {
  if !PluginsAvailable {
    m.logger.Warningf("plugins are unavailable on platform %s - Plugin manager startup has been skipped", runtime.GOOS)
    return nil
  }

  files, err := ioutil.ReadDir(m.path)
  if err != nil {
    if os.IsNotExist(err) {
      m.logger.Warningf("plugin directory \"%s\" does not exist", m.path)
      return nil
    }

    return err
  }

  for _, file := range files {
    if !strings.HasSuffix(file.Name(), pluginExt) {
      continue
    }

    path := filepath.Join(m.path, file.Name())
    m.Load(path)
  }

  return nil
}

// loads a plugin from the specified path
func (m *Manager) Load(path string) {
  plugin, err := Load(path)
  if err != nil {
    m.logger.Errorf("failed to load plugin from path \"%s\": %s", path, err)
  } else {
    m.Plugins = append(m.Plugins, plugin)
    err = m.Context.merge(plugin.Context)
    if err != nil {
      m.logger.Errorf("failed to register one or more components of plugin \"%s\" v%s (defined by file \"%s\"): %s", plugin.Metadata.Name, plugin.Metadata.Version, path, err)
    }

    m.logger.Infof("loaded plugin \"%s\" v%s from file %s", plugin.Metadata.Name, plugin.Metadata.Version, path)
  }
}
