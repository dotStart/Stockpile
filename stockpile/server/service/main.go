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
package service

//go:generate protoc -I ../rpc --go_out=plugins=grpc:../rpc common.proto events.proto profile.proto server.proto

import (
  "net"

  "github.com/dotStart/Stockpile/stockpile/cache"
  "github.com/dotStart/Stockpile/stockpile/server/rpc"
  "github.com/op/go-logging"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
)

// Represents an RPC server
type Server struct {
  logger *logging.Logger
  cache  *cache.Cache

  srv *grpc.Server
}

// Constructs a new RPC server instance and starts it
func NewServer(cache *cache.Cache) (*Server, error) {
  logger := logging.MustGetLogger("rpc")

  return &Server{
    logger: logger,
    cache:  cache,
  }, nil
}

// Starts listening on an arbitrary socket
func (s *Server) Listen(listener net.Listener) {
  s.srv = grpc.NewServer()
  grpc.NewServer()
  rpc.RegisterProfileServiceServer(s.srv, NewProfileService(s.cache))
  rpc.RegisterServerServiceServer(s.srv, NewServerService(s.cache))
  rpc.RegisterEventServiceServer(s.srv, NewEventService(s.cache))
  reflection.Register(s.srv)
  s.srv.Serve(listener)
}

// Stops listening
func (s *Server) Stop() {
  s.srv.Stop()
}

// destroys the server instance permanently
func (s *Server) Destroy() {
  s.srv.Stop()
  s.cache.Close()
}
