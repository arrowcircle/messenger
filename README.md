# Messenger Service

[![Build Status](https://travis-ci.org/arrowcircle/chat.svg?branch=master)](https://travis-ci.org/arrowcircle/chat)[![Code Climate](https://codeclimate.com/github/arrowcircle/chat/badges/gpa.svg)](https://codeclimate.com/github/arrowcircle/chat)

# Installation

Get binary and run.

# Requirements

* Postgresql 9+

# Configuration

All configuration is set via env variables

* `CHAT_DATABASE_URL` - path to database, default: `postgres:///chat_development?sslmode=disable`
* `CHAT_BIND_ADDRESS` - host and port to run, default: `localhost:8080`
