#!/bin/bash

export MASTODON_INSTANCE_URL="your.mastodon.instance"
export MASTODON_TOKEN="your_access_token"
export SERVICE_URL="localhost"
export DATABASE_URL="data.db"

go build
./questionbasket