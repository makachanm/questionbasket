#!/bin/bash

export DATABASE_URL="data.db"

go build
./questionbasket -init