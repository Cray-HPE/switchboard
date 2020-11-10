#!/bin/bash

# Copyright 2020 Hewlett Packard Enterprise Development LP

echo "Configure PAM to use sssd..."
pam-config -a --sss --mkhomedir

echo "Generating broker host keys..."
ssh-keygen -A

echo "Checking for UAI_CREATION_CLASS..."
if ! [ -z $UAI_CREATION_CLASS ]; then
    echo UAI_CREATION_CLASS=$UAI_CREATION_CLASS >> /etc/environment
fi

echo "Starting sshd..."
/usr/sbin/sshd -f /etc/switchboard/sshd_config

echo "Starting sssd..."
sssd

sleep infinity
