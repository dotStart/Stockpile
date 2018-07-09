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
package legacy

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "time"

  "github.com/dotStart/Stockpile/entity"
  "github.com/google/uuid"
)

func (s *Server) handleBlacklist(w http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" || req.URL.Path != "/v1/blacklist" {
    http.NotFound(w, req)
    return
  }

  defer req.Body.Close()
  enc, err := ioutil.ReadAll(req.Body)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to read hostname from request body: %s", err), http.StatusBadRequest)
    return
  }

  hostname := string(enc)
  blacklist, err := s.cache.GetBlacklist()
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to retrieve blacklist: %s", err), http.StatusServiceUnavailable)
    return
  }
  blacklisted, err := blacklist.IsBlacklisted(hostname)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to check blacklist: %s", err), http.StatusServiceUnavailable)
    return
  }

  result := struct {
    Blacklisted bool   `json:"blacklisted"`
    Hostname    string `json:"hostname"`
  }{
    Blacklisted: blacklisted,
    Hostname:    hostname,
  }
  j, err := json.Marshal(result)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to encode result: %s", err), http.StatusServiceUnavailable)
  }

  w.Write([]byte(j))
}

func (s *Server) handleLogin(w http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" || req.URL.Path != "/v1/login" {
    http.NotFound(w, req)
    return
  }

  query := req.URL.Query()
  if query.Get("username") == "" {
    http.Error(w, "illegal login request: username parameter is required", http.StatusBadRequest)
    return
  }
  if query.Get("serverId") == "" {
    http.Error(w, "illegal login request: serverId parameter is required", http.StatusBadRequest)
    return
  }

  profile, err := s.cache.Login(query.Get("username"), query.Get("serverId"), "") // ip not supported in legacy
  if err != nil {
    http.Error(w, fmt.Sprintf("login server responded with error: %s", err), http.StatusServiceUnavailable)
    return
  }

  if req.Header.Get("X-Forward") != "" {
    j, err := profile.Mojang()
    if err != nil {
      http.Error(w, fmt.Sprintf("failed to encode response: %s", err), http.StatusServiceUnavailable)
      return
    }
    w.Write(j)
    return
  }

  props := make([]*entity.ProfileProperty, 0)
  for _, prop := range profile.Properties {
    props = append(props, prop)
  }

  enc := struct {
    CacheTimestamp int64                     `json:"cacheTimestamp"`
    Identifier     uuid.UUID                 `json:"identifier"`
    Name           string                    `json:"name"`
    Properties     []*entity.ProfileProperty `json:"properties"`
  }{
    CacheTimestamp: time.Now().Unix(), // unsupported in this environment, substitute a valid value
    Identifier:     profile.Id,
    Name:           profile.Name,
    Properties:     props,
  }

  j, err := json.Marshal(enc)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to encode response: %s", err), http.StatusServiceUnavailable)
    return
  }
  w.Write(j)
}
