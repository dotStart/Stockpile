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
  "io/ioutil"
  "net/url"
  "strings"

  "github.com/dotStart/Stockpile/entity"
  "github.com/op/go-logging"
)

var blacklistLogger = logging.MustGetLogger("blacklist")

// retrieves the server blacklist
func (a *MojangAPI) GetBlacklist() (*entity.Blacklist, error) {
  res, err := a.execute("GET", "https://sessionserver.mojang.com/blockedservers", nil)
  if err != nil {
    return nil, err
  }

  if res.StatusCode == 204 {
    return entity.NewBlacklist(make([]string, 0))
  }

  defer res.Body.Close()
  encoded, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil, err
  }

  blacklistStr := string(encoded)
  blacklistStr = strings.Replace(blacklistStr, "\r", "", -1)

  return entity.NewBlacklist(strings.Split(blacklistStr, "\n"))
}

// performs the server-side phase of the online handshake
func (a *MojangAPI) Login(displayName string, serverId string, ip string) (*entity.Profile, error) {
  if ip != "" {
    ip = "&ip=" + url.QueryEscape(ip)
  }
  res, err := a.execute("GET", fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/hasJoined?username=%s&serverId=%s%s", url.QueryEscape(displayName), url.QueryEscape(serverId), ip), nil)
  if err != nil {
    return nil, err
  }

  profile := &entity.Profile{}
  defer res.Body.Close()
  err = profile.Read(res.Body)
  if err != nil {
    return nil, err
  }

  return profile, nil
}
