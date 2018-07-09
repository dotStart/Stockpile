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
  "strings"
  "time"

  "github.com/dotStart/Stockpile/entity"
  "github.com/google/uuid"
)

func (s *Server) handleName(w http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" && req.Method != "POST" && req.Method != "DELETE" {
    http.NotFound(w, req)
    return
  }

  var query string
  if req.Method == "GET" || (req.Method == "DELETE" && req.URL.Path != "/v1/name") {
    if !strings.HasPrefix(req.URL.Path, "/v1/name/") {
      http.NotFound(w, req)
      return
    }

    query = strings.TrimLeft(req.URL.Path, "/v1/name/")
    if strings.Contains(query, "/") || query == "" {
      http.NotFound(w, req)
      return
    }
  } else {
    defer req.Body.Close()
    data, err := ioutil.ReadAll(req.Body)
    if err != nil {
      http.Error(w, fmt.Sprintf("failed to read input: %s", err), http.StatusBadRequest)
      return
    }

    query = string(data)
  }

  at := time.Now()
  if req.Method == "DELETE" {
    err := s.cache.PurgeProfileId(query, at)
    if err != nil {
      http.Error(w, fmt.Sprintf("failed to purge profile association: %s", err), http.StatusServiceUnavailable)
      return
    }
    w.WriteHeader(http.StatusNoContent)
    return
  }

  profileId, err := s.cache.GetProfileId(query, time.Now())
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to retrieve profile association: %s", err), http.StatusServiceUnavailable)
    return
  }

  if profileId == nil {
    http.NotFound(w, req)
    return
  }

  enc := struct {
    Id uuid.UUID `json:"identifier"`
  }{
    Id: profileId.Id,
  }

  j, err := json.Marshal(enc)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to encode response: %s", err), http.StatusServiceUnavailable)
    return
  }
  w.Write(j)
}

func (s *Server) handleProfile(w http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" && req.Method != "POST" && req.Method != "DELETE" {
    http.NotFound(w, req)
    return
  }

  var query string
  if req.Method == "GET" || (req.Method == "DELETE" && req.URL.Path != "/v1/profile") {
    if !strings.HasPrefix(req.URL.Path, "/v1/profile/") {
      http.NotFound(w, req)
      return
    }

    query = strings.TrimLeft(req.URL.Path, "/v1/profile/")
    if strings.Contains(query, "/") {
      http.NotFound(w, req)
      return
    }
  } else {
    defer req.Body.Close()
    data, err := ioutil.ReadAll(req.Body)
    if err != nil {
      http.Error(w, fmt.Sprintf("failed to read input: %s", err), http.StatusBadRequest)
      return
    }

    query = string(data)
  }

  var err error
  var id uuid.UUID
  if entity.IsProfileId(query) {
    id, err = entity.ParseId(query)
    if err != nil {
      http.Error(w, fmt.Sprintf("failed to decode id: %s", err), http.StatusBadRequest)
      return
    }
  } else {
    profileId, err := s.cache.GetProfileId(query, time.Now())
    if err != nil {
      http.Error(w, fmt.Sprintf("failed to retrieve profile association: %s", err), http.StatusServiceUnavailable)
      return
    }

    id = profileId.Id
  }

  if req.Method == "DELETE" {
    err := s.cache.PurgeProfile(id)
    if err != nil {
      http.Error(w, fmt.Sprintf("failed to purge profile: %s", err), http.StatusServiceUnavailable)
      return
    }
    w.WriteHeader(http.StatusNoContent)
    return
  }

  profile, err := s.cache.GetProfile(id)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to retrieve profile: %s", err), http.StatusServiceUnavailable)
    return
  }

  if profile == nil {
    http.NotFound(w, req)
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
