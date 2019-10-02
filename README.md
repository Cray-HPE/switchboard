![picture](img/switchboard.jpg)

## Switchboard
Switchboard is a lightweight utility to redirect users into a User Access Instance

### Usage
Switchboard works best as a ForceCommand setting in an sshd_config file.

```bash
Match User !root,*
	PermitTTY yes
	ForceCommand /opt/cray/switchboard.sh
```

This will redirect all users (except root) through the switchboard logic. The 
possible outcomes are:

1. Start a UAI if one is not already running and SSH to it once it is available.
2. SSH to a UAI already running if only one UAI is found
3. Choose a UAI to SSH to if multiple are found

Switchboard will autogenerate SSH keys if the user does not already have them
in their home directory at ~/.ssh/id_rsa.
