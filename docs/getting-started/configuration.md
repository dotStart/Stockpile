---
title: Configuration - Getting Started
layout: docs-navigation
---

# Configuration

The application provides various configuration options which provide further
customization options. Note that the application will automatically choose a
suitable configuration in development mode while production mode requires an
actual configuration file to be present.

Custom configuration files are written in the `hcl` format and typically loaded
using the `-config=<file or directory>` flag. You may either specify a single
file (for instance: `production.hcl`) or an entire directory of configuration
files. When a directory is specified, the files will be loaded in alphabetical
order (previously defined properties will be overridden).

## Supported Properties

* **plugin-dir**<br />
  *Default:* `plugins`<br />
  Defines the location from which plugins are loaded (only considered on 64-Bit
  versions of Linux)
* **bind-address**<br />
  *Default:* `127.0.0.1:36623`<br />
  Specifies the ip address and port on which the server listens for gRPC and
  REST requests
* **ui**<br />
  *Default:* `false` (`true` in development mode)<br />
  En- or disables the web interface
* **legacy-api**<br />
  *Default:* `false` (`true` in development mode)<br />
  En- or disables support for the legacy (Stockpile v1.0) REST API
* **storage**<br />
  *Default:* `mem` (unset in production mode)<br />
  Refer to the [Storage Backend](storage-backend.html) page for more information
* **ttl**<br />
  Refer to the section below for more information

### TTL Configuration Properties

* **name**<br />
  *Default:* `888h`<br />
  Specifies how long name <-> profile associations are to be persisted
* **name-history**<br />
  *Default:* `180h`<br />
  Specifies how long name histories are retained
* **profile**<br />
  *Default:* `168h`<br />
  Defines how long profiles are retained
* **blacklist**<br />
  *Default:* `168h`<br />
  Defines how long the server blacklist is persisted

## Example Configuration

```
plugin-dir = "plugins"
bind-address = "127.0.0.1:36623"
ui = true
legacy-api = true

storage "mem" {
}

ttl {
  name = "888h"
  name-history = "180h"
  profile = "168h"
  blacklist = "168h"
}
```

## Next Steps

You may now proceed to the [Storage Backend](storage-backend.html) section.
