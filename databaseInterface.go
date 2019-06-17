package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"time"
)

// Function to load a database, or set up a new database, then pass it back
func LoadDB(DBPath string) *sql.DB {
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		DualErr(err)
	}
	// Initialize Default table consisting of:
	// AccountNumber    OwnerName   Password    DiscordName     TokValue    DcBalance   CcBalance   ArBalance
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS accounts (AccountNumber TEXT PRIMARY KEY, OwnerName TEXT, Password TEXT, DiscordName TEXT, TokValue INTEGER, DcBalance INTEGER, CcBalance INTEGER, ArBalance INTEGER)")
	statement.Exec()
    DualDebug("Initialized Database")
	return db
}

//Simple Function to check if a string is in a slice
func stringNotInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return false
		}
	}
	return true
}

func genAccountNumber(db *sql.DB) string {
	// Query db for the AccountNumber of all active accounts and add them to a slice of strings
	var AccountNumber string
    DualDebug("Checking accounts for AccountNumber")
	accounts, err := db.Query("SELECT AccountNumber FROM accounts;")
	ActiveAccounts := make([]string, 1)
	for accounts.Next() {
		err := accounts.Scan(&AccountNumber)
		if err != nil {
			DualErr(err)
		}
		ActiveAccounts = append(ActiveAccounts, AccountNumber)
	}
	//Infinite loop until generation succeeds
	for i := 0; i < 1; {
        DualDebug("Attempting to generate account number")
		//Seed rng
		rand.Seed(time.Now().UnixNano() + 123432)
		//Generate an account number
		baseNumber := rand.Intn(50000)
		//Convert the account number into a string and then into bytes
		stringNumber := strconv.Itoa(baseNumber)
		bytesNum := []byte(stringNumber)
		//Loop over the account number and add leading zeros until it's a five digit number
		for zerosNeeded := len(bytesNum); zerosNeeded < 5; zerosNeeded++ {
			zeroArray1 := []byte("0")
			bytesNum = append(zeroArray1, bytesNum...)
		}
		NewAccountNumber := string(bytesNum)
		if err != nil {
			DualErr(err)
		}
		// Check if the generated account number exists in the list of active accounts, if not, exit loop
		if stringNotInSlice(NewAccountNumber, ActiveAccounts) {
           DualDebug("Success!") 
			return NewAccountNumber

		}
       DualDebug("Failed! Trying again") 
	}
	// Dummy return
	return ""
}

//Function to create a main account that depends on a username and password
func CreateMainAccount(db *sql.DB, OwnerName string, Password []byte) error {
    DualInfo("Creating main account")
	var Owner string
	var TokValue int
	accounts, err := db.Query("SELECT OwnerName FROM accounts;")
	Owners := make([]string, 1)
	for accounts.Next() {
		err := accounts.Scan(&Owner)
		if err != nil {
			DualErr(err)
		}
		Owners = append(Owners, Owner)
	}
	if stringNotInSlice(OwnerName, Owners) {
		//Generate Token ID so that tokens can be revoked
		rand.Seed(time.Now().UnixNano() + 123432)
		TokValue = rand.Intn(999)
		if err != nil {
			DualErr(err)
		}
	} else {
		return errors.New("Account already exists!")
	}
	//Generate new Account Number
	AccountNumber := genAccountNumber(db)
	//Create a password hash
	hash, err := bcrypt.GenerateFromPassword(Password, 15)
	//Create the account entry in the database
	statement, _ := db.Prepare("INSERT INTO accounts VALUES ($1, $2, $3, $4, $5, $6, $7, $8);")
	statement.Exec(AccountNumber, OwnerName, hash, "none", TokValue, 0, 0, 10)
    return nil
}
