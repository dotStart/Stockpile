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
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "net/url"
  "time"

  "github.com/google/uuid"
)

var unixEpoch = time.Unix(0, 0)

// represents a single profile id mapping between a display name and a mapping at a given time
// note that lastSeenAt and validUntil may be set to UNIX epoch if the initial mapping is requested
type ProfileId struct {
  Id          uuid.UUID
  RawId       string `json:"id"`
  Name        string `json:"name"`
  FirstSeenAt time.Time
  LastSeenAt  time.Time
  ValidUntil  time.Time
}

// represents a single name change within a profile's history
// note that changedToAt and validUntil may be set to UNIX epoch when the entry represents the initial account name
type NameChange struct {
  Name           string `json:"name"`
  ChangedToAt    time.Time
  RawChangedToAt int64  `json:"changedToAt"`
  ValidUntil     time.Time
}

func (p *ProfileId) init(seen time.Time) error {
  id, err := ToStandardId(p.RawId)
  if err != nil {
    return err
  }

  p.Id = id
  p.FirstSeenAt = seen
  p.LastSeenAt = seen

  if seen != unixEpoch {
    p.ValidUntil = CalculateNameGracePeriodEnd(seen)
  } else {
    p.ValidUntil = unixEpoch
  }

  return nil
}

// evaluates whether the profile is still valid at the given time
func (p *ProfileId) IsValid(at time.Time) bool {
  return !p.FirstSeenAt.After(at) && p.ValidUntil.After(at)
}

func (c *NameChange) init() error {
  c.ChangedToAt = time.Unix(c.RawChangedToAt, 0)

  if c.RawChangedToAt != 0 {
    c.ValidUntil = CalculateNameGracePeriodEnd(c.ChangedToAt)
  } else {
    c.ValidUntil = unixEpoch
  }

  return nil
}

// retrieves the profile id (and some associated attributes) for a given display name at the specified time
// - if the UNIX epoch (e.g. zero) is passed instead of a real time, the initial account name will be checked (assuming
//   that the account in question is a legacy account or has changed its name at least once)
// - if no profile matches the specified name, nil will be returned instead
func (a *MojangAPI) GetId(name string, at time.Time) (*ProfileId, error) {
  res, err := a.execute("GET", fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s?at=%d", url.PathEscape(name), at.Unix()), nil)
  if err != nil {
    return nil, err
  }

  if res.StatusCode == 204 {
    return nil, nil
  }

  profile := &ProfileId{}
  defer res.Body.Close()
  err = json.NewDecoder(res.Body).Decode(profile)
  if err != nil {
    return nil, err
  }

  err = profile.init(at)
  if err != nil {
    return nil, err
  }
  return profile, nil
}

// resolves a list of multiple names at the current time
// only 100 names may be resolved at a time
func (a *MojangAPI) BulkGetId(names []string) ([]ProfileId, error) {
  if len(names) > 100 {
    return nil, errors.New("cannot request more than 100 names")
  }

  payload, err := json.Marshal(names)
  if err != nil {
    return nil, err
  }

  at := time.Now()
  res, err := a.execute("POST", "https://api.mojang.com/profiles/minecraft", bytes.NewBuffer(payload))
  if err != nil {
    return nil, err
  }

  if res.StatusCode == 204 { // TODO: verify whether this case actually occurs
    return make([]ProfileId, 0), nil
  }

  profiles := make([]ProfileId, 0)
  defer res.Body.Close()
  err = json.NewDecoder(res.Body).Decode(&profiles)
  if err != nil {
    return nil, err
  }

  for _, profile := range profiles {
    err = profile.init(at)
    if err != nil {
      return nil, err
    }
  }
  return profiles, nil
}

// retrieves the complete name change history for a given profile
// the initial account name is indicated by the lack of its timestamp (e.g. if set to UNIX epoch)
func (a *MojangAPI) GetHistory(id uuid.UUID) ([]NameChange, error) {
  res, err := a.execute("GET", fmt.Sprintf("https://api.mojang.com/user/profiles/%s/names", ToMojangId(id)), nil)
  if err != nil {
    return nil, err
  }

  if res.StatusCode == 204 { // TODO: Verify whether this case actually occurs (e.g. is the API consistent)
    return nil, nil
  }

  history := make([]NameChange, 0)
  defer res.Body.Close()
  err = json.NewDecoder(res.Body).Decode(&history)
  if err != nil {
    return nil, err
  }

  for _, change := range history {
    err = change.init()
    if err != nil {
      return nil, err
    }
  }
  return history, nil
}
