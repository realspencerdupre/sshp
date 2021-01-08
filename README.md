# SSHP
SSH login manager
 
SSHP keeps track of multiple SSH logins.

## Installation

To install on Ubuntu 20.04


```
# The following line automatically creates an ssh key if you don't already have one
test ! -e ~/.ssh/id_rsa && ssh-keygen

# Install the go language if you don't already have it
sudo apt install golang-go

# Install sshp
go get github.com/realspencerdupre/sshp
```

## Usage

### Adding a host
```
sshp add
```
Follow the prompts to add a new host.

```
Username user
Host (or IP) 127.0.0.1
✔ Owner: Me
✔ Description: server 3
Password ********
```

Your password is not saved, only used for the first connection.

### Connecting to selected host
```
sshp
```
Will launch a prompt to pick a host, looking something like this:
```
Select Day
  > server 3 (Me)
    Other server (Client1)
```
