#!/bin/bash

# Set the environment variables for the Go app.
export DB_USERNAME=root
export DB_PASSWORD=supersecret
export DB_HOST=DB_IP_ADDRESS
# export DB_UNIX_SOCKET=/cloudsql/PROJECT_ID:REGION:INSTANCE_NAME
export DB_PORT=3306
export DB_NAME=handson
export PROJECT_ID=PROJECT_ID
export PORT=8080

# for local debug
# export GOOGLE_APPLICATION_CREDENTIALS=/Users/USERNAME/.config/gcloud/application_default_credentials.json

# Run the Go app.
go run main.go
