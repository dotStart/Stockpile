plugin-dir = "plugins"
bind-address = "0.0.0.0:36623"

storage "redis" {
  address = "localhost:6379"
  // password = "admin1234"
  database = 0
}