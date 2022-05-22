#!/bin/bash
export HTTP_ADDR=localhost:8080

### DB settings

# Postgres settings
export PG_URL=postgres://postgres:postgres@localhost/postgres?sslmode=disable
export PG_MIGRATIONS_PATH=file://../../store/pg/migrations

# Logger settings
export LOG_LEVEL=debug