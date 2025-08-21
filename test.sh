#!/bin/bash

export MASTODON_INSTANCE_URL="your.mastodon.instance"
export MASTODON_TOKEN="your_access_token"
export SERVICE_URL="your.questionbasket.service"

go build
./questionbasket