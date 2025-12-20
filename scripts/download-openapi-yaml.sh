#!/bin/bash

source .env

curl -o openapi.yaml --location -g --request POST "https://api.apidog.com/v1/projects/${APIDOG_PROJECT_ID}/export-openapi?locale=en-US" \
--header 'X-Apidog-Api-Version: 2024-03-28' \
--header "Authorization: Bearer ${APIDOG_ACCESS_TOKEN}" \
--header 'Content-Type: application/json' \
--data-raw '{
  "scope": {
    "type": "ALL"
  },
  "options": {
    "includeApidogExtensionProperties": false,
    "addFoldersToTags": false
  },
  "oasVersion": "3.0",
  "exportFormat": "YAML"
}'
