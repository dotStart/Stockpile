---
title: Installation - Redis Plugin
layout: docs-navigation
---

# Installation

The redis plugin is part of the standard distribution for 64-Bit versions of
Linux and will be loaded automatically when present within the
`plugin-directory`.

Once loaded, you may specify the `redis` storage backend:

```
storage "redis" {
  address = "localhost:6379"
  db = 0
}
```

For more information, refer to the [Configuration](configuration.html) section.
