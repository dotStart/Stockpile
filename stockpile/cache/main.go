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
package cache

import (
  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/dotStart/Stockpile/stockpile/storage"
  "github.com/op/go-logging"
)

// provides an abstraction layer between callers, the caching system and the upstream API
type Cache struct {
  logger   *logging.Logger
  upstream *mojang.MojangAPI
  storage  storage.StorageBackend
}

// creates a new cache client using
func New(upstream *mojang.MojangAPI, storage storage.StorageBackend) *Cache {
  return &Cache{
    logger:   logging.MustGetLogger("cache"),
    upstream: upstream,
    storage:  storage,
  }
}

func (c *Cache) Close() error {
  return c.storage.Close()
}
