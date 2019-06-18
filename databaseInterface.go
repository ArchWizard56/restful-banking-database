package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
    crand "crypto/rand"
	"strconv"
	"time"
    "fmt"
    "encoding/base64"
)

type Account struct {
    Username string
    Number string
    CcBal int
    DcBal int
    ArBal int
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
	db, err := sql.Open("sqlite3", DBPath)
	if err != nil {
		DualErr(err)
	}
	// Initialize Default table consisting of:
	// AccountNumber    OwnerName   Password    DiscordName     TokValue    DcBalance   CcBalance   ArBalance
	statement, _ := db.Prepare("CREATE TABLE IF NOT EXISTS accounts (AccountNumber TEXT PRIMARY KEY, OwnerName TEXT, Password TEXT, DiscordName TEXT, TokValue TEXT, DcBalance INTEGER, CcBalance INTEGER, ArBalance INTEGER)")
	statement.Exec()
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
	statement.Exec(AccountNumber, OwnerName, hash, "none", TokValue, 0, 0, 10)
    NewAccount := Account{OwnerName, AccountNumber, 0,0,10}
    return  NewAccount, nil
}
//Verify an account is valid, given the username and password
func IsAccountValid (db *sql.DB, username string, password string) (bool, error) {
    //Query database for matching user accounts
	accounts, err := db.Query("SELECT Password FROM accounts WHERE OwnerName = $1;",username)
    if err != nil {
        DualWarning(fmt.Sprintf("%v", err))
        return false, err
    }
    var hash string
    for accounts.Next(){
    //Gather the hash
    err := accounts.Scan(&hash)
    if err != nil {
        DualWarning(fmt.Sprintf("%v", err))
        return false, err
    }
    //Compare the hash, and return appropriatly
    err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err == nil {
        return true, nil
    }
    return false, nil
}
    return false, nil
}
//Remove all tokens in the cache that haven't been used in 3 hours
func CleanTokenCache () {
    for k, v := range TokenCache {
        timeElapsed := time.Now().Unix() - v.StoredAt
        if timeElapsed > 10800 {
            delete(TokenCache, k)
            DualDebug("Deleting old token from memory")
        }
    }
}
func GetToken (db *sql.DB, username string) (string, error) {
    //First check the token values stored in memory
    tokValueHolder, ok := TokenCache[username]
    if ok {
        DualDebug("Found Token in Memory, using it")
        tokValueHolder.StoredAt = time.Now().Unix()
        CleanTokenCache()
        return tokValueHolder.TokValue, nil
    //If the token isn't stored in memory, get it from the database
    } else    {
    DualDebug("Didn't find Token in Memory, accessing database")
	accounts, err := db.Query("SELECT TokValue FROM accounts WHERE OwnerName = $1;",username)
    if err != nil {
        DualDebug(fmt.Sprintf("%v", err))
        return "", err
    }
    var TokValue string
    for accounts.Next(){
    //Gather the Token
    err := accounts.Scan(&TokValue)
    if err != nil {
        DualWarning(fmt.Sprintf("%v", err))
        return "", err 
    }
    CleanTokenCache()
    TokenCache[username] = TokenValueHolder{TokValue,time.Now().Unix()}
    DualDebug("Found Token Value!")
    return TokValue, nil
}
}
return "", nil
}
func GetAccounts (db *sql.DB, username string) ([]Account, error) {
    //Query database for matching user accounts
	accounts, err := db.Query("SELECT * FROM accounts WHERE OwnerName = $1;",username)
    if err != nil {
        DualWarning(fmt.Sprintf("%v", err))
        return []Account{Account{}}, err
    }
    var account Account
    var Accounts []Account
    var dummy []byte
    var dummyTokValue []byte
    var dummyName []byte
    for accounts.Next(){
    //Gather the hash
    err := accounts.Scan(&account.Number, &account.Username, &dummyTokValue, &dummy, &dummyName, &account.CcBal, &account.DcBal, &account.ArBal)
    if err != nil {
        DualWarning(fmt.Sprintf("%v", err))
        return []Account{Account{}}, err
    }
    Accounts = append(Accounts, account)
}
return Accounts, nil
}
func ChangeToken (db *sql.DB, username string) (error) {
    TokValue, err := GenerateRandomString(32)
		if err != nil {
			DualErr(err)
		}
	statement, _ := db.Prepare("UPDATE accounts SET TokValue = $1 WHERE OwnerName = $2")
	statement.Exec(TokValue,username)
    TokenCache[username] = TokenValueHolder{TokValue,time.Now().Unix()}
    CleanTokenCache()
    return nil
}
