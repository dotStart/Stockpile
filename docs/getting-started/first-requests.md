---
title: Your first Requests - Getting Started
layout: docs-navigation
---

# Your first Requests

With the development server ready to go, you are ready to send a few requests
using the command line client.

Let's try resolving the profileId assigned to a display name at the moment using
the `stockpile get-id` command:

```
$ stockpile get-id dotStart

Key         | Value
------------+-------------------------------------
Id          | d71a5dac-4e71-443b-8158-4389c269e44d
------------+-------------------------------------
Name        | dotStart
------------+-------------------------------------
FirstSeenAt | 2018-07-08 22:35:44 +0200 CEST
------------+-------------------------------------
LastSeenAt  | 2018-07-08 22:35:44 +0200 CEST
------------+-------------------------------------
ValidUntil  | 2018-08-14 22:35:44 +0200 CEST
```

If you wish, you may additionally pass the `--time=<date & time>` parameter to
the command in order to check an association at a specific time:

```
$ stockpile get-id --time=2017-10-02T15:11:52+01:00 dotStart
Key         | Value
------------+-------------------------------------
Id          | d71a5dac-4e71-443b-8158-4389c269e44d
------------+-------------------------------------
Name        | dotStart
------------+-------------------------------------
FirstSeenAt | 2017-10-02 16:11:52 +0200 CEST
------------+-------------------------------------
LastSeenAt  | 2017-10-02 16:11:52 +0200 CEST
------------+-------------------------------------
ValidUntil  | 2017-11-08 15:11:52 +0100 CET
```

Note that the time flag currently also accepts two special values:

* `now` which refers to the current date and time (this is the default if
  unspecified)
* `0` which resolves the original owner of a name

## Profile Lookups

The same general process applies to player profiles via the
`stockpile get-profile` command:

```
$ stockpile get-profile d71a5dac-4e71-443b-8158-4389c269e44d
Key        | Value
-----------+----------------------------------------------------------------
Id         | d71a5dac-4e71-443b-8158-4389c269e44d
-----------+----------------------------------------------------------------
Name       | dotStart
-----------+----------------------------------------------------------------
Properties | textures:
           |   Name: textures
           |   Value: eyJ...X0=
           |   Signature: OZS...Y8=
-----------+----------------------------------------------------------------
Textures   | Timestamp: 2018-07-08 22:58:52 +0200 CEST
           | ProfileId: d71a5dac-4e71-443b-8158-4389c269e44d
           | ProfileName: dotStart
           | Textures: SKIN:
           |             http://textures.minecraft.net/texture/d38a4d64311...
           |           CAPE:
           |             
```

## Command Documentation

All commands in Stockpile come with a general description and one or more usage
examples. Generally all command documentations may be accessed via
`stockpile help <command>`:

```
$ stockpile help get-profile
Usage: stockpile get-profile [options] <id>

This command retrieves a profile from a Stockpile server:

  $ stockpile get-profile d71a5dac-4e71-443b-8158-4389c269e44d

Available command specific flags:

  -server-address string
    	specifies the address of the target server (default "0.0.0.0:36623")
```

A complete list of all commands is available via `stockpile commands`.

## Next

You may now continue with the next step:
[Configuration](configuration.html)
