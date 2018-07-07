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
  "time"

  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/dotStart/Stockpile/stockpile/server/rpc"
  "github.com/google/subcommands"
  "golang.org/x/net/context"
)

type IdCommand struct {
  ClientCommand

  flagTimestamp string
}

func (*IdCommand) Name() string {
  return "get-id"
}

func (*IdCommand) Synopsis() string {
  return "queries a user's profile id from a Stockpile server"
}

func (*IdCommand) Usage() string {
  return `Usage: stockpile get-id [options] <name> [name2] [name3] ...

This command retrieves the profile identifier (and some additional data) from a Stockpile server:

  $ stockpile get-id dotStart

When unspecified, the names will be resolved at the current time. You may, however, also specify a
custom time to resolve a name at:

  $ stockpile get-id --time=2016-05-02T15:04:05Z07:00 dotStart

In addition, you may also retrieve the profile identifiers of multiple names at once:

  $ stockpile get-id dotStart MiniDigger

Note that the current time will always be substituted when multiple names are passed.

Available command specific flags:

`
}

func (c *IdCommand) SetFlags(f *flag.FlagSet) {
  c.ClientCommand.SetFlags(f)
  f.StringVar(&c.flagTimestamp, "time", "now", "defines the time at which this name should be resolved (for instance \""+time.RFC3339+"\"; may also be \"now\" for the current time or \"zero\" for the initial account name)")
}

func (c *IdCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
  client, err := c.createClient()
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to establish a connection to server \"%s\": %s\n", c.flagServerAddress, err)
    return 1
  }

  if f.NArg() == 0 {
    fmt.Fprintf(os.Stderr, "illegal command invocation: display-name is required\n")
    return 1
  }

  if f.NArg() == 1 {
    var timestamp time.Time
    if c.flagTimestamp == "now" {
      timestamp = time.Now()
    } else if c.flagTimestamp == "0" {
      timestamp = time.Unix(0, 0)
    } else {
      timestamp, err = time.Parse(time.RFC3339, c.flagTimestamp)
      if err != nil {
        fmt.Fprintf(os.Stderr, "illegal request time \"%s\": %s\n", c.flagTimestamp, err)
        return 1
      }
    }

    profileService := rpc.NewProfileServiceClient(client)
    res, err := profileService.GetId(ctx, &rpc.GetIdRequest{
      Name:      f.Arg(0),
      Timestamp: timestamp.Unix(),
    })
    if err != nil {
      fmt.Fprintf(os.Stderr, "command execution has failed: %s\n", err)
      return 1
    }

    if !res.IsPopulated() {
      fmt.Fprintf(os.Stderr, "no such profile\n")
      return 1
    }

    profile, err := rpc.ProfileIdFromRpc(res)
    if err != nil {
      fmt.Fprintf(os.Stderr, "failed to convert profile: %s", err)
      return 1
    }
    writeTable(os.Stdout, *profile)
    return 0
  }

  profileService := rpc.NewProfileServiceClient(client)
  res, err := profileService.BulkGetId(ctx, &rpc.BulkIdRequest{
    Names: f.Args(),
  })
  if err != nil {
    fmt.Fprintf(os.Stderr, "command execution has failed: %s\n", err)
    return 1
  }

  if !res.IsPopulated() {
    fmt.Fprintf(os.Stderr, "no such profile")
    return 1
  }

  profileIds, err := rpc.ProfileIdsFromRpcArray(res.Ids)
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to decode one or more profile ids: %s", err)
    return 1
  }

  writeTable(
    os.Stdout,
    struct {
      Ids []*mojang.ProfileId
    }{
      profileIds,
    },
  )
  return 0
}
