bind-address = "127.0.0.1:36623"

// file storage is technically suited for small production deployments, however, a proper storage
// server like redis is recommended for higher volumes
storage "file" {
  path = "data"
}
