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
  "context"
  "flag"
  "fmt"
  "net"
  "net/http"
  "os"

  "github.com/dotStart/Stockpile/stockpile/cache"
  "github.com/dotStart/Stockpile/stockpile/metadata"
  "github.com/dotStart/Stockpile/stockpile/mojang"
  "github.com/dotStart/Stockpile/stockpile/plugin"
  "github.com/dotStart/Stockpile/stockpile/server"
  "github.com/dotStart/Stockpile/stockpile/server/service"
  "github.com/dotStart/Stockpile/stockpile/server/ui"
  "github.com/google/subcommands"
  "github.com/op/go-logging"
  "github.com/soheilhy/cmux"
)

type ServerCommand struct {
  flagConfig       string
  flagDevelopment  bool
  flagLogLevel     string
  flagCorsOverride string
}

func (*ServerCommand) Name() string {
  return "server"
}

func (*ServerCommand) Synopsis() string {
  return "starts a new Stockpile server instance"
}

func (*ServerCommand) Usage() string {
  return `Usage: stockpile server [options]

This command starts a new Stockpile server instance which relies to API requests. By default, Stockpile
will only start the RPC server.

Start a server with a configuration file:

  $ stockpile server -config=/etc/stockpile/config.hcl

Run in development mode:

  $ stockpile server -dev

For a full list of examples, please refer to the documentation.

Available command specific flags:

`
}

func (c *ServerCommand) SetFlags(f *flag.FlagSet) {
  f.StringVar(&c.flagConfig, "config", "", "specifies a configuration file or directory")
  f.BoolVar(&c.flagDevelopment, "dev", false, "enables development mode")
  f.StringVar(&c.flagLogLevel, "log-level", "info", "specifies a log level")
  f.StringVar(&c.flagCorsOverride, "cors-override", "", "specifies a host from which CORS requests are permitted")
}

func (c *ServerCommand) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
  var cfg *server.Config
  if c.flagDevelopment {
    cfg = server.DevelopmentConfig()
  }
  if c.flagConfig != "" {
    fileCfg, err := server.LoadConfig(c.flagConfig)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Error: Failed to read the configuration file(s): %s", err)
      return 1
    }

    if cfg == nil {
      cfg = fileCfg
    } else {
      cfg.Merge(fileCfg)
    }
  }

  level, err := logging.LogLevel(c.flagLogLevel)
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: Illegal log level \"%s\": %s", c.flagLogLevel, err)
    return 1
  }

  format := logging.MustStringFormatter(`%{color}%{time:15:04:05.000} [%{level:.4s}] %{module} : %{color:reset} %{message}`)
  backend := logging.AddModuleLevel(logging.NewBackendFormatter(logging.NewLogBackend(os.Stdout, "", 0), format))
  backend.SetLevel(level, "")
  logging.SetBackend(backend)

  fmt.Printf("==> Stockpile Configuration\n\n")
  fmt.Printf("   Server Address: %s\n", *cfg.BindAddress)
  fmt.Printf("          Version: %s\n", metadata.VersionFull())
  fmt.Printf("      Commit Hash: %s\n", metadata.CommitHash())
  fmt.Printf("        Log Level: %s\n", c.flagLogLevel)
  fmt.Printf("  Storage Backend: %s\n", cfg.Storage.Type)
  fmt.Printf("              PID: %d\n\n", os.Getpid())

  fmt.Printf("==> TTL Configuration\n\n")
  fmt.Printf("           Names: %s\n", cfg.Ttl.Name)
  fmt.Printf("         Profile: %s\n", cfg.Ttl.Profile)
  fmt.Printf("    Name History: %s\n", cfg.Ttl.NameHistory)
  fmt.Printf("       Blacklist: %s\n\n", cfg.Ttl.Profile)

  var log = logging.MustGetLogger("stockpile")

  if c.flagDevelopment {
    log.Warning("Stockpile is running in development mode")
  }

  // initialize the shared network listener and mux first so we can detect potential binding errors early on
  listener, err := net.Listen("tcp", *cfg.BindAddress)
  if err != nil {
    log.Fatalf("Failed to listen on %s (TCP): %s", *cfg.BindAddress, err)
  }

  mux := cmux.New(listener)

  // initialize the plugin system and cache manager
  pluginManager := plugin.NewManager(*cfg.PluginDir)
  pluginManager.LoadAll()

  storageFactory := pluginManager.Context.GetStorageBackend(cfg.Storage.Type)
  if storageFactory == nil {
    log.Fatalf("no such storage backend: %s", cfg.Storage.Type)
  }
  storage, err := storageFactory(cfg)
  if err != nil {
    log.Fatalf("failed to initialize storage backend \"%s\": %s", err)
  }
  log.Infof("Using database plugin: %s", cfg.Storage.Type)
  cacheImpl := cache.New(mojang.New(), storage)

  // initialize the RPC server at all times (only differ between mux policies depending on whether the legacy API or UI
  // is enabled)
  rpcPolicy := cmux.Any()
  if *cfg.UiEnabled {
    rpcPolicy = cmux.HTTP2HeaderField("content-type", "application/grpc")
  }
  rpcServer, err := service.NewServer(pluginManager, cacheImpl)
  if err != nil {
    log.Fatalf("Failed to initialize grpc server: %s", err)
  }
  go rpcServer.Listen(mux.Match(rpcPolicy))
  defer rpcServer.Destroy()
  log.Info("Enabled grpc server")

  if *cfg.UiEnabled {
    httpMux := http.NewServeMux()

    if c.flagCorsOverride != "" {
      log.Warningf("CORS override configured: %s", c.flagCorsOverride)
    }

    // instance currently unused
    ui.NewServer(httpMux, c.flagCorsOverride, cacheImpl)

    httpSrv := &http.Server{
      Handler: httpMux,
    }
    go httpSrv.Serve(mux.Match(cmux.Any()))
    log.Info("Enabled ui")
  }

  mux.Serve()
  return 0
}
