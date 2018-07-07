Stockpile
=========

Lightweight and protocol aware API proxy for the Mojang APIs.

**Looking for Stockpile 1.0?** Its source code may still be accessed [here](https://github.com/dotStart/Stockpile/tree/v1.0.0-SNAPSHOT)

Key Features
------------

* Customizable (broad configuration options and plugin support)
* gRPC based (clients may be generated for most popular languages)
* Completely Open Source

Installation & Setup
--------------------

1. Download a binary release from the [releases page](https://github.com/dotStart/Stockpile/releases)
or [build from source](#Building)
2. Extract or copy the Stockpile executable and plugin directory (if applicable) to a custom
directory (convention for Linux systems is `/opt/stockpile`)
3. Create a configuration file (for examples, refer to the `docs` directory in the source
distribution)
4. Start Stockpile via `./stockpile server -config=path/to/myconfig.hcl`

Note that you may additionally launch Stockpile in development mode via `./stockpile server -dev` in
order to automatically select a set of acceptable default parameters without requiring a customized
configuration file.

Once the server is running, you may access its functionality via the `stockpile` executable. Run
`./stockpile help` for more information on the respective client commands. If configured you may
also access the Web UI via `<your ip or hostname>:36623` (or wherever you configured Stockpile to
listen) assuming that it has been enabled (the Web UI is enabled by default in development mode).

Prerequisites
-------------

The following applications are required to be present within your system PATH when **building** the
application from source (some may be omitted if you are willing to skip some of the build steps):

* git (latest recommended)
* go (1.10 or newer)
* node & npm (latest recommended)

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
