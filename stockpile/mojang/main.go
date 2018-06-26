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
  "io"
  "net/http"
  "runtime"

  "github.com/dotStart/Stockpile/stockpile/metadata"
  "github.com/op/go-logging"
)

type MojangAPI struct {
  logger *logging.Logger
  http   *http.Client
}

// Creates a new Mojang API client
func New() *MojangAPI {
  return &MojangAPI{
    logger: logging.MustGetLogger("api"),
    http:   &http.Client{},
  }
}

// Executes an HTTP request
func (a *MojangAPI) execute(method string, uri string, body io.Reader) (*http.Response, error) {
  req, err := http.NewRequest(method, uri, body)
  if err != nil {
    return nil, err
  }

  req.Header.Set("user-agent", fmt.Sprintf("Stockpile/%s (Go/%s; %s; +https://github.com/dotStart/Stockpile)", metadata.VersionFull(), runtime.Version(), metadata.Brand()))
  req.Header.Set("content-type", "application/json")

  a.logger.Debugf("Sending request: %s %s", method, uri)
  res, err := a.http.Do(req)
  if err != nil {
    return nil, err
  }

  statusCategory := res.StatusCode / 100
  a.logger.Debugf("Server responded with status code %d (category %d)", res.StatusCode, statusCategory)

  if statusCategory == 2 || res.StatusCode == 404 {
    return res, nil
  }
  if statusCategory == 4 {
    return nil, fmt.Errorf("client error (code %d): %s", res.StatusCode, uri)
  }
  if statusCategory == 5 {
    return nil, fmt.Errorf("server error (code %d): %s", res.StatusCode, uri)
  }
  return nil, fmt.Errorf("unknown error (code %d): %s", res.StatusCode, uri)
}
