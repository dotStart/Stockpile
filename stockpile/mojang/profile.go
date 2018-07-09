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
package mojang

import (
  "fmt"

  "github.com/dotStart/Stockpile/entity"
  "github.com/google/uuid"
)

// retrieves a single profile from the server
func (a *MojangAPI) GetProfile(id uuid.UUID) (*entity.Profile, error) {
  res, err := a.execute("GET", fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s?unsigned=false", entity.ToMojangId(id)), nil)
  if err != nil {
    return nil, err
  }

  if res.StatusCode == 204 || res.StatusCode == 404 {
    return nil, nil
  }

  profile := &entity.Profile{}
  defer res.Body.Close()
  err = profile.Read(res.Body)
  return profile, err
}
