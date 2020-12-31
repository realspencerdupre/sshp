# SSHP
SSH login manager
 
SSHP keeps track of multiple SSH logins.

## Installation

First install Go version 1.15.6 or greater, then


```
git clone https://github.com/realspencerdupre/sshp.git

cd sshp

go build

```

## Usage
```
./sshp add
```
Launches the add host form. Sshp will ask you for several pieces of information, then save the host in the config.


```
./sshp
```
Will launch a prompt to pick a host, like this:
```
Select Day
  > MyWebsite (Me)
    ClientProject (Client1)

```
