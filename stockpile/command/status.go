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
  "os"

  "github.com/dotStart/Stockpile/rpc"
  "github.com/golang/protobuf/ptypes/empty"
  "github.com/google/subcommands"
  "golang.org/x/net/context"
)

type StatusCommand struct {
  ClientCommand
}

func (*StatusCommand) Name() string {
  return "status"
}

func (*StatusCommand) Synopsis() string {
  return "displays the current status of a Stockpile server"
}

func (*StatusCommand) Usage() string {
  return `Usage: stockpile status [options]

This command displays basic status information of a given Stockpile server:

  $ stockpile status

Available command specific flags:

`
}

func (c *StatusCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
  client, err := c.createClient()
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to establish a connection to server \"%s\": %s\n", c.flagServerAddress, err)
    return 1
  }

  systemService := rpc.NewSystemServiceClient(client)
  status, err := systemService.GetStatus(ctx, &empty.Empty{})
  if err != nil {
    fmt.Fprintf(os.Stderr, "server responded with error: %s", err)
    return 1
  }

  writeTable(os.Stdout, *status)
  return 0
}
