# openHAB cli client [![Build Status](https://travis-ci.org/dereulenspiegel/openhab-cli.svg?branch=master)](https://travis-ci.org/dereulenspiegel/openhab-cli)

This is a simple command line client for [openHAB](http://www.openhab.org/).
The intention is to create a client which interacts well with shell scripts and
that is easy to use for command line jockeys who want to control their home
without having to open a browser or flip out the smartphone.

## Installation

´go get github.com/dereulenspiegel/openhab-cli/oh´

The openHAB cli commands expect a minimal configuration in your home directory.
You have to create the file ~/.openhab-cli with the following content:

```yaml
host: http://example.openhab:8080
```

## Usage

* List all items: `oh list` or shorter `oh l`
* Get state of an item: `oh state <item name>` or shorter `oh s <item name>`
* Send command: `oh command <item name> <command>` or shorter `oh c <item name> <command>`

More to follow
