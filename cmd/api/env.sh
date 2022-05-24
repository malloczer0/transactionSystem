#!/bin/bash
export HTTP_ADDR=localhost:8080

### DB settings

# Postgres settings
# shellcheck disable=SC2125
export PG_URL=postgres://postgres:postgres@localhost/postgres?sslmode=disable
export PG_MIGRATIONS_PATH=file://../../store/pg/migrations
export PG_USER=postgres
export PG_PASSWORD=postgres
# Logger settings
export LOG_LEVEL=debug