# Copyright 2020 Hewlett Packard Enterprise Development LP

#!/bin/bash

echo "Configure PAM to use sssd..."
pam-config -a --sss --mkhomedir

echo "Generating broker host keys..."
ssh-keygen -A

echo "Starting sssd..."
sssd

echo "Starting sshd..."
/usr/sbin/sshd -f /etc/switchboard/sshd_config

sleep infinity
