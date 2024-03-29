package main

import (
	crand "crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
    "net/http"
	"time"
)

type Account struct {
	Username string
	Number   string
	CcBal    int
	DcBal    int
	ArBal    int
}

type Transfer struct {
	Username    string
	FromAccount string
	ToAccount   string
	Amount      int
	Type        string
}

type TokenValueHolder struct {
	TokValue string
	StoredAt int64
}

var TokenCache map[string]TokenValueHolder

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := crand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// Function to load a database, or set up a new database, then pass it back
func LoadDB(DBPath string) *sql.DB {
	db, err := sql.Open("sqlite3", DBPath+"?cache=shared&mode=rwc")
	db.SetMaxOpenConns(10)
	if err != nil {
		DualErr(err)
	}
	// Initialize Default table consisting of:
	// AccountNumber    OwnerName   Password    DiscordName     TokValue    DcBalance   CcBalance   ArBalance
	var statement *sql.Stmt
	statement, err = db.Prepare("CREATE TABLE IF NOT EXISTS accounts (AccountNumber TEXT PRIMARY KEY, OwnerName TEXT, Password TEXT, DiscordName TEXT, TokValue TEXT, DcBalance INTEGER, CcBalance INTEGER, ArBalance INTEGER)")
	if err != nil {
		DualErr(err)
	}
	_, err = statement.Exec()
	if err != nil {
		DualErr(err)
	}
	DualDebug("Initialized Database")
	TokenCache = make(map[string]TokenValueHolder)
	DualDebug("Initialized Token Cache")
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
	defer accounts.Close()
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
func CreateMainAccount(db *sql.DB, OwnerName string, Password []byte) (Account, error) {
	DualInfo("Creating main account")
	var Owner string
	var TokValue string
	accounts, err := db.Query("SELECT OwnerName FROM accounts;")
	defer accounts.Close()
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
		TokValue, err = GenerateRandomString(32)
		if err != nil {
			DualErr(err)
		}
	} else {
		return Account{}, errors.New("Account already exists!")
	}
	//Generate new Account Number
	AccountNumber := genAccountNumber(db)
	//Create a password hash
	hash, err := bcrypt.GenerateFromPassword(Password, 15)
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
		return Account{}, err
	}
	//Create the account entry in the database
	statement, _ := db.Prepare("INSERT INTO accounts VALUES ($1, $2, $3, $4, $5, $6, $7, $8);")
	statement.Exec(AccountNumber, OwnerName, hash, "none", TokValue, 0, 0, 0)
	NewAccount := Account{OwnerName, AccountNumber, 0, 0, 10}
	return NewAccount, nil
}
func CreateSubAccount(db *sql.DB, username string) (Account, error) {
	DualInfo(fmt.Sprintf("Creating sub account for user %s", username))
	accounts, err := db.Query("SELECT Password FROM accounts WHERE OwnerName = $1;", username)
	defer accounts.Close()
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
		return Account{}, err
	}
	var hash string
	for accounts.Next() {
		//Gather the hash
		err := accounts.Scan(&hash)
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
			return Account{}, err
		}
		break
	}
	accounts.Close()
	//Get token value
	var TokValue string
	TokValue, err = GetToken(db, username)
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
		return Account{}, err
	}
	//Generate new Account Number
	AccountNumber := genAccountNumber(db)
	//Create the account entry in the database
	var statement *sql.Stmt
	statement, err = db.Prepare("INSERT INTO accounts VALUES ($1, $2, $3, $4, $5, $6, $7, $8);")
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
		return Account{}, err
	}
	_, err = statement.Exec(AccountNumber, username, hash, "none", TokValue, 0, 0, 0)
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
		return Account{}, err
	}
	NewAccount := Account{username, AccountNumber, 0, 0, 0}
	return NewAccount, nil

	return Account{}, nil
}

//Verify an account is valid, given the username and password
func IsAccountValid(db *sql.DB, username string, password string) (bool, error) {
	//Query database for matching user accounts
	accounts, err := db.Query("SELECT Password FROM accounts WHERE OwnerName = $1;", username)
	defer accounts.Close()
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
		return false, err
	}
	var hash string
	for accounts.Next() {
		//Gather the hash
		err := accounts.Scan(&hash)
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
			return false, err
		}
		//rowCount = rowCount + 1
		break
	}
	//Compare the hash, and return appropriatly
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		accounts.Close()
		return true, nil
	}
	accounts.Close()
	return false, nil
}

