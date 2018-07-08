---
title: Development Server - Getting Started
layout: docs-navigation
---

# Development Server

Once Stockpile is installed on your machine, you may proceed to running a
development instance to test and tinker with.

Simply execute the `stockpile server -dev` command:

```
$ stockpile server -dev
==> Stockpile Configuration

   Server Address: 127.0.0.1:36623
          Version: 2.0+git-4f4954f
      Commit Hash: 4f4954f96ebe83f1276e889b37c50b1283dee479
        Log Level: info
  Storage Backend: mem
              PID: 15572

==> TTL Configuration

           Names: 888h0m0s
         Profile: 168h0m0s
    Name History: 180h0m0s
       Blacklist: 168h0m0s

12:17:06.704 [WARN] stockpile :  Stockpile is running in development mode
12:17:06.704 [WARN] plugin :  plugins are unavailable on platform windows - Plugin manager startup has been skipped
12:17:06.704 [INFO] stockpile :  using database plugin: mem
12:17:06.704 [INFO] stockpile :  grpc server enabled
12:17:06.704 [WARN] stockpile :  legacy api enabled
12:17:06.704 [INFO] stockpile :  web ui enabled
```

The actual output of the command may differ slightly depending on the platform
on which you are running the application or its respective version.

Note that Stockpile does not fork (e.g. will stay in foreground). Thus you will
need to open another terminal (or a tab if supported) to execute any commands as
of this point.

## Verifying the Server State

In order to verify whether the server is correctly running (and listening for
commands), you may execute the `stockpile status` command:

```
$ stockpile status
Key            | Value
---------------+-----------------------------------------
Brand          | vanilla
---------------+-----------------------------------------
Version        | 2.0
---------------+-----------------------------------------
VersionFull    | 2.0+git-4f4954f
---------------+-----------------------------------------
CommitHash     | 4f4954f96ebe83f1276e889b37c50b1283dee479
---------------+-----------------------------------------
BuildTimestamp | 1531045291
```

Again, your command output may differ slightly.

You are additionally encouraged to open the web interface at
[http://localhost:36623](http://localhost:36623) in your favourite browser as it
provides valuable feedback while developing applications against Stockpile and
further illustrates the inner workings of the application.

## Next Steps

You may now move on to the [next section](first-requests.html).
