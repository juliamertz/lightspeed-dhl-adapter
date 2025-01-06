#!/usr/bin/env bash

key=""
secret=""

CLUSTER=https://api.webshopapp.com/nl
AUTH="${key}:${secret}"

# Create new webhook
# curl -X POST $CLUSTER/webhooks.json \
#   -u $AUTH \
#   -d webhook[isActive]="true" \
#   -d webhook[itemGroup]="orders" \
#   -d webhook[itemAction]="created" \
#   -d webhook[language]="nl" \
#   -d webhook[format]="json" \
#   -d webhook[address]="https://nettenshop.juliamertz.dev/webhook"

# Delete webhook
# webhook_id=4421635
# curl -X DELETE $CLUSTER/webhooks/${webhook_id}.json -u $AUTH

# Get webhooks
# curl $CLUSTER/webhooks.json -u $AUTH | jq
