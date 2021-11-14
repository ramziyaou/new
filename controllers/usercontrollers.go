package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"rest-go-demo/database"
	"rest-go-demo/entity"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var (
	InvalidInput      error = errors.New("Invalid username or password")
	DuplicateID       error = errors.New("ID exists")
	DuplicateUsername error = errors.New("Username exists")
	Commit            error = errors.New("Transaction's been committed")
	ZeroID            error = errors.New("ID can't be 0")
)

// Creates and returns new user based on passed arguments
func NewUser(strId, username, password string) (entity.User, error) {
	var user entity.User

	if username == "" || password == "" {
		return user, InvalidInput
	}

	// Convert ID into int
	id, err := strconv.Atoi(strId)
	if err != nil {
		return user, err
	}
	if id == 0 {
		return user, ZeroID
	}

	// Check if ID exists
	if err := database.Connector.Where("id = ?", id).First(&user).Error; err == nil {
		return user, DuplicateID
	}

	// Check if username exists
	if err := database.Connector.Where("username = ?", username).First(&user).Error; err == nil {
		return user, DuplicateUsername
	}

	// If all checks above passed
	return entity.User{
		ID:       id,
		Username: username,
		Password: password,
	}, nil
}

//Get user info
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from path
	vars := mux.Vars(r)
	strId := vars["id"]

	// Incorrect ID format
	id, err := strconv.Atoi(strId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "ID is not integer")
		return
	}

	// Check if user exists
	var user entity.User
	if err := database.Connector.Where("id = ?", id).First(&user).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "User with ID %d not found\n", id)
		return
	}

	// Display user record if found
	fmt.Fprintf(w, "ID: %d, username: %s\nWallets:\n", user.ID, user.Username)
	if len(user.Wallets) == 0 {
		fmt.Fprintln(w, "empty")
	}
	for _, v := range strings.Split(user.Wallets, " ") {
		fmt.Fprintln(w, v)
	}
}

//Save user
func SaveUser(w http.ResponseWriter, r *http.Request) {
	// Get ID from path
	vars := mux.Vars(r)
	strId := vars["id"]

	// Get username and password from params
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := NewUser(strId, username, password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, err.Error())
		return
	}

	// If no error
	database.Connector.Create(user)
	fmt.Fprintln(w, "Success")
}

// Get wallet amount
func GetWallet(w http.ResponseWriter, r *http.Request) {
	var user entity.User
	// Get username from value passed in request context upon authentication
	username := r.Context().Value("username").(string)
	database.Connector.Where("username = ?", username).First(&user)

	// Get wallet name from path
	vars := mux.Vars(r)
	wallet := vars["name"]

	// Check if user has any wallets
	if user.Wallets == "" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w, "You have no cryptowallets")
		return
	}
	// Get wallet info if everything above ok
	var v entity.CryptoWallet
	if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&v).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Couldn't find wallet %s under your username", wallet)
		return
	}

	fmt.Fprintf(w, "Wallet %s, amount: %d", wallet, v.Amount)
}

// Create new wallet
func SaveWallet(w http.ResponseWriter, r *http.Request) {
	// Get wallet name from path
	vars := mux.Vars(r)
	wallet := vars["name"]
	var user entity.User
	// Get username from value passed in request context upon authentication
	username := r.Context().Value("username").(string)
	if err := database.Connector.Where("username = ?", username).First(&user).Error; err != nil {
		log.Println("SaveWallet unexpected error:", err)
		return
	}

	// Check if wallet already exists
	for _, name := range strings.Split(user.Wallets, " ") {
		if name == wallet {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Wallet under name %s already exists", wallet)
			return
		}
	}

	// Create wallet under given name if everything above ok
	user.Wallets += wallet + " "

	// Save updated user changes
	if err := database.Connector.Save(&user).Error; err != nil {
		log.Println("SaveWallet unexpected error: ", err)
		return
	}

	// Create a wallet in DB
	v := entity.NewWallet(wallet)
	v.Username = username
	if err := database.Connector.Create(&v).Error; err != nil {
		log.Println("SaveWallet unexpected error:", err)
	}
	// Save wallet status in DB
	ss := entity.NewStartStop(username, wallet)
	if err := database.Connector.Create(&ss).Error; err != nil {
		log.Println("SaveWallet unexpected error:", err)
		return
	}
	fmt.Fprintln(w, "Success")
}

