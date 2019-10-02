#!/bin/bash

# Set format to json so jq can be used
export CRAY_FORMAT=json
export CRAY_AUTH_LOGIN_USERNAME=$USER

function create_uai() {
  UAS_CREATE=$(cray uas create --publickey ~/.ssh/id_rsa.pub)
}

# Automatically generate rsa keys if they don't exist.
# These keys are required to create and log in to UAIs
function verify_ssh_keys_exist() {
  yes no | ssh-keygen -f ~/.ssh/id_rsa -N "" > /dev/null
}

echo "Checking for authentication with Keycloak..."
cray init --overwrite --no-auth \
          --hostname https://api-gw-service-nmn.local > /dev/null

# Check if a simple command returns 401. If it doesn't,
# we don't need to force the user to authenticate again 
# as they may already have a valid token
if cray uas mgr-info list 2>&1 | grep --silent "401 Unauthorized"; then
  echo "Log in as $USER to Keycloak..."
  cray auth login
fi

echo "Checking for running UAIs..."
UAS_LIST=$(cray uas list)

# Make sure we are able to parse the UAS_LIST with jq
if ! echo $UAS_LIST | jq -e > /dev/null; then
  echo "Could not parse list of UAIs..."
  exit 1
fi

# TODO use this nifty trick
# jq '.[] | select(.uai_status=="Running: Ready")'

NUM_UAI=$(echo $UAS_LIST | jq '.|length')

if [ $NUM_UAI -lt 1 ]; then
  echo "Creating a UAI..."
  create_uai
  exit 0
fi

if [ $NUM_UAI -eq 1 ]; then
  echo "Using existing UAI connection string..."
  $(echo $UAS_LIST | jq -r '.[0] | .uai_connect_string')
  exit 0
fi

if [ $NUM_UAI -gt 1 ]; then
  echo $UAS_LIST | jq -r '.[] | "\(.uai_name) \t \(.uai_status) \t \(.uai_age) \t \(.uai_img)"' | awk '{print NR, "\t", $0}'
  read -p "Select a UAI by number: " selection
  selection="$(($selection-1))"
  # TODO fix this
  #if ! (( 0 <= $selection < $NUM_UAI )); then
  #  echo "Invalid selection"
  #  exit 1
  #fi
  echo "Logging in to UAI:"
  echo $UAS_LIST | jq -r --arg INDEX $selection '.[$INDEX|tonumber] | .uai_connect_string'
  $(echo $UAS_LIST | jq -r --arg INDEX $selection '.[$INDEX|tonumber] | .uai_connect_string')
  exit 0
fi
