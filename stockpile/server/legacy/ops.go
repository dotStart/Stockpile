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
  "net/http"
  "os"
  "time"

  "github.com/dotStart/Stockpile/stockpile/metadata"
)

func (s *Server) handleRoot(w http.ResponseWriter, req *http.Request) {
  if req.Method != "GET" || (req.URL.Path != "/v1" && req.URL.Path != "/v1/") {
    http.NotFound(w, req)
    return
  }

  api := struct {
    State   string `json:"state"`
    Version uint8  `json:"version"`
  }{
    State:   "DEPRECATED",
    Version: 1,
  }

  enc := struct {
    Api     interface{} `json:"api"`
    Name    string      `json:"name"`
    Vendor  string      `json:"vendor"`
    Version string      `json:"version"`
  }{
    Api:     api,
    Name:    "Stockpile",
    Vendor:  "Torchmind",
    Version: metadata.VersionFull(),
  }

  j, err := json.Marshal(enc)
  if err != nil {
    http.Error(w, fmt.Sprintf("failed to encode response: %s", err), http.StatusServiceUnavailable)
    return
  }
  w.Write(j)
}

func (s *Server) handleServerShutdown(w http.ResponseWriter, req *http.Request) {
  if req.Method != "POST" || req.URL.Path != "/v1/shutdown" {
    http.NotFound(w, req)
    return
  }

  w.WriteHeader(http.StatusNoContent)
  go func() {
    s.logger.Infof("Graceful shutdown has been requested via legacy API")
    time.Sleep(time.Second * 5)
    s.logger.Infof("Performing shutdown")
    os.Exit(0)
  }()
}
