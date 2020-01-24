![Deliverzes Kigi](logo.svg)

# Deliverzes

`Deliverzes` is courier delivery message, transforming your http request to telegram bot chat.

## Overview

## Prerequisite

To run this program, you will need

### List System & App Dependencies

```$xslt
- Golang 1.10+
- Go Mod Enabled
- Badger K/V DB
- Telebot
- GotEnv
```

## How to Run

- Copy environment file from `env.example` to be `.env` or use `ENVIRONMENT VARIABLE` directly
- Verify and download dependencies `go mod download && go mod verify`
- Run the app `go run main.go`


## Deployment

Create directory for Mounted storage
`mkdir /tmp/badger`

Run the latest Image
`docker run -v /tmp/badger:/tmp/badger --env APP_URI=0.0.0.0 --env HTTP_PORT=8000 --env TELEGRAM_TOKEN=YOUR_TELEGRAM_TOKEN --env TELEGRAM_POLLING_INTERVAL=1 --env DB_PATH=/tmp/badger -p 8000:8000 registry.gitlab.com/aeidelos/deliverzes:latest`
