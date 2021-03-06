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
  "encoding/json"
  "io"
  "time"

  "github.com/google/uuid"
)

// represents a single profile id mapping between a display name and a mapping at a given time
// note that lastSeenAt and validUntil may be set to UNIX epoch if the initial mapping is requested
type ProfileId struct {
  Id          uuid.UUID
  Name        string
  FirstSeenAt time.Time
  LastSeenAt  time.Time
  ValidUntil  time.Time
}

// represents a Mojang compatible representation of a profile id
type restProfileId struct {
  Id   string `json:"id"`
  Name string `json:"name"`
}

// represents a serializable version of the profile Id object
type serializableProfileId struct {
  restProfileId
  FirstSeenAt int64 `json:"firstSeenAt"`
  LastSeenAt  int64 `json:"lastSeenAt"`
  ValidUntil  int64 `json:"validUntil"`
}

func (p *ProfileId) Serialize() ([]byte, error) {
  enc := serializableProfileId{
    restProfileId: restProfileId{
      Id:   p.Id.String(),
      Name: p.Name,
    },
    FirstSeenAt: p.FirstSeenAt.Unix(),
    LastSeenAt:  p.LastSeenAt.Unix(),
    ValidUntil:  p.ValidUntil.Unix(),
  }

  return json.Marshal(&enc)
}

func SerializeProfileIdArray(profileIds []*ProfileId) ([]byte, error) {
  enc := make([]serializableProfileId, len(profileIds))
  for i, profileId := range profileIds {
    enc[i] = serializableProfileId{
      restProfileId: restProfileId{
        Id:   profileId.Id.String(),
        Name: profileId.Name,
      },
      FirstSeenAt: profileId.FirstSeenAt.Unix(),
      LastSeenAt:  profileId.LastSeenAt.Unix(),
      ValidUntil:  profileId.ValidUntil.Unix(),
    }
  }
  return json.Marshal(&enc)
}

func (p *ProfileId) Deserialize(enc []byte) error {
  parsed := serializableProfileId{}
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
  p.FirstSeenAt = time.Unix(parsed.FirstSeenAt, 0)
  p.LastSeenAt = time.Unix(parsed.LastSeenAt, 0)
  p.ValidUntil = time.Unix(parsed.ValidUntil, 0)
  return nil
}

func DeserializeProfileIdArray(enc []byte) ([]*ProfileId, error) {
  parsed := make([]serializableProfileId, 0)
  err := json.Unmarshal(enc, &parsed)
  if err != nil {
    return nil, err
  }

  res := make([]*ProfileId, len(parsed))
  for i, profileId := range parsed {
    id, err := uuid.Parse(profileId.Id)
    if err != nil {
      return nil, err
    }

    res[i] = &ProfileId{
      Id:          id,
      Name:        profileId.Name,
      FirstSeenAt: time.Unix(profileId.FirstSeenAt, 0),
      LastSeenAt:  time.Unix(profileId.LastSeenAt, 0),
      ValidUntil:  time.Unix(profileId.ValidUntil, 0),
    }
  }
  return res, nil
}

func (p *ProfileId) Read(reader io.Reader, at time.Time) error {
  parsed := restProfileId{}
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
  p.FirstSeenAt = at
  p.LastSeenAt = at
  p.ValidUntil = CalculateNameGracePeriodEnd(at)
  return nil
}

func ReadProfileIdArray(reader io.Reader) ([]*ProfileId, error) {
  parsed := make([]restProfileId, 0)
  err := json.NewDecoder(reader).Decode(&parsed)
  if err != nil {
    return nil, err
  }

  at := time.Now()
  res := make([]*ProfileId, len(parsed))
  for i, profileId := range parsed {
    id, err := ParseId(profileId.Id)
    if err != nil {
      return nil, err
    }

    res[i] = &ProfileId{
      Id:          id,
      Name:        profileId.Name,
      FirstSeenAt: at,
      LastSeenAt:  at,
      ValidUntil:  CalculateNameGracePeriodEnd(at),
    }
  }
  return res, nil
}

// updates the time at which this id has been discovered
func (p *ProfileId) UpdateDiscovery(at time.Time) {
  p.FirstSeenAt = at
  p.LastSeenAt = at
  p.ValidUntil = CalculateNameGracePeriodEnd(at)
}

// updates the last time at which this id has been encountered and the respective expiration times
// if the passed time is set before the current last encounter, the method will return immediately
// without changing the profile state
func (p *ProfileId) UpdateExpiration(seen time.Time) {
  if p.LastSeenAt.After(seen) {
    return // nothing to do
  }

  p.LastSeenAt = seen
  p.ValidUntil = CalculateNameGracePeriodEnd(seen)
}

