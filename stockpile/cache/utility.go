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
  "time"

  "github.com/dotStart/Stockpile/stockpile/mojang"
)

// adjusts the name associations for the data discovered through a profile request
func (c *Cache) updateNameMapping(profile *mojang.Profile) error {
  at := time.Now()

  mapping := &mojang.ProfileId{
    Id:          profile.Id,
    Name:        profile.Name,
    FirstSeenAt: at,
  }
  mapping.UpdateExpiration(at)

  c.storage.PutProfileId(mapping)

  c.Events <- &Event{
    Type: ProfileIdEvent,
    Key: &ProfileIdKey{
      Name: profile.Name,
      At:   at,
    },
    Object: mapping,
  }
  return nil
}
