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

  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/dotStart/Stockpile/stockpile/server/rpc"
  "github.com/google/subcommands"
  "golang.org/x/net/context"
)

type HistoryCommand struct {
  ClientCommand

  flagTimestamp string
}

func (*HistoryCommand) Name() string {
  return "name-history"
}

func (*HistoryCommand) Synopsis() string {
  return "queries a profile's name history"
}

func (*HistoryCommand) Usage() string {
  return `Usage: stockpile name-history [options] <id>

This command retrieves the name history of a profile from a Stockpile server:

  $ stockpile name-history d71a5dac-4e71-443b-8158-4389c269e44d

Note that Mojang formatted UUIDs are also supported:

  $ stockpile name-history d71a5dac4e71443b81584389c269e44d

Available command specific flags:

`
}

func (c *HistoryCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
  client, err := c.createClient()
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to establish a connection to server \"%s\": %s\n", c.flagServerAddress, err)
    return 1
  }

  if f.NArg() != 1 {
    fmt.Fprintf(os.Stderr, "illegal command invocation: id is required\n")
    return 1
  }

  id, err := mojang.ParseId(f.Arg(0))
  if err != nil {
    fmt.Fprintf(os.Stderr, "illegal profile id: %s (is this a UUID?)", id)
    return 1
  }

  profileService := rpc.NewProfileServiceClient(client)
  res, err := profileService.GetNameHistory(ctx, &rpc.IdRequest{
    Id: id.String(),
  })
  if err != nil {
    fmt.Fprintf(os.Stderr, "command execution has failed: %s\n", err)
    return 1
  }

  if !res.IsPopulated() {
    fmt.Fprintf(os.Stderr, "no such profile")
    return 1
  }

  history := rpc.NameHistoryFromRpc(res)
  writeTable(os.Stdout, *history)
  return 0
}
