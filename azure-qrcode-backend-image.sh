#!/bin/bash

set -euo pipefail

echo "Logging into trpcontainers"

az acr login -n trpcontainers

TAG="qrcode-backend"

echo "Now building ${TAG}"

docker build -t $TAG .

AZURETAG="trpcontainers.azurecr.io/${TAG}"

docker tag $TAG $AZURETAG

docker push $AZURETAG

echo "Image push for ${AZURETAG} complete"