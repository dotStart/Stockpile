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
  "net/http"

  "github.com/dotStart/Stockpile/stockpile/cache"
  "github.com/op/go-logging"
)

type Server struct {
  logger *logging.Logger
  cache  *cache.Cache
}

func NewServer(httpMux *http.ServeMux, cache *cache.Cache) (*Server) {
  srv := &Server{
    logger: logging.MustGetLogger("legacy"),
    cache:  cache,
  }

  httpMux.HandleFunc("/v1/shutdown", srv.handleServerShutdown)
  httpMux.HandleFunc("/v1/blacklist", srv.handleBlacklist)
  httpMux.HandleFunc("/v1/login", srv.handleLogin)
  httpMux.HandleFunc("/v1/name", srv.handleName)
  httpMux.HandleFunc("/v1/name/", srv.handleName)
  httpMux.HandleFunc("/v1/profile", srv.handleProfile)
  httpMux.HandleFunc("/v1/profile/", srv.handleProfile)
  httpMux.HandleFunc("/v1/", srv.handleRoot)

  return srv
}
