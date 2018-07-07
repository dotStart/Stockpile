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
package ui

import (
  "net/http"
  "time"

  "github.com/dotStart/Stockpile/stockpile/cache"
  "github.com/dotStart/Stockpile/stockpile/metadata"
  "github.com/googollee/go-socket.io"
  "github.com/op/go-logging"
)

//go:generate go-bindata-assetfs -pkg ui -o data.gen.go -prefix ../../ui/dist -ignore semantic.(js|css|min.js) ../../ui/dist/...

type Server struct {
  logger *logging.Logger
  io     *socketio.Server
  cache  *cache.Cache

  rateLimitTicker *time.Ticker

  corsOverride string
}

func NewServer(httpSrv *http.ServeMux, corsOverride string, cacheImpl *cache.Cache) (*Server, error) {
  io, err := socketio.NewServer(nil)
  if err != nil {
    return nil, err
  }

  srv := &Server{
    logger: logging.MustGetLogger("ui"),
    io:     io,
    cache:  cacheImpl,

    rateLimitTicker: time.NewTicker(time.Minute),

    corsOverride: corsOverride,
  }

  io.On("connection", srv.onSocketConnect)
  io.On("disconnection", srv.onSocketDisconnect)
  go srv.forwardRateLimit()
  go srv.forwardCacheEvents()

  httpSrv.HandleFunc("/ui/socket.io/", srv.handleSocket)
  httpSrv.Handle("/ui/", http.StripPrefix("/ui/", http.FileServer(assetFS())))
  httpSrv.HandleFunc("/", srv.handleRootRequest)
  return srv, nil
}

// handles requests to the root endpoint (either by redirecting to the UI endpoint or by responding
// with a 404)
func (s *Server) handleRootRequest(w http.ResponseWriter, req *http.Request) {
  if req.URL.Path != "/" {
    http.NotFound(w, req)
    return
  }

  http.Redirect(w, req, "/ui/", http.StatusFound)
}

func (s *Server) handleSocket(w http.ResponseWriter, req *http.Request) {
  if s.corsOverride != "" {
    w.Header().Set("Access-Control-Allow-Origin", s.corsOverride)
    w.Header().Set("Access-Control-Allow-Credentials", "true")
  }
  s.io.ServeHTTP(w, req)
}

// forwards the current rate limit to connected clients
func (s *Server) forwardRateLimit() {
  for range s.rateLimitTicker.C {
    s.io.BroadcastTo("ui", "rate-limit", s.cache.GetRateLimitAllocation())
  }
}

// forwards all cache events to connected clients
func (s *Server) forwardCacheEvents() {
  for e := range s.cache.Events {
    // TODO: socket.io eats serialization errors here - use this to debug until this issue is fixed
    /*_, err := json.Marshal(e)
    if err != nil {
      s.logger.Errorf("ENCODING ERROR: %s %v", err, e)
    } */
    s.logger.Debugf("forwarding event of type %T (using key %T) to active clients", e.Object, e.Key)
    s.io.BroadcastTo("ui", "cache", e)
  }
}

func (s *Server) onSocketConnect(io socketio.Socket) {
  s.logger.Debugf("client %s (id: %s) established websocket connection", io.Request().RemoteAddr, io.Id())
  io.Join("ui")
  io.Emit(
    "system",
    struct {
      Version string `json:"version"`
    }{
      Version: metadata.VersionFull(),
    },
  )
  io.Emit("rate-limit", s.cache.GetRateLimitAllocation())
}

func (s *Server) onSocketDisconnect(io socketio.Socket) {
  s.logger.Debugf("client %s (id %s) has disconnected", io.Request().RemoteAddr, io.Id())
}
