![picture](img/switchboard.jpg)

## Switchboard
Switchboard is a lightweight utility to redirect users into a User Access Instance.

### Requirements
Switchboard sits at the top of a fairly large software stack. Some fundamental 
requirements are:
* Kubernetes
* cray-uas-mgr
* Keycloak

Switchboard also requires that users are able to ssh AND authenticate with
`cray auth login`.


### Installation
An rpm will eventually be available.

cd /opt/cray/
git clone ssh://git@stash.us.cray.com:7999/~alanm/switchboard.git
/usr/sbin/sshd -f switchboard/sshd_config

### Usage
Switchboard may be run as a interactive script but it works best as a 
ForceCommand setting in an sshd_config file:

```bash
Match User !root,*
	PermitTTY yes
	ForceCommand /opt/cray/switchboard/bin/switchboard.sh
```

This `Match User` block may be added to the sshd listening on port 22 or it
may also be used for an sshd server listening on an alternate port.

This will redirect all users (except root) through the switchboard logic. The 
possible outcomes are:

1. Start a UAI if one is not already running and SSH to it once it is available.
2. SSH to a UAI already running if only one UAI is found
3. Choose a UAI to SSH to if multiple are found

Switchboard will autogenerate SSH keys if the user does not already have them
in their home directory at `~/.ssh/id_rsa`. The user may also be prompted for
an additional password via `cray auth login` if they do not already have a valid
token. Depending on the refresh timeout of the token, subsequent logins may not
ask for a second password.

### Example
```bash
arbus:~ $ ssh -p 40 slice-uan01
Warning: Permanently added '[slice-sms]:40,[172.30.52.70]:40' (ECDSA) to the list of known hosts.
Password:
Checking for authentication with Keycloak...
Verifying ssh keys exist...
Checking for running UAIs...
1 	 uai-alanm-017a3df1 	 Running: Ready 	 16m 	 sms.local:5000/cray/cray-uas-sles15-slurm:latest
2 	 uai-alanm-fddee5bb 	 Running: Ready 	 15m 	 sms.local:5000/cray/cray-uas-sles15-pbs:latest
Select a UAI by number: 1
Logging in to UAI:
ssh alanm@172.30.52.72 -p 32283 -i ~/.ssh/id_rsa
   ______ ____   ___ __  __   __  __ ___     ____
  / ____// __ \ /   |\ \/ /  / / / //   |   /  _/
 / /    / /_/ // /| | \  /  / / / // /| |   / /
/ /___ / _, _// ___ | / /  / /_/ // ___ | _/ /
\____//_/ |_|/_/  |_|/_/   \____//_/  |_|/___/

alanm@uai-alanm-017a3df1-54d985fb7b-5k7w2:~>
```
