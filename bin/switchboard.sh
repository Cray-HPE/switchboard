#!/bin/bash

if [ $EUID -eq 0 ]; then
   echo "switchboard.sh must not be run as root" 
   exit 1
fi

# Set format to json so jq can be used
export CRAY_AUTH_LOGIN_USERNAME=$USER
READY_RETRIES=15
SPIN='-\|/'

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

echo "Verifying ssh keys exist..."
verify_ssh_keys_exist

# Check if a simple command returns 401. If it doesn't,
# we don't need to force the user to authenticate again 
# as they may already have a valid token
#TODO Fix tokens expiring with /userinfo
#if cray uas mgr-info list 2>&1 | grep --silent "401 Unauthorized"; then
if cray uas list 2>&1 | grep --silent "Token not valid for UAS"; then
  echo "cray auth login --username $USER..."
  #TODO add retries 
  cray auth login
fi

echo "Checking for running UAIs..."
UAS_LIST=$(cray uas list --format json)

# Make sure we are able to parse the UAS_LIST with jq
if ! echo $UAS_LIST | jq -e .[] 2>&1 /dev/null; then
  echo "Could not parse list of UAIs..."
  exit 1
fi

NUM_UAI=$(echo $UAS_LIST | jq '.|length')

if [ $NUM_UAI -lt 1 ]; then
  echo "Creating a UAI..."
  create_uai
  for i in $(seq 1 $READY_RETRIES); do
    if cray uas list --format json | jq -e '.[] | select(.uai_status=="Running: Ready")'; then
      break
    elif [ $i -eq $READY_RETRIES ]; then
      echo "Timed out waiting for UAI"
      exit 1
    else
      printf "\r${SPIN:$(((i+1)%4)):1}"
      sleep 1
    fi
  done
  # We got here so there is a UAI in "Running: Ready"
  $(cray uas list --format json | jq -r '.[0] | .uai_connect_string') $SSH_ORIGINAL_COMMAND
  exit 0
fi

if [ $NUM_UAI -eq 1 ]; then
  echo "Using existing UAI connection string..."
  $(echo $UAS_LIST | jq -r '.[0] | .uai_connect_string') $SSH_ORIGINAL_COMMAND
  exit 0
fi

if [ $NUM_UAI -gt 1 ]; then
  # Print a table of UAIs and prepend a number (awk) so 
  # users are able to select a UAI by number
  echo $UAS_LIST | jq -r '.[] | "\(.uai_name) \t \(.uai_status) \t \(.uai_age) \t \(.uai_img)"' | awk '{print NR, "\t", $0}'
  read -p "Select a UAI by number: " selection
  selection="$(($selection-1))"
  # TODO fix this
  #if ! (( 0 <= $selection < $NUM_UAI )); then
  #  echo "Invalid selection"
  #  exit 1
  #fi
  $(echo $UAS_LIST | jq -r --arg INDEX $selection '.[$INDEX|tonumber] | .uai_connect_string') $SSH_ORIGINAL_COMMAND
  exit 0
fi
