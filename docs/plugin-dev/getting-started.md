---
title: Getting Started - Plugin Development
layout: docs-navigation
---

# Getting Started

Stockpile plugins are written in [go](https://golang.org/) and use the built-in
plugin API (note that this feature is currently only available on Linux and Mac
OS). It is highly recommended to use the included `Vagrantfile` when developing
plugins on Windows or other operating systems.

## Metadata

While plugins may define an arbitrary set of symbols, they are required to
define two special symbols. The first of these is the `metadata` variable which
exposes basic human readable information about your plugin.

The structure is formatted like this:

```go
package main

import (
  "github.com/dotStart/Stockpile/stockpile/plugin"
)

var Metadata = plugin.Metadata{
  Name:    "My Plugin",
  Version: "0.1.0",
  Authors: []string{"Jane Doe"},
  Website: "https://www.example.org/my-plugin",
}
```

In short, the following properties are currently available:

* **Name**<br />
  Defines a human readable name which is used to identify your plugin within the
  web UI and client commands
* **Version**<br />
  Specifies a revision (we recommend formatting all version numbers using the
  [Semantic Versioning](https://semver.org/) format)
* **Authors**<br />
  Identifies one or more authors (either by real name or alias) who were
  involved with the development of the plugin
* **Website**<br />
  Specifies the website from which this plugin is distributed (this site
  typically also provides detailed documentation on the plugin)

## Initializer Function

In addition to the `metadata` variable, plugins are also required to define the
`InitializePlugin` function:

```go
func InitializePlugin(ctx *plugin.Context) error {
  ctx.RegisterStorageBackend("mybackend", NewCustomStorageBackend)
  return nil
}
```

The respective components of your plugin are registered with the plugin context
using this function. The following implementation types may be registered at the
moment:

* **Storage Backend**<br />
  *Registration Function:* `RegisterStorageBackend`<br />
  *Implementation Type:* `github.com/dotStart/Stockpile/stockpile/database.StorageBackend`<br />
  Implementations which handle the permanent storage and retrieval of cache
  entries

## Building

To produce a plugin file which may be loaded by Stockpile, you will need to
compile your package(s) using the following command:
`go build  -buildmode=plugin -o myplugin.so`

## Next Steps

You may continue with the [Storage Backend](storage-backend.html) section.
