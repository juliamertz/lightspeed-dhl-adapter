#!/usr/bin/env bash

CLUSTER=https://api.webshopapp.com/nl
AUTH="${LIGHTSPEED_KEY}:${LIGHTSPEED_SECRET}"

# # Create new webhook
# curl -X POST $CLUSTER/webhooks.json \
#   -u $AUTH \
#   -d webhook[isActive]="true" \
#   -d webhook[itemGroup]="orders" \
#   -d webhook[itemAction]="created" \
#   -d webhook[language]="nl" \
#   -d webhook[format]="json" \
#   -d webhook[address]="https://nettenshop-staging.juliamertz.dev/webhook"

# # Disable webhook
# export WEBHOOK_ID=4680733
# curl -X PUT "$CLUSTER/webhooks/$WEBHOOK_ID.json" \
#   -u $AUTH \
#   -d webhook[isActive]="true"

# # Delete webhook
# webhook_id=4421635
# curl -X DELETE $CLUSTER/webhooks/${webhook_id}.json -u $AUTH

# # Get webhooks
# curl $CLUSTER/webhooks.json -u $AUTH | jq
