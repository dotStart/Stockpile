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

  "github.com/dotStart/Stockpile/stockpile/server/rpc"
  "github.com/google/subcommands"
  "golang.org/x/net/context"
)

type BlacklistCommand struct {
  ClientCommand
}

func (*BlacklistCommand) Name() string {
  return "check-blacklist"
}

func (*BlacklistCommand) Synopsis() string {
  return "checks whether an address has been blacklisted using a remote Stockpile server"
}

func (*BlacklistCommand) Usage() string {
  return `Usage: stockpile check-blacklist <address> [address2] [address3] ...

This command evaluates whether a given server address has been blacklisted:

  $ stockpile check-blacklist example.org

You may also pass multiple addresses at once, if desired:

  $ stockpile get-id example.org example.net

Available command specific flags:

`
}

func (c *BlacklistCommand) SetFlags(f *flag.FlagSet) {
  c.ClientCommand.SetFlags(f)
}

func (c *BlacklistCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
  client, err := c.createClient()
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to establish a connection to server \"%s\": %s\n", c.flagServerAddress, err)
    return 1
  }

  if f.NArg() == 0 {
    fmt.Fprintf(os.Stderr, "illegal command invocation: address is required\n")
    return 1
  }

  serverService := rpc.NewServerServiceClient(client)
  res, err := serverService.CheckBlacklist(ctx, &rpc.CheckBlacklistRequest{
    Addresses: f.Args(),
  })
  if err != nil {
    fmt.Fprintf(os.Stderr, "command execution has failed: %s\n", err)
    return 1
  }

  if len(res.MatchedAddresses) == 0 {
    fmt.Fprintf(os.Stderr, "none of the passed addresses match")
    return 0
  }

  for _, match := range res.MatchedAddresses {
    fmt.Printf("%s matches the blacklist\n", match)
  }
  return 0
}
