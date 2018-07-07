Stockpile
=========

Lightweight and protocol aware API proxy for the Mojang APIs.

**Looking for Stockpile 1.0?** Its source code may still be accessed [here](https://github.com/dotStart/Stockpile/tree/v1.0.0-SNAPSHOT)

Key Features
------------

* Customizable (broad configuration options and plugin support)
* gRPC based (clients may be generated for most popular languages)
* Completely Open Source

Building
--------

1. `go get -d -u github.com/dotStart/Stockpile/...`
2. `cd $(go env GOPATH)/src/github.com/dotStart/Stockpile`
3. `make`

The resulting binaries will be located in the `build` directory.

**Note:** Binaries built on Windows will not provide support for plugins (nor will they include the
standard plugins). If you wish to build a version for distribution (or execution in a production
environment), please build them on Linux (for instance, using the Vagrant configuration included
within this repository).

License
-------

```
Copyright [year] [name] <[email]>
and other copyright owners as documented in the project's IP log.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
