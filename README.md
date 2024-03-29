# restful-banking-database
This project provides an api for the managment of a sqlite database containing accounts and three currencies.
## ToDo
- [ ] https

## Installation
### Note: for Non-Linux users, these instructions might be off, especially for Windows users not using Bash. Feel free to amend this documentation with instructions for other platforms
Make sure you have `go` installed:
```sh
$ which go
```
If that doesn't give an output, or errors out, you probably should install go from your package manager or the [golang project](https://golang.org/).

If you're on Arch Linux, use:
```
# pacman -S go
```
First, clone the repo and cd into the directory:
```sh
$ git clone https://github.com/ArchWizard56/restful-banking-database.git
$ cd restful-banking-database
```
Next, run `make depend`, which will pull the necessary go dependencies, followed by `make`:
```sh
$ make depend
$ make
```
The binary will be placed in `bin/`. Move the binary wherever you want it:
```sh
$ mv bin/restful-banking-database your-location/
```
You'll have to copy the configuration files into the directory containing your binary (the cloned secrets file is named `default_secrets.json`:
```sh
$ cp config.json your-location/
$ cp default_secrets.json your-location/secrets.json
```
Make sure you change the secret value in `secrets.json` from `CHANGEME`to a secret key that only you know. Here's how you would do it with `sed`:
```sh
$ sed -i 's/CHANGEME/<yoursecret>/g' your-location/secrets.json
```
Finally, you can run the program:
```sh
$ cd your-location
$ ./restful-banking-database
```
The default port is `8050` and the default database is called `accounts.db`.
## Command line options
You can run the binary with the `-h` flag to show help
```sh
$ ./bin/restful-banking-database -h
Usage of ./bin/restful-banking-database:
  -c string
        the location of the config file to use (default "config.json")
  -d    toggles debug output
  -s string
        the location of the secrets file to use (default "secrets.json")
```
You can use `-c` to use a different location for the configuration file, `-s` to use a different location for the secrets file, and `-d` to toggle all of the debug output.
## Configuration file parameters

config.json:
```json
{
    "port": 8050,
    "database": "accounts.db"
}
```
Change the `port` parameter to set the port that the application listens on, and change the `database` parameter to change the location of the database file that the application uses. There's no need to initialize the database; the application will handle all of that.

secrets.json:
```json
{
    "jwtsecret":"CHANGEME"
}
```
You really should change the `jwtsecret` parameter in order to secure the jwt token generation.
## Updating the database
User accounts don't have a balance when created, so you'll have to manually update the balance of a new account to the initial value. This is done to prevent a user from creating many accounts and transfering the created amount to another account. In order to manually update the balance you can enter the database with:
```sh
$ sqlite3 <your-database>
```
The default name of the database is `accounts.db`, but it can be changed in the config file. After entering the database you can update the balance with a sql statement like this:
```sql
UPDATE accounts SET <balance type> = <default amount> WHERE AccountNumber = <account number>;
```
As an example:
```sh
$ sqlite3 accounts.db
> UPDATE accounts SET ArBalance = 10 WHERE AccountNUmber = 11198;
```
## API Documentation
See [the documentation](API.md)
