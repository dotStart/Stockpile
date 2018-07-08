---
title: Storage Backends - Getting Started
layout: docs-navigation
---

# Storage Backends

Stockpile provides support for various storage backends (typically via the
plugin API).

As part of its core, Stockpile provides the following storage backends:

* `mem` - memory based storage system
  All data in this storage backend will automatically get cleared when the
  application is stopped. As such, this storage backend should only ever be used
  in test environments. Default in development mode.
* `file` file based storage system
  Cache entries will be permanently stored in a directory structure (serialized
  as json). Recommended for small production deployments.

In addition, we maintain a set of plugins which provide support for third party
storage systems:

* `redis` - Redis KV store support
  Permanently stores cache entries on a Redis server inside of your network.

Note that plugin based storage backends are currently only available on 64-Bit
versions of Linux due to restrictions in the application runtime.

## Configuration

Only one storage backend is enabled at a time. They are defined via the
`storage` property block.

For memory storage:

```
storage "mem" {
}
```

For file based storage:

```
storage "file" {
  path = "data"
}
```

For redis:

```
storage "redis" {
  address = "localhost:6379"
  password = "admin1234"
  database = 0
}
```

## Next Steps

This concludes the getting started guide. For more information refer to the
remaining documentation pages.
