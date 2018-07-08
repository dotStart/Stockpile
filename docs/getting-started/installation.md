---
title: Installation - Getting Started
layout: docs-navigation
---

# Getting Started

To complete this guide, you will be required to install a copy of Stockpile on
your development machine (if you wish to test or develop plugins, you will
additionally be required to run the application on a Linux (virtual) machine).

## Installation

To install Stockpile, find the [appropriate package](../) for your system and
download it. Stockpile is packaged as a gzip compressed tarball.

After downloading Stockpile, extract the package to a directory of your choosing
(it is recommended to add the directory to your `$PATH` variable if you wish for
quick access). The application is packaged as a single executable called
`stockpile` (or `stockpile.exe` if you are on Windows) which provides the server
implementation and a selected set of client commands for testing and maintenance
purposes. If you are working with Linux, you will additionally receive a set of
plugins (within the `plugins` directory). It is up to you whether you wish to
extract them as well.

## Verifying the Installation

Once you extracted Stockpile, you may verify whether the application has been
correctly installed. When you invoke `stockpile` (or `stockpile.exe`) without
any arguments you should see an output similar to this:

```
$ stockpile
Usage: stockpile <flags> <subcommand> <subcommand args>

Subcommands:
	commands         list all command names
	help             describe subcommands and their syntax
	server           starts a new Stockpile server instance

Subcommands for Client:
	check-blacklist  checks whether an address has been blacklisted using a remote Stockpile server
	get-id           queries a user's profile id from a Stockpile server
	get-profile      queries a user's profile from a Stockpile server
	listen           listens for events on a Stockpile server
	name-history     queries a profile's name history
	plugins          displays a list of commands loaded on the Stockpile server
	status           displays the current status of a Stockpile server


Use "stockpile flags" for a list of top-level flags
```

If the command execution results in an error citing that the command could not
be found, you should verify your `$PATH` variable.

Otherwise Stockpile is ready to go.

## Next Steps

You may continue with the [next section](dev-server.html) now.
