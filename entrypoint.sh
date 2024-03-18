#!/bin/sh

echo "APP_ENV = ${APP_ENV}" > .env
echo "APP_PORT = ${APP_PORT}" >> .env

/app/main