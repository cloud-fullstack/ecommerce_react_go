#!/bin/bash

export DB_HOST=${secrets.DB_HOST}
export DB_PORT=${secrets.DB_PORT}
export DB_USER=${secrets.DB_USER}
export DB_PASSWORD=${secrets.DB_PASSWORD}
export DB_NAME=${secrets.DB_NAME}
export SPACES_ACCESS_KEY=${secrets.SPACES_ACCESS_KEY}
export SPACES_SECRET_KEY=${secrets.SPACES_SECRET_KEY}



# Print out values for debugging
echo "DB_HOST: $DB_HOST"
echo "DB_PORT: $DB_PORT"
echo "DB_USER: $DB_USER"
echo "DB_PASSWORD: $DB_PASSWORD"
echo "DB_NAME: $DB_NAME"
echo "SPACES_ACCESS_KEY: $SPACES_ACCESS_KEY"
echo "SPACES_SECRET_KEY: $SPACES_SECRET_KEY"
echo "STATIC_PATH: $STATIC_PATH"
echo "FRON
