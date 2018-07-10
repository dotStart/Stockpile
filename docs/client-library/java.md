---
title: Java - Client Libraries
layout: docs-navigation
---

# Java Client Library

The Java client library is provided
[as part of an additional repository](https://github.com/dotStart/Stockpile-Java/tree/develop/client).

## Getting Started

The base of the client is the `io.github.dotstart.stockpile.Stockpile` class
which provides multiple constructors for various different scenarios:

```java
import io.github.dotstart.stockpile.Stockpile;
import io.github.dotstart.stockpile.entity.profile.ProfileId;
import java.time.Instant;

public class MyApplication {
  public static void main(String[] args) throws InterruptedException {
    try (Stockpile client = new Stockpile("127.0.0.1")) {
      ProfileId profileId = client.profileOperations().getProfileId("dotStart", Instant.now());
      System.out.println("ProfileId: " + profileId.getId())
    }
  }
}
```
