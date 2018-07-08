---
title: Introduction
layout: docs-navigation
---

# Introduction

Welcome to the introductory guide to Stockpile. This guide has been designed to
get you setup and running as fast as possible. Specifically this part of the
documentation contains information on Stockpile's features, potential use cases
and how to install the program.

## What is Stockpile?

Stockpile is a protocol aware caching system which has been designed to
temporarily store data from the [Mojang API](http://wiki.vg/Mojang_API) in order
to prevent unnecessary requests and simplify tasks such as the retrieval of
name associations.

Specifically Stockpile provides:

* **Persisted name associations**<br />
  Within the default configuration, name associations will be stored for the
  maximum permitted time (e.g. names will not suddenly disappear when a player
  changes their name)
* **Rate limit management**<br />
  Especially on systems where multiple servers are running on the same physical
  machine, you may end up exceeding the rate limit of 600 requests per 10
  minutes (or 1 request per minute for profiles) - Stockpile alleviates this
  issue by making the information globally accessible to all instances even when
  many implementations request information
* **Centralized Storage**<br />
  Stockpile provides a central point of storage and thus gets rid of the need to
  implement custom caching logic in your code
* **Events**<br />
  Automatically perform tasks within your infrastructure when new elements are
  written to the cache

## Next Steps

You may continue by reading the
[Getting Started](getting-started/installation.html) guide.
