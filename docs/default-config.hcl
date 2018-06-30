bind-address = "0.0.0.0:36623"
legacy-api = false

// no storage backend in default - required for actual operation
// example:
// storage "mem" {
//   param1 = false
//   param2 = "hostname:port"
// }

ttl {
  name = "888h"
  name-history = "180h"
  profile = "168h"
  blacklist = "168h"
}
