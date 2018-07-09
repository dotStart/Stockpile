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

  "github.com/dotStart/Stockpile/entity"
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

  stream, err := client.EventChannel(func(err error) {
    fmt.Fprintf(os.Stderr, "failed to poll events: %s", err)
    os.Exit(1)
  })

  fmt.Fprintf(os.Stderr, "listening for events:\n")
  for event := range stream {

    var entry string
    switch event.Type {
    case entity.ProfileIdEvent:
      profileId, _ := event.ProfileIdPayload()
      entry = fmt.Sprintf("updated name association for name \"%s\" to profile %s (valid from %s until %s)", profileId.Name, profileId.Id, profileId.FirstSeenAt, profileId.ValidUntil)
    case entity.NameHistoryEvent:
      id, _ := event.IdKey()
      entry = fmt.Sprintf("updated name history for profile %s", id)
    case entity.ProfileEvent:
      profile, _ := event.ProfilePayload()
      entry = fmt.Sprintf("updated profile %s (display name: \"%s\")", profile.Id, profile.Name)
    case entity.BlacklistEvent:
      entry = fmt.Sprintf("updated blacklist")
    default:
      entry = "Unknown Event"
    }

    fmt.Fprintf(os.Stdout, "[%s] %s\n", time.Now(), entry)
  }

  return 0
}
