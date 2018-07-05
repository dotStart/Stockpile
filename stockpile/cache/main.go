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
  "sync/atomic"
  "time"

  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/dotStart/Stockpile/stockpile/storage"
  "github.com/op/go-logging"
)

// provides an abstraction layer between callers, the caching system and the upstream API
type Cache struct {
  logger   *logging.Logger
  upstream *mojang.MojangAPI
  storage  storage.StorageBackend

  resetTicker    *time.Ticker
  requestCounter uint64
}

// creates a new cache client using
func New(upstream *mojang.MojangAPI, storage storage.StorageBackend) *Cache {
  cache := &Cache{
    logger:      logging.MustGetLogger("cache"),
    upstream:    upstream,
    storage:     storage,
    resetTicker: time.NewTicker(time.Minute),
  }
  go cache.resetRequestCounter()
  return cache
}

// increments the request counter by one
func (c *Cache) incrementRequestCounter() {
  atomic.AddUint64(&c.requestCounter, 1)
}

// regularly clears the request counter
func (c *Cache) resetRequestCounter() {
  for range c.resetTicker.C {
    atomic.StoreUint64(&c.requestCounter, 0)
  }
}

// retrieves the amount of requests which have been submitted to the upstream servers within the
// last minute
func (c *Cache) GetRateLimitAllocation() uint64 {
  return atomic.LoadUint64(&c.requestCounter)
}

func (c *Cache) Close() error {
  c.resetTicker.Stop()
  return c.storage.Close()
}
