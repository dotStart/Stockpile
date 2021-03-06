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
package entity

import (
  "encoding/base64"
  "encoding/json"
  "io"
  "time"

  "github.com/google/uuid"
)

type Profile struct {
  Id         uuid.UUID
  Name       string
  Properties map[string]*ProfileProperty
  Textures   *ProfileTextures
}

type restProfile struct {
  Id         string             `json:"id"`
  Name       string             `json:"name"`
  Properties []*ProfileProperty `json:"properties"`
}

type serializableProfile struct {
  restProfile
  Textures *serializableProfileTextures
}

func (p *Profile) Serialize() ([]byte, error) {
  i := 0
  props := make([]*ProfileProperty, len(p.Properties))
  for _, prop := range p.Properties {
    props[i] = prop
    i++
  }

  var tex *serializableProfileTextures = nil
  if p.Textures != nil {
    tex = &serializableProfileTextures{
      Timestamp:   p.Textures.Timestamp.Unix(),
      ProfileId:   p.Textures.ProfileId.String(),
      ProfileName: p.Textures.ProfileName,
      Textures:    p.Textures.Textures,
    }
  }

  enc := serializableProfile{
    restProfile: restProfile{
      Id:         p.Id.String(),
      Name:       p.Name,
      Properties: props,
    },
    Textures: tex,
  }
  return json.Marshal(&enc)
}

func (p *Profile) Deserialize(enc []byte) error {
  parsed := serializableProfile{}
  err := json.Unmarshal(enc, &parsed)
  if err != nil {
    return err
  }

  id, err := uuid.Parse(parsed.Id)
  if err != nil {
    return err
  }
  p.Id = id
  p.Name = parsed.Name

  p.Properties = make(map[string]*ProfileProperty)
  for _, prop := range parsed.Properties {
    p.Properties[prop.Name] = prop
  }

  p.Textures = nil
  if parsed.Textures != nil {
    id, err := uuid.Parse(parsed.Textures.ProfileId)
    if err != nil {
      return err
    }

    p.Textures = &ProfileTextures{
      Timestamp:   time.Unix(parsed.Textures.Timestamp, 0),
      ProfileId:   id,
      ProfileName: parsed.Textures.ProfileName,
      Textures:    parsed.Textures.Textures,
    }
  }
  return nil
}

func (p *Profile) Read(reader io.Reader) error {
  parsed := restProfile{}
  err := json.NewDecoder(reader).Decode(&parsed)
  if err != nil {
    return err
  }

  id, err := ParseId(parsed.Id)
  if err != nil {
    return err
  }
  p.Id = id
  p.Name = parsed.Name

  p.Properties = make(map[string]*ProfileProperty)
  for _, prop := range parsed.Properties {
    p.Properties[prop.Name] = prop
  }

  p.Textures = nil
  texProp := p.Properties["textures"]
  if texProp != nil {
    extractedValue, err := base64.StdEncoding.DecodeString(texProp.Value)
    if err != nil {
      return err
    }

    parsedProp := restProfileTextures{}
    err = json.Unmarshal(extractedValue, &parsedProp)
    if err != nil {
      return err
    }

    id, err := ParseId(parsedProp.ProfileId)
    if err != nil {
      return err
    }

    textures := make(map[string]string)
    for key, spec := range parsedProp.Textures {
      textures[key] = spec.Url
    }

    p.Textures = &ProfileTextures{
      Timestamp:   time.Unix(parsedProp.Timestamp/1000, parsedProp.Timestamp%1000*1000000),
      ProfileId:   id,
      ProfileName: parsedProp.ProfileName,
      Textures:    textures,
    }
  }
  return nil
}

// converts a profile into its original REST representation
// TODO: this implementation is used purely for the legacy API and should be removed once the legacy
// API has been removed
func (p *Profile) Mojang() ([]byte, error) {
  enc := restProfile{}

  enc.Id = ToMojangId(p.Id)
  enc.Name = p.Name
  enc.Properties = make([]*ProfileProperty, 0)
  for _, property := range p.Properties {
    enc.Properties = append(enc.Properties, property)
  }

  return json.Marshal(enc)
}

type ProfileProperty struct {
  Name      string `json:"name"`
  Value     string `json:"value"`
  Signature string `json:"signature"`
}

func (p *ProfileProperty) Serialize() ([]byte, error) {
  return json.Marshal(p)
}

func (p *ProfileProperty) Deserialize(enc []byte) error {
  return json.Unmarshal(enc, p)
}

func (p *ProfileProperty) read(reader io.Reader) error {
  return json.NewDecoder(reader).Decode(p)
}

type ProfileTextures struct {
  Timestamp   time.Time
  ProfileId   uuid.UUID
  ProfileName string
  Textures    map[string]string
}

type restProfileTextures struct {
  Timestamp   int64                             `json:"timestamp"`
  ProfileId   string                            `json:"profileId"`
  ProfileName string                            `json:"profileName"`
  Textures    map[string]restProfileTextureSpec `json:"textures"`
}

type serializableProfileTextures struct {
  Timestamp   int64  `json:"timestamp"`
  ProfileId   string `json:"profileId"`
  ProfileName string `json:"profileName"`
  Textures    map[string]string
}

type restProfileTextureSpec struct {
  Url string `json:"url"`
}

func (t *ProfileTextures) Serialize() ([]byte, error) {
  enc := &serializableProfileTextures{
    Timestamp:   t.Timestamp.Unix(),
    ProfileId:   t.ProfileId.String(),
    ProfileName: t.ProfileName,
    Textures:    t.Textures,
  }

  return json.Marshal(enc)
}

func (t *ProfileTextures) Deserialize(enc []byte) error {
  parsed := serializableProfileTextures{}
  err := json.Unmarshal(enc, &parsed)
  if err != nil {
    return err
  }

  id, err := uuid.Parse(parsed.ProfileId)
  if err != nil {
    return err
  }

  t.Timestamp = time.Unix(parsed.Timestamp, 0)
  t.ProfileId = id
  t.ProfileName = parsed.ProfileName
  t.Textures = parsed.Textures
  return nil
}

func (t *ProfileTextures) Read(reader io.Reader) error {
  parsed := restProfileTextures{}
  err := json.NewDecoder(reader).Decode(&parsed)
  if err != nil {
    return err
  }

  id, err := ParseId(parsed.ProfileId)
  if err != nil {
    return err
  }

  t.Timestamp = time.Unix(parsed.Timestamp, 0)
  t.ProfileId = id
  t.ProfileName = parsed.ProfileName

  textures := make(map[string]string)
  for key, spec := range parsed.Textures {
    textures[key] = spec.Url
  }

  return nil
}