// evaluates whether the profile is still valid at the given time
func (p *ProfileId) IsValid(at time.Time) bool {
  return !p.FirstSeenAt.After(at) && p.ValidUntil.After(at)
}

// evaluates whether two profileIds theoretically overlap
//
// two profiles are considered to overlap if their validity period overlaps at any point in time or
// if their assignments are equal while less than 30 days have passed (e.g. it is impossible for
// another user to claim and unclaim the name in the meantime due to the grace period)
// TODO: I have no clue how and whether Mojang handles theft of names with content creators
func (p *ProfileId) IsOverlappingWith(other *ProfileId) bool {
  return p.IsValid(other.FirstSeenAt) || p.IsValid(other.ValidUntil) || (p.Id == other.Id && p.ValidUntil.Add(NameChangeRateLimitPeriod).After(p.FirstSeenAt))
}

// encapsulates a name history
type NameChangeHistory struct {
  History []*NameChange
}

func (h *NameChangeHistory) Serialize() ([]byte, error) {
  return SerializeNameChangeArray(h.History)
}

func (h *NameChangeHistory) Deserialize(enc []byte) error {
  history, err := DeserializeNameChangeArray(enc)
  if err != nil {
    return err
  }
  h.History = history
  return nil
}

func (h *NameChangeHistory) Read(reader io.Reader) error {
  history, err := ReadNameChangeArray(reader)
  if err != nil {
    return err
  }
  h.History = history
  return nil
}

// represents a single name change within a profile's history
// note that changedToAt and validUntil may be set to UNIX epoch when the entry represents the initial account name
type NameChange struct {
  Name        string
  ChangedToAt time.Time
  ValidUntil  time.Time
}

// represents a Mojang compatible version of a name change record
type restNameChange struct {
  Name        string `json:"name"`
  ChangedToAt int64  `json:"changedToAt"`
}

// represents a serializable version of the name change object
type serializableNameChange struct {
  restNameChange
  ValidUntil int64 `json:"validUntil"`
}

func (p *NameChange) Serialize() ([]byte, error) {
  enc := serializableNameChange{
    restNameChange: restNameChange{
      Name:        p.Name,
      ChangedToAt: p.ChangedToAt.Unix(),
    },
    ValidUntil: p.ValidUntil.Unix(),
  }

  return json.Marshal(&enc)
}

func SerializeNameChangeArray(history []*NameChange) ([]byte, error) {
  enc := make([]*serializableNameChange, len(history))
  for i, change := range history {
    enc[i] = &serializableNameChange{
      restNameChange: restNameChange{
        Name:        change.Name,
        ChangedToAt: change.ChangedToAt.Unix(),
      },
      ValidUntil: change.ValidUntil.Unix(),
    }
  }

  return json.Marshal(&enc)
}

func (p *NameChange) Deserialize(enc []byte) error {
  parsed := serializableNameChange{}
  err := json.Unmarshal(enc, &parsed)
  if err != nil {
    return err
  }

  p.Name = parsed.Name
  p.ChangedToAt = time.Unix(parsed.ChangedToAt, 0)
  p.ValidUntil = time.Unix(parsed.ValidUntil, 0)
  return nil
}

func DeserializeNameChangeArray(enc []byte) ([]*NameChange, error) {
  parsed := make([]serializableNameChange, 0)
  err := json.Unmarshal(enc, &parsed)
  if err != nil {
    return nil, err
  }

  res := make([]*NameChange, len(parsed))
  for i, change := range parsed {
    res[i] = &NameChange{
      Name:        change.Name,
      ChangedToAt: time.Unix(change.ChangedToAt, 0),
      ValidUntil:  time.Unix(change.ValidUntil, 0),
    }
  }
  return res, nil
}

func (p *NameChange) Read(reader io.Reader) error {
  parsed := restNameChange{}
  err := json.NewDecoder(reader).Decode(&parsed)
  if err != nil {
    return err
  }

  p.Name = parsed.Name
  p.ChangedToAt = time.Unix(parsed.ChangedToAt/1000, parsed.ChangedToAt%1000*1000000)
  p.ValidUntil = CalculateNameGracePeriodEnd(p.ChangedToAt)
  return nil
}

func ReadNameChangeArray(reader io.Reader) ([]*NameChange, error) {
  parsed := make([]*restNameChange, 0)
  err := json.NewDecoder(reader).Decode(&parsed)
  if err != nil {
    return nil, err
  }

  res := make([]*NameChange, len(parsed))
  for i, change := range parsed {
    at := time.Unix(change.ChangedToAt/1000, change.ChangedToAt%1000*1000000)

    res[i] = &NameChange{
      Name:        change.Name,
      ChangedToAt: at,
      ValidUntil:  CalculateNameGracePeriodEnd(at),
    }
  }
  return res, nil
}