//Remove all tokens in the cache that haven't been used in 3 hours
func CleanTokenCache() {
	for k, v := range TokenCache {
		timeElapsed := time.Now().Unix() - v.StoredAt
		if timeElapsed > 10800 {
			delete(TokenCache, k)
			DualDebug("Deleting old token from memory")
		}
	}
}
func GetToken(db *sql.DB, username string) (string, error) {
	//First check the token values stored in memory
	tokValueHolder, ok := TokenCache[username]
	var TokValue string
	if ok {
		DualDebug("Found Token in Memory, using it")
		tokValueHolder.StoredAt = time.Now().Unix()
		CleanTokenCache()
		return tokValueHolder.TokValue, nil
		//If the token isn't stored in memory, get it from the database
	} else {
		DualDebug("Didn't find Token in Memory, accessing database")
		accounts, err := db.Query("SELECT TokValue FROM accounts WHERE OwnerName = $1;", username)
		defer accounts.Close()
		if err != nil {
			DualDebug(fmt.Sprintf("%v", err))
			return "", err
		}
		for accounts.Next() {
			//Gather the Token
			err := accounts.Scan(&TokValue)
			if err != nil {
				DualWarning(fmt.Sprintf("%v", err))
				return "", err
			}
			CleanTokenCache()
			TokenCache[username] = TokenValueHolder{TokValue, time.Now().Unix()}
			DualDebug("Found Token Value!")
			break
		}
		accounts.Close()
	}
	return TokValue, nil
}
func GetAccounts(db *sql.DB, username string) ([]Account, error) {
	DualDebug("Got request to list accounts")
	//Query database for matching user accounts
	accounts, err := db.Query("SELECT AccountNumber, OwnerName, DcBalance, CcBalance, ArBalance FROM accounts WHERE OwnerName = $1;", username)
	defer accounts.Close()
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
		return []Account{Account{}}, err
	}
	//Initialize variables
	var account Account
	var Accounts []Account
	for accounts.Next() {
		//Gather the hash
		DualDebug("Found Account")
		err := accounts.Scan(&account.Number, &account.Username, &account.DcBal, &account.CcBal, &account.ArBal)
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
			return []Account{Account{}}, err
		}
		Accounts = append(Accounts, account)
	}
	DualDebug("Done looping over accounts")
	accounts.Close()
	return Accounts, nil
}
func ChangeToken(db *sql.DB, username string) error {
	//Generate a new token string
	TokValue, err := GenerateRandomString(32)
	if err != nil {
		DualErr(err)
	}
	//Set the token in the database for the user to the new token
	statement, _ := db.Prepare("UPDATE accounts SET TokValue = $1 WHERE OwnerName = $2")
	statement.Exec(TokValue, username)
	//Update the token cache with the user's new token
	TokenCache[username] = TokenValueHolder{TokValue, time.Now().Unix()}
	//Clear old tokens
	CleanTokenCache()
	return nil
}
func TransferFunc(db *sql.DB, transfer Transfer) (error, int){
	//Make sure the user isn't trying to transfer a negative number
	DualInfo(fmt.Sprintf("%s is beginning a transfer from %s to %s of amount %d %s", transfer.Username, transfer.FromAccount, transfer.ToAccount, transfer.Amount, transfer.Type))
	if transfer.Amount < 0 {
		DualNotice("User attempted to transfer negative number")
        DualInfo("Transfer Failed")
		return errors.New("Cannot transfer negative number. Nice Try ;)"), http.StatusBadRequest
	}
	//Getting the account to transfer from
	//Query db for matching account
	DualDebug("Accessing account database")
	accounts, err := db.Query("SELECT OwnerName, DcBalance, CcBalance, ArBalance FROM accounts WHERE AccountNumber = $1;", transfer.FromAccount)
	defer accounts.Close()
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
        DualInfo("Transfer Failed")
		return err, http.StatusInternalServerError
	}
	//Empty account object
	var fromAccount Account
	fromAccount.Number = transfer.FromAccount
	rowCount := 0
	for accounts.Next() {
		//Populate accout
		err := accounts.Scan(&fromAccount.Username, &fromAccount.DcBal, &fromAccount.CcBal, &fromAccount.ArBal)
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
            DualInfo("Transfer Failed")
            return err, http.StatusInternalServerError
		}
		rowCount = rowCount + 1
		break
	}
	DualDebug("Got database info")
	accounts.Close()
	if rowCount == 0 {
		DualDebug("Account not found")
        DualInfo("Transfer Failed")
		return errors.New("Account not found"), http.StatusBadRequest
	}
	//Repeat for to account
	DualDebug("Accessing account database")
	accounts, err = db.Query("SELECT OwnerName, DcBalance, CcBalance, ArBalance FROM accounts WHERE AccountNumber = $1;", transfer.ToAccount)
	defer accounts.Close()
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
        DualInfo("Transfer Failed")
		return err, http.StatusInternalServerError
	}
	//Empty account object
	var toAccount Account
	toAccount.Number = transfer.ToAccount
	rowCount = 0
	for accounts.Next() {
		//Populate accout
		err := accounts.Scan(&toAccount.Username, &toAccount.DcBal, &toAccount.CcBal, &toAccount.ArBal)
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
            DualInfo("Transfer Failed")
			return err, http.StatusInternalServerError
		}
		rowCount = rowCount + 1
		break
	}
	DualDebug("Got database info")
	accounts.Close()
	if rowCount == 0 {
		DualDebug("Account not found")
        DualInfo("Transfer Failed")
		return errors.New("Account not found"), http.StatusBadRequest
	}
	//Make sure the user owns the requested account to transfer from
	if fromAccount.Username != transfer.Username {
		DualDebug("Cannot Transfer from unowned account")
        DualInfo("Transfer Failed")
		return errors.New("Cannot Transfer From Unowned Account"), http.StatusForbidden
	}
	//Make sure the user has the balance to complete the transfer
	var statement1 *sql.Stmt
	var statement2 *sql.Stmt
	if transfer.Type == "ArBal" {
		if fromAccount.ArBal < transfer.Amount {
			DualDebug("Not enough balance to complete transfer")
            DualInfo("Transfer Failed")
			return errors.New("Not enough balance"), http.StatusBadRequest
		}
		statement1, err = db.Prepare("UPDATE accounts SET ArBalance = ArBalance - $1 WHERE AccountNumber = $2")
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
            DualInfo("Transfer Failed")
			return err, http.StatusInternalServerError
		}
		statement2, err = db.Prepare("UPDATE accounts SET ArBalance = ArBalance + $1 WHERE AccountNumber = $2")
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
            DualInfo("Transfer Failed")
			return err, http.StatusInternalServerError
		}
	} else if transfer.Type == "DcBal" {
		if fromAccount.DcBal < transfer.Amount {
			DualDebug("Not enough balance to complete transfer")
            DualInfo("Transfer Failed")
			return errors.New("Not enough balance"), http.StatusBadRequest
		}
		statement1, err = db.Prepare("UPDATE accounts SET DcBalance = DcBalance - $1 WHERE AccountNumber = $2")
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
            DualInfo("Transfer Failed")
			return err, http.StatusInternalServerError
		}
		statement2, err = db.Prepare("UPDATE accounts SET DcBalance = DcBalance + $1 WHERE AccountNumber = $2")
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
            DualInfo("Transfer Failed")
			return err, http.StatusInternalServerError
		}
	} else if transfer.Type == "CcBal" {
		if fromAccount.CcBal < transfer.Amount {
			DualDebug("Not enough balance to complete transfer")
            DualInfo("Transfer Failed")
			return errors.New("Not enough balance"), http.StatusBadRequest
		}
		statement1, err = db.Prepare("UPDATE accounts SET CcBalance = CcBalance - $1 WHERE AccountNumber = $2")
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
            DualInfo("Transfer Failed")
			return err, http.StatusInternalServerError
		}
		statement2, err = db.Prepare("UPDATE accounts SET CcBalance = CcBalance + $1 WHERE AccountNumber = $2")
		if err != nil {
			DualWarning(fmt.Sprintf("%v", err))
            DualInfo("Transfer Failed")
			return err, http.StatusInternalServerError
		}
	} else {
		DualDebug("Invalid Balance Type")
        DualInfo("Transfer Failed")
		return errors.New("Invalid Balance Type"), http.StatusBadRequest
	}
	DualDebug("Checks succeeded; beginning transfer")
	//Complete the transfer
	_, err = statement1.Exec(transfer.Amount, fromAccount.Number)
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
        DualInfo("Transfer Failed")
		return err, http.StatusInternalServerError
	}
	DualDebug("Checks succeeded; beginning transfer")
	_, err = statement2.Exec(transfer.Amount, toAccount.Number)
	if err != nil {
		DualWarning(fmt.Sprintf("%v", err))
        DualInfo("Transfer Failed")
		return err, http.StatusInternalServerError
	}
	DualInfo("Transfer Completed!")
	return nil, http.StatusOK
}
