# SSHP
SSH login manager
 
SSHP keeps track of multiple SSH logins.

## Installation

To install on Ubuntu 20.04


```
sudo apt install git golang-go

go get github.com/realspencerdupre/sshp

git clone https://github.com/realspencerdupre/sshp.git

cd sshp

go build

```

## Usage

### Adding
```
./sshp add

```
Follow the prompts to add a new host.

```
Username user
Host (or IP) 127.0.0.1
✔ Owner: Me█
✔ Description: server 3█
Password ********
```

### Connecting
```
./sshp
```
Will launch a prompt to pick a host, looking something like this:
```
Select Day
  > MyWebsite (Me)
    ClientProject (Client1)

```
