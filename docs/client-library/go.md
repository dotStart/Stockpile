---
title: Go - Client Libraries
layout: docs-navigation
---

# Go Client Library

The Go client library is provided
[as part of the main repository](https://github.com/dotStart/Stockpile/tree/develop/client)
and used as a base for the built-in CLI commands.

It is highly recommended to use a dependency management tool, such as
[dep](https://golang.github.io/dep/) as the API may change every once in a
while.

## Getting Started

The base of the client is the `github.com/dotStart/Stockpile/client.Stockpile`
struct which may be initialized using the `New` function:

```go
import (
  "fmt"
  "github.com/dotStart/Stockpile/client"
)

func main() {
  client, err := client.New("localhost:36623")
  if err != nil {
    panic(err)
  }

  profileId, err := client.GetProfileId("dotStart", time.Now())
  if err != nil {
    panic(err)
  }

  fmt.Printf("ProfileId: %s\n", profileId.Id)
}
```

## Next Steps

For more information, refer to the client's
[godoc](https://godoc.org/github.com/dotStart/Stockpile/client)
