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

  "github.com/dotStart/Stockpile/entity"
  "github.com/google/uuid"
)

// retrieves the profile id (and some associated attributes) for a given display name at the specified time
// - if the UNIX epoch (e.g. zero) is passed instead of a real time, the initial account name will be checked (assuming
//   that the account in question is a legacy account or has changed its name at least once)
// - if no profile matches the specified name, nil will be returned instead
func (a *MojangAPI) GetId(name string, at time.Time) (*entity.ProfileId, error) {
  res, err := a.execute("GET", fmt.Sprintf("https://api.mojang.com/users/profiles/minecraft/%s?at=%d", url.PathEscape(name), at.Unix()), nil)
  if err != nil {
    return nil, err
  }

  if res.StatusCode == 204 {
    a.logger.Debugf("server reported no association for name \"%s\" at time %s", name, at)
    return nil, nil
  }

  profile := &entity.ProfileId{}
  defer res.Body.Close()
  err = profile.Read(res.Body, at)
  return profile, err
}

// resolves a list of multiple names at the current time
// only 100 names may be resolved at a time
func (a *MojangAPI) BulkGetId(names []string) ([]*entity.ProfileId, error) {
  if len(names) > 100 {
    return nil, errors.New("cannot request more than 100 names")
  }

  payload, err := json.Marshal(names)
  if err != nil {
    return nil, err
  }

  res, err := a.execute("POST", "https://api.mojang.com/profiles/minecraft", bytes.NewBuffer(payload))
  if err != nil {
    return nil, err
  }

  if res.StatusCode == 204 { // TODO: verify whether this case actually occurs
    return make([]*entity.ProfileId, 0), nil
  }

  defer res.Body.Close()
  return entity.ReadProfileIdArray(res.Body)
}

// retrieves the complete name change history for a given profile
// the initial account name is indicated by the lack of its timestamp (e.g. if set to UNIX epoch)
func (a *MojangAPI) GetHistory(id uuid.UUID) (*entity.NameChangeHistory, error) {
  res, err := a.execute("GET", fmt.Sprintf("https://api.mojang.com/user/profiles/%s/names", entity.ToMojangId(id)), nil)
  if err != nil {
    return nil, err
  }

  if res.StatusCode == 204 { // TODO: Verify whether this case actually occurs (e.g. is the API consistent)
    return nil, nil
  }

  history := &entity.NameChangeHistory{}
  defer res.Body.Close()
  err = history.Read(res.Body)
  if err != nil {
    return nil, err
  }

  return history, nil
}
