#!/usr/bin/env bash

URL="https://api-gw.dhlparcel.nl"
USER_ID=""
API_KEY=""

# # Authenticate with dhl
# token=$(curl "$URL/authenticate/api-key" \
#   -X POST \
#   -H "Content-Type: application/json" \
#   --data "{ \"userId\": \"$USER_ID\", \"key\": \"$API_KEY\" }" \
#   | jq '.accessToken' | xargs)

# # get orders with wrong account id
# curl -v -X 'GET' \
#   'https://api-gw.dhlparcel.nl/drafts?accountNumbers=<account-id>' \
#   -H 'accept: application/json' \
#   -H "Authorization: Bearer $token" | jq > wrong-orders.json

# # Update / create order
# curl -v -X 'POST' \
#   'https://api-gw.dhlparcel.nl/drafts' \
#   -H 'accept: application/json' \
#   -H "Authorization: Bearer $token" \
#   -H 'Content-Type: application/json' \
#   -d "$(cat ./new-order.json)"

# # Update orders with wrong account id
# json_data=$(cat ./wrong-orders.json)
# echo "$json_data" | jq -c --arg newId "$NEW_ACCOUNT_ID" '.[] | .accountId = $newId' | while read -r obj; do
#   curl -v -X 'POST' \
#     'https://api-gw.dhlparcel.nl/drafts' \
#     -H 'accept: application/json' \
#     -H "Authorization: Bearer $token" \
#     -H 'Content-Type: application/json' \
#     -d "$obj"
# done
