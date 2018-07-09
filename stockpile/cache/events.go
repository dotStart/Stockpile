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

import "github.com/dotStart/Stockpile/entity"

type Listener struct {
  cache *Cache
  C     chan *entity.Event
}

// frees all resources associated with this listener
func (e *Listener) Close() {
  e.cache.removeListener(e)
}

// registers a new event listener with the cache
func (c *Cache) NewListener() *Listener {
  c.listenerMutex.Lock()
  defer c.listenerMutex.Unlock()

  listener := &Listener{
    cache: c,
    C:     make(chan *entity.Event),
  }
  c.listeners = append(c.listeners, listener)
  return listener
}

// removes a listener from the global list
func (c *Cache) removeListener(listener *Listener) {
  c.listenerMutex.Lock()
  defer c.listenerMutex.Unlock()

  for i, l := range c.listeners {
    if l == listener {
      c.listeners = append(c.listeners[:i], c.listeners[i+1:]...)
      break
    }
  }
}

// distributes cache events to all registered listeners
func (c *Cache) deliverEvents() {
  for e := range c.events {
    for _, listener := range c.listeners {
      listener.C <- e
    }
  }
}