// Start cryptomining
func StartMining(w http.ResponseWriter, r *http.Request) {
	// Get username from request context
	if err := dbTransaction(database.Connector, w, r); err != nil && err != Commit {
		fmt.Println(err)
	} else if err == Commit || err == nil {
		fmt.Println("Transaction completed")
	}
}

// Transaction for mining
func dbTransaction(db *gorm.DB, w http.ResponseWriter, r *http.Request) error {
	// Begin transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	// Get username from request context
	username := r.Context().Value("username")

	// Get wallet name from path
	vars := mux.Vars(r)
	wallet := vars["name"]

	// Get wallet from DB
	var v entity.CryptoWallet
	if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&v).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Couldn't find wallet %s under your username\n", wallet)
		return err
	}
	// Check wallet status
	var ss entity.StartStopCheck
	// Check below is redundant
	if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&ss).Error; err != nil {
		log.Println("StartMining unexpected error:", err)
		return err
	}
	if ss.Start {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Mining not stopped yet")
		return fmt.Errorf("Mining not stopped yet")
	}

	// Update wallet status to mining started
	if err := database.Connector.Model(&ss).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"start": true}).Error; err != nil {
		fmt.Println("dbTransaction unexpected error:", err)
		return err
	}
	v.RLock()
	defer v.RUnlock()
	fmt.Fprintf(w, "Starting mining, current amount: %d\n", v.Amount)

	// Wait for stop instructions, otherwise keep mining
	go func(v *entity.CryptoWallet) {
		wg := &sync.WaitGroup{}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		wg.Add(1)
		go func(ctx context.Context, v *entity.CryptoWallet, wg *sync.WaitGroup) {
			for {
				v.Mine()
				select {
				case <-ctx.Done():
					if err := database.Connector.Model(&v).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"amount": v.Amount}).Error; err != nil {
						fmt.Println("wallet amount update error", err) // nil
						return
					}
					wg.Done()
					log.Println("Mining actually stopped for wallet", wallet)
					return
				default:
				}
			}
		}(ctx, v, wg)
		for !ss.Stop {
			time.Sleep(time.Second * 10)
			// Check if stop has been updated by StopMining function
			if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&ss).Error; err != nil {
				w.WriteHeader(http.StatusNotFound)
				log.Println("StartMining:", err)
				return
			}
		}
		// Cancel context if stop signal's been received
		cancel()
		log.Println("Canceling context")
		wg.Wait()
		if err := database.Connector.Model(&ss).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"start": false, "stop": false}).Error; err != nil {
			log.Println("StartMining:", err)
			return
		}
		// Finish our goroutine
		return
	}(&v)
	// Check if rollback took place (happens when we call StartMining consecutively without StopMining)
	if tx.Error != Commit {
		tx.Commit().Error = Commit
		return Commit
	} else {
		tx.Rollback()
	}
	return nil
}

// Stop cryptomining
func StopMining(w http.ResponseWriter, r *http.Request) {
	// Get username from request context
	username := r.Context().Value("username")

	// Get wallet name from path
	vars := mux.Vars(r)
	wallet := vars["name"]

	// Check if mining's started
	var ss entity.StartStopCheck
	if err := database.Connector.Where("username = ?", username).Where("name = ?", wallet).First(&ss).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Couldn't find wallet %s under your username\n", wallet)
		return
	}

	if !ss.Start {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Can't stop mining that hasn't started or wallet doesn't exist\n")
		return
	}

	log.Println("Closing wallet", wallet)
	// Send signal to StartMining method to stop mining
	if err := database.Connector.Model(&ss).Where("username = ?", username).Where("name = ?", wallet).Updates(map[string]interface{}{"stop": true}).Error; err != nil {
		fmt.Println("StopMining unexpected error:", err.Error())
		return
	}
	// Можно убрать этот слип, он предотвращает юзера нажимать на стоп несколько раз. Ждать окончания майнинга придется и без этого слипа.
	time.Sleep(time.Second * 10)
	fmt.Fprintln(w, "Stopping mining for wallet, takes a while")
}
