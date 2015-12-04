# Messenger Service

[![Build Status](https://travis-ci.org/arrowcircle/messenger.svg?branch=master)](https://travis-ci.org/arrowcircle/messenger)

# Installation

Get binary and run.

# Requirements

* Postgresql 9+

# Configuration

All configuration is set via env variables

* `CHAT_DATABASE_URL` - path to database, default: `postgres:///chat_development?sslmode=disable`
* `CHAT_BIND_ADDRESS` - host and port to run, default: `localhost:8080`
