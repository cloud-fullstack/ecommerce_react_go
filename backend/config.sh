#!/bin/bash

export DB_HOST="${DB_HOST}"
export DB_PORT="${DB_PORT}"
export DB_USER="${DB_USER}"
export DB_PASSWORD="${DB_PASSWORD}"
export DB_NAME="${DB_NAME}"
export SPACES_ACCESS_KEY="${SPACES_ACCESS_KEY}"
export SPACES_SECRET_KEY="${SPACES_SECRET_KEY}"

# Debug output
echo "DB_HOST: $DB_HOST"
echo "DB_PORT: $DB_PORT"
echo "DB_USER: $DB_USER"
echo "DB_PASSWORD: $DB_PASSWORD"
echo "DB_NAME: $DB_NAME"
echo "SPACES_ACCESS_KEY: $SPACES_ACCESS_KEY"
echo "SPACES_SECRET_KEY: $SPACES_SECRET_KEY"
echo "STATIC_PATH: $STATIC_PATH"
