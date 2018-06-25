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
package command

import (
  "flag"
  "fmt"

  "github.com/dotStart/Stockpile/stockpile/server"
  "google.golang.org/grpc"
)

type ClientCommand struct {
  flagServerAddress string
}

// creates a new grpc client using the command configuration
func (c *ClientCommand) createClient() (*grpc.ClientConn, error) {
  client, err := grpc.Dial(c.flagServerAddress, grpc.WithInsecure())
  if err != nil {
    return nil, err
  }

  return client, nil
}

func (c *ClientCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&c.flagServerAddress, "server-address", fmt.Sprintf("%s:%d", server.DefaultAddress, server.DefaultPort), "specifies the address of the target server")
  // TODO: TLS
}
