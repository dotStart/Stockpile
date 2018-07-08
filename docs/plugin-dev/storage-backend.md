---
title: Storage Backend - Plugin Development
layout: docs-navigation
---

# Storage Backend

Storage backends are registered with the plugin context using the
`RegisterStorageBackend` function. Each storage backend requires a factory
method which takes a set of arbitrary user properties and converts them into an
implementation of
`github.com/dotStart/Stockpile/stockpile/database.StorageBackend`.

For instance:

```go
package main

import (
  "github.com/dotStart/Stockpile/stockpile/storage"
  "github.com/hashicorp/hcl2/gohcl"
  "github.com/op/go-logging"
)

type myStorageBackend struct {
  logger *logging.Logger
  cfg    *MyStorageBackendConfig
}

type MyStorageBackendConfig struct {
  myparam string `hcl:"my-param,attr"`
}

func NewRedisStorageBackend(cfg *server.Config) (storage.StorageBackend, error) {
  myCfg := &MyStorageBackendConfig{}
  diag := gohcl.DecodeBody(cfg.Storage.Parameters, nil, myCfg)
  if diag.HasErrors() {
    return nil, fmt.Errorf("illegal backend configuration: %s", diag.Error())
  }

  return &myStorageBackend{
    logger: logging.MustGetLogger("my-backend"),
    cfg: myCfg,    
  }, nil
}

// ...
```

For more information on the methods within the `StorageBackend` interface,
please refer to the
`github.com/dotStart/Stockpile/stockpile/database.StorageBackend` definition
within the source distribution.
