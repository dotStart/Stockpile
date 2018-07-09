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
  "github.com/dotStart/Stockpile/rpc"
  "github.com/google/subcommands"
  "golang.org/x/net/context"
)

type ProfileCommand struct {
  ClientCommand
}

func (*ProfileCommand) Name() string {
  return "get-profile"
}

func (*ProfileCommand) Synopsis() string {
  return "queries a user's profile from a Stockpile server"
}

func (*ProfileCommand) Usage() string {
  return `Usage: stockpile profile [options] <id>

This command retrieves a profile from a Stockpile server:

  $ stockpile get-profile d71a5dac-4e71-443b-8158-4389c269e44d

Available command specific flags:

`
}

func (c *ProfileCommand) SetFlags(f *flag.FlagSet) {
  c.ClientCommand.SetFlags(f)
}

func (c *ProfileCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
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
  res, err := profileService.GetProfile(ctx, &rpc.IdRequest{
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

  profile, err := rpc.ProfileFromRpc(res)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to decode profile: %s", err)
    return 1
  }

  writeTable(os.Stdout, *profile)
  return 0
}
