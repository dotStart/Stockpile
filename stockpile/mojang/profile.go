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
  "encoding/json"
  "fmt"
  "time"

  "github.com/google/uuid"
)

type Profile struct {
  Id            uuid.UUID
  RawId         string            `json:"id"`
  Name          string
  Properties    map[string]*ProfileProperty
  RawProperties []ProfileProperty `json:"properties"`
  Textures      *ProfileTextures
}

type ProfileProperty struct {
  Name      string `json:"name"`
  Value     string `json:"value"`
  Signature string `json:"signature"`
}

type ProfileTextures struct {
  Timestamp    time.Time
  RawTimestamp int64                         `json:"timestamp"`
  ProfileId    uuid.UUID
  RawProfileId string                        `json:"profileId"`
  ProfileName  string                        `json:"profileName"`
  Textures     map[string]string
  RawTextures  map[string]ProfileTextureSpec `json:"textures"`
}

type ProfileTextureSpec struct {
  Url string `json:"url"`
}

func (p *Profile) init() error {
  id, err := ToStandardId(p.RawId)
  if err != nil {
    return err
  }

  p.Id = id

  p.Properties = make(map[string]*ProfileProperty)
  for _, prop := range p.RawProperties {
    p.Properties[prop.Name] = &prop
  }

  textures := p.Properties["textures"]
  if textures != nil {
    p.Textures = &ProfileTextures{}
    err = json.Unmarshal([]byte(textures.Value), p.Textures)
    if err != nil {
      return err
    }
  }

  return nil
}

func (p *ProfileTextures) init() error {
  id, err := ToStandardId(p.RawProfileId)
  if err != nil {
    return err
  }

  p.Timestamp = time.Unix(p.RawTimestamp/1000, p.RawTimestamp%1000*1000000)
  p.ProfileId = id

  p.Textures = make(map[string]string)
  for key, spec := range p.RawTextures {
    p.Textures[key] = spec.Url
  }

  return nil
}

// retrieves a single profile from the server
func (a *MojangAPI) GetProfile(id uuid.UUID) (*Profile, error) {
  res, err := a.execute("GET", fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", id.String()), nil)
  if err != nil {
    return nil, err
  }

  if res.StatusCode == 204 || res.StatusCode == 404 {
    return nil, nil
  }

  profile := &Profile{}
  defer res.Body.Close()
  err = json.NewDecoder(res.Body).Decode(profile)
  if err != nil {
    return nil, err
  }

  err = profile.init()
  if err != nil {
    return nil, err
  }
  return profile, nil
}
