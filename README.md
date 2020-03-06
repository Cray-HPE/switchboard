![picture](img/switchboard.jpg)

## Switchboard
Switchboard is a convenience utility to automate the process of managing and connecting
to User Access Instances. When used as a CLI, switchboard offers the following subcommands:
$ switchboard start|list|delete

### Requirements

* The Cray CLI command is present in the users PATH
* The Shasta API gateway is known and accessible 
* The user is configured in Keycloak and knows their password
* The user has RBAC sufficient to create, list, and delete UAIs

### Installation

```bash
zypper install cray-switchboard
# Optionally start switchboard with sshd
systemctl start cray-switchboard-sshd
```

An ansible role is also available that will simply start the service.
switchboard/tasks/main.yml:
```bash
- name: Start the cray-switchboard-sshd service
  systemd:
    state: started
    enabled: yes
    name: cray-switchboard-sshd
```

### Usage
switchboard may be run as a interactive command but it works best as a 
ForceCommand setting in an sshd_config file:

```bash
Match User !root,*
	PermitTTY yes
	ForceCommand /opt/cray/switchboard/src/switchboard
```

This `Match User` block may be added to the sshd listening on port 22 but by 
default switchboard is configured to start an sshd service listening on
port 203.

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
arbus:~ $ ssh -p 203 slice-uan01
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

The switchboard command will also perform the same logic when used interactively:
```bash
alanm@uan01:~> switchboard
Checking for authentication with Keycloak...
Verifying ssh keys exist...
Checking for running UAIs...
Using existing UAI connection string...
Warning: Permanently added '[172.30.52.72]:30980' (ECDSA) to the list of known hosts.
   ______ ____   ___ __  __   __  __ ___     ____
  / ____// __ \ /   |\ \/ /  / / / //   |   /  _/
 / /    / /_/ // /| | \  /  / / / // /| |   / /
/ /___ / _, _// ___ | / /  / /_/ // ___ | _/ /
\____//_/ |_|/_/  |_|/_/   \____//_/  |_|/___/

alanm@uai-alanm-6034192a-7968b8b47b-f58n6:~>
```
