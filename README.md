# restful-banking-database
This project provides an api for the managment of a sqlite database containing accounts and three types of balance.

## Installation
### Note: for non linux users, these instructions might be off, especially for windows users not using the subsystem. Feel free to amend this documentation with instructions for other platforms
Make sure you have `go` installed:
```
$ which go
```
If that doesn't give an output, or errors out, you probably should install go from your package manager or the [golang project](https://golang.org/).

If you're on Arch Linux, use:
```
# pacman -S go
```
First, clone the repo and cd into the directory:
```
$ git clone https://github.com/ArchWizard56/restful-banking-database.git
$ cd restful-banking-database
```
Next, run  `make`, which will pull the necessary go dependencies:
```
$ make
```
The binary will be placed in `bin/`. Move the binary wherever you want it:
```
$ mv bin/restful-banking-database your-location/
```
You'll have to copy the configuration files into the directory containing your binary (the secret file is named `default_secrets.json`:
```
$ cp config.json your-location/
$ cp default_secrets.json your-location/secrets.json
```
Make sure you change the secret value in `secrets.json` from `CHANGEME`to a secret key that only you know. Here's how you would do it with `sed`:
```
$ sed -i 's/CHANGEME/yoursecret/g' your-location/secrets.json
```
Finally, you can run the program:
```
$ cd your-location
$ ./restful-banking-database
```
The default port is `8050` and the default database is called `accounts.db`.
