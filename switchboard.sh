#!/bin/bash

function uas_output_parse() {
  echo $1 | grep $2| grep -o '".*' | tr -d '"'
}

export CRAY_FORMAT=json

echo "Authenticate with Cray DC LDAP..."
cray init --overwrite --no-auth \ 
          --hostname https://api-gw-service-nmn.local > /dev/null

# Check if a simple command returns 401. If it doesn't
# we don't need to force the user to authenticate again 
# as they may already have a valid token
if cray uas mgr-info list 2>&1 | grep --silent "401 Unauthorized"; then
  cray auth login --username $USER
fi

# Automatically generate rsa keys if they don't exist.
# These keys are required to create and log in to
# User Access Instances
yes no | ssh-keygen -f ~/.ssh/id_rsa -N "" > /dev/null

echo "Checking for running UAI..."
UAS_LIST=$(cray uas list --format json)

# Make sure we are able to parse the UAS_LIST with jq
if ! echo $UAS_LIST | jq -e > /dev/null; then
  echo "Could not parse list of UAIs..."
  exit 1
fi

NUM_UAI=$(echo $UAS_LIST | jq '.|length')
echo "Found $NUM_UAI instances:"
echo $UAS_LIST | jq -r '.[] | .uai_name'

if [ $NUM_UAI -gt 0 ]; then
  $(echo $UAS_LIST | jq -r '.[0] | .uai_connect_string')
fi
