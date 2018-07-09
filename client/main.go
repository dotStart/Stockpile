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
package client

import (
  "log"
  "os"

  "github.com/dotStart/Stockpile/rpc"
  "google.golang.org/grpc"
)

type Stockpile struct {
  Logger *log.Logger
  client *grpc.ClientConn

  eventService   rpc.EventServiceClient
  profileService rpc.ProfileServiceClient
  serverService  rpc.ServerServiceClient
  systemService  rpc.SystemServiceClient
}

// creates a new client for the specified server address
func New(address string) (*Stockpile, error) {
  client, err := grpc.Dial(address, grpc.WithInsecure()) // TODO: SSL support
  if err != nil {
    return nil, err
  }

  return &Stockpile{
    client:         client,
    Logger:         log.New(os.Stderr, "[stockpile]", log.Flags()),
    eventService:   rpc.NewEventServiceClient(client),
    profileService: rpc.NewProfileServiceClient(client),
    serverService:  rpc.NewServerServiceClient(client),
    systemService:  rpc.NewSystemServiceClient(client),
  }, nil
}

// closes the connection to the remote server
func (s *Stockpile) Close() error {
  return s.client.Close()
}
