bind-address = "127.0.0.1:36623"
legacy-api = true

storage "mem" {
  param1 = false
  param2 = "hostname:port"
}

ttl {
  name = "888h"
  name-history = "180h"
  profile = "168h"
  blacklist = "168h"
}
