---
title: Configuration - Redis Plugin
layout: docs-navigation
---

# Configuration

The following configuration properties are currently available for the Redis
storage backend:

* **address**<br />
  *Default:* unset<br />
  *Required:* yes<br />
  Specifies the address on which the Redis server listens
* **password**<br />
  *Default:* unset<br />
  *Required:* no<br />
  Defines a password for authentication purposes
* **db**<br />
  *Default:* unset<br />
  *Required:* yes<br />
  Specifies the database to store cache entries in

## Example Configuration

```
storage "redis" {
  address = "127.0.0.1:6379"
  password = "stockpilerulez"
  db = 1
}
```
