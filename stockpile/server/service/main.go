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

//go:generate protoc -I ../rpc --go_out=plugins=grpc:../rpc common.proto profile.proto server.proto

import (
  "fmt"
  "net"
  "path"

  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/dotStart/Stockpile/stockpile/plugin"
  "github.com/dotStart/Stockpile/stockpile/server"
  "github.com/dotStart/Stockpile/stockpile/server/rpc"
  "github.com/op/go-logging"
  "google.golang.org/grpc"
  "google.golang.org/grpc/reflection"
)

// Represents an RPC server
type Server struct {
  logger  *logging.Logger
  cfg     *server.Config
  storage plugin.StorageBackend

  srv *grpc.Server
}

// Constructs a new RPC server instance and starts it
func NewServer(config *server.Config) (*Server, error) {
  logger := logging.MustGetLogger("rpc")

  var storage plugin.StorageBackend
  if config.Storage.Type == "mem" {
    logger.Warning("Using in-memory storage")
    storage = plugin.NewMemoryStorageBackend(config)
  } else {
    plg, err := plugin.Open(path.Join("plugins", config.Storage.Type))
    if err != nil {
      return nil, err
    }

    if !plg.HasStorageBackendImplementation() {
      return nil, fmt.Errorf("selected plugin \"%s\" does not provide a storage backend implementation", config.Storage.Type)
    }

    storage, err = plg.CreateStorageBackend(config)
    if err != nil {
      return nil, err
    }

    logger.Infof("Using database plugin: %s", config.Storage.Type)
  }

  return &Server{
    logger:  logger,
    cfg:     config,
    storage: storage,
  }, nil
}

// Starts listening on an arbitrary socket
func (s *Server) Listen(listener net.Listener) {
  api := mojang.New()

  s.srv = grpc.NewServer()
  grpc.NewServer()
  rpc.RegisterProfileServiceServer(s.srv, NewProfileService(api, s.cfg, s.storage))
  rpc.RegisterServerServiceServer(s.srv, NewServerService(api, s.cfg, s.storage))
  reflection.Register(s.srv)
  s.srv.Serve(listener)
}

// Stops listening
func (s *Server) Stop() {
  s.srv.Stop()
}
