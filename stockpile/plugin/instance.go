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
  "fmt"
  "plugin"
  "reflect"
  "runtime"
)

// represents a loaded plugin instance
type Plugin struct {
  handle *plugin.Plugin

  Metadata Metadata
  Context  *Context
}

// loads an arbitrary plugin from the specified location
func Load(path string) (*Plugin, error) {
  // at the moment the go plugin architecture isn't particularly consistent so we won't be able to
  // load anything unless we're on Linux or Mac OS
  if !PluginsAvailable {
    return nil, fmt.Errorf("plugins are not supported on platform %s", runtime.GOOS)
  }

  handle, err := plugin.Open(path)
  if err != nil {
    return nil, fmt.Errorf("cannot open plugin \"%s\": %s", path, err)
  }

  metadataHandle, err := handle.Lookup("Metadata")
  if err != nil {
    return nil, fmt.Errorf("plugin \"%s\" does not expose its metadata: %s", path, err)
  }

  metadata, ok := metadataHandle.(*Metadata)
  if !ok {
    return nil, fmt.Errorf("plugin \"%s\" defines an illegal metadata symbol: expected *Metadata but got %s", path, reflect.TypeOf(metadataHandle))
  }

  initializerHandle, err := handle.Lookup("InitializePlugin")
  if err != nil {
    return nil, fmt.Errorf("plugin \"%s\" does not expose an initializer: %s", path, err)
  }

  initializer, ok := initializerHandle.(Initializer)
  if !ok {
    return nil, fmt.Errorf("plugin \"%s\" defines an illegal initializer symbol: expected func(*plugin.Context) error but got %s", path, reflect.TypeOf(initializerHandle))
  }

  ctx := &Context{
    storage: make(map[string]StorageBackendFactory),
  }

  err = initializer(ctx)
  if err != nil {
    return nil, fmt.Errorf("plugin \"%s\" failed to initialize: %s", path, err)
  }

  return &Plugin{
    handle:   handle,
    Metadata: *metadata,
    Context:  ctx,
  }, nil
}

// represents the context associated with a given plugin
// we use this instance to simplify the registration of plugin implementations
type Context struct {
  storage map[string]StorageBackendFactory
}

// merges two context instances with each other
func (c *Context) merge(other *Context) error {
  for key, factory := range other.storage {
    current := c.storage[key]
    if current != nil {
      return fmt.Errorf("storage backend with identifier \"%s\" is already defined")
    }

    c.storage[key] = factory
  }

  return nil
}

// retrieves the storage backend factory for the specified identifier
func (c *Context) GetStorageBackend(id string) StorageBackendFactory {
  return c.storage[id]
}

// registers a new storage backend with the context
func (c *Context) RegisterStorageBackend(id string, factory StorageBackendFactory) error {
  current := c.storage[id]
  if current != nil {
    return fmt.Errorf("storage backend with id \"%s\" has already been registered", id)
  }

  c.storage[id] = factory
  return nil
}
