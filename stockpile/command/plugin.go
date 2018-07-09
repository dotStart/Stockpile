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

type PluginCommand struct {
  ClientCommand
}

func (*PluginCommand) Name() string {
  return "plugins"
}

func (*PluginCommand) Synopsis() string {
  return "displays a list of commands loaded on the Stockpile server"
}

func (*PluginCommand) Usage() string {
  return `Usage: stockpile plugins [options]

This command displays the list of loaded plugins on a given Stockpile server:

  $ stockpile plugins

Available command specific flags:

`
}

func (c *PluginCommand) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
  client, err := c.createClient()
  if err != nil {
    fmt.Fprintf(os.Stderr, "failed to establish a connection to server \"%s\": %s\n", c.flagServerAddress, err)
    return 1
  }

  systemService := rpc.NewSystemServiceClient(client)
  res, err := systemService.GetPlugins(ctx, &empty.Empty{})
  if err != nil {
    fmt.Fprintf(os.Stderr, "server responded with error: %s", err)
    return 1
  }

  pluginList := rpc.PluginMetadataListFromRpc(res)
  fmt.Fprintf(os.Stdout, "server has %d plugin(s) loaded:\n\n", len(pluginList))

  for _, plugin := range pluginList {
    writeTable(os.Stdout, *plugin)
    fmt.Fprintf(os.Stdout, "\n")
  }
  return 0
}
