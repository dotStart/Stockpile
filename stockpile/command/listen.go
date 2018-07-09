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
  "io"
  "os"
  "time"

  "github.com/dotStart/Stockpile/stockpile/cache"
  "github.com/dotStart/Stockpile/rpc"
  "github.com/golang/protobuf/ptypes/empty"
  "github.com/google/subcommands"
  "golang.org/x/net/context"
)

type ListenCommand struct {
  ClientCommand
}

func (*ListenCommand) Name() string {
  return "listen"
}

func (*ListenCommand) Synopsis() string {
  return "listens for events on a Stockpile server"
}

func (*ListenCommand) Usage() string {
  return `Usage: stockpile listen [options]

This command listens to all record changes processed by a given Stockpile server:

  $ stockpile listen

Available command specific flags:

`
}

func (c *ListenCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
  client, err := c.createClient()
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to establish a connection to server \"%s\": %s\n", c.flagServerAddress, err)
    return 1
  }

  eventService := rpc.NewEventServiceClient(client)
  stream, err := eventService.StreamEvents(ctx, &empty.Empty{})
  if err != nil {
    fmt.Fprintf(os.Stderr, "server responded with error: %s", err)
    return 1
  }

  fmt.Fprintf(os.Stderr, "listening for events:\n")
  for true {
    event, err := stream.Recv()
    if err != nil {
      if err == io.EOF {
        fmt.Fprintf(os.Stderr, "--- END OF STREAM ---")
        break
      }

      fmt.Fprintf(os.Stderr, "failed to poll event from server: %s", err)
      return 1
    }

    decoded, err := rpc.EventFromRpc(event)
    if err != nil {
      fmt.Fprintf(os.Stderr, "failed to convert event: %s\n", err)
    }

    var entry string
    switch decoded.Type {
    case cache.ProfileIdEvent:
      profileId, _ := decoded.ProfileIdPayload()
      entry = fmt.Sprintf("updated name association for name \"%s\" to profile %s (valid from %s until %s)", profileId.Name, profileId.Id, profileId.FirstSeenAt, profileId.ValidUntil)
    case cache.NameHistoryEvent:
      id, _ := decoded.IdKey()
      entry = fmt.Sprintf("updated name history for profile %s", id)
    case cache.ProfileEvent:
      profile, _ := decoded.ProfilePayload()
      entry = fmt.Sprintf("updated profile %s (display name: \"%s\")", profile.Id, profile.Name)
    case cache.BlacklistEvent:
      entry = fmt.Sprintf("updated blacklist")
    default:
      entry = "Unknown Event"
    }

    fmt.Fprintf(os.Stdout, "[%s] %s\n", time.Now(), entry)
  }

  return 0
}
