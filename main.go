package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type State struct {
	Name    string `json:"name"`
	Balance uint   `json:"balance"`
}

var symbols = map[string]uint{
	"♠": 4,
	"♥": 7,
	"♦": 12,
	"♣": 20,
}

// multipliers for each symbol line
var multipliers = map[string]uint{
	"♠": 20,
	"♥": 10,
	"♦": 5,
	"♣": 2,
}

// In memory db
var (
	name    = "Anonymous Player"
	balance = uint(200)
	slotMu  sync.Mutex
)

func handleState(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		slotMu.Lock()

		var response = map[string]interface{}{
			"name":        name,
			"balance":     balance,
			"symbols":     symbols,
			"multipliers": multipliers,
		}

		json.NewEncoder(w).Encode(response)
		slotMu.Unlock()

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleName(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		slotMu.Lock()
		json.NewEncoder(w).Encode(map[string]interface{}{"name": name})
		slotMu.Unlock()

	case http.MethodPost:
		var newName struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&newName); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
			return
		}

		if newName.Name == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Name is required"})
			return
		}

		slotMu.Lock()
		name = newName.Name
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"name": name})
		slotMu.Unlock()

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func handleBet(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var newBet struct {
			Bet uint `json:"bet"`
		}

		if err := json.NewDecoder(r.Body).Decode(&newBet); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Invalid request body"})
			return
		}

		if newBet.Bet > uint(balance) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "You cannot bet more than your balance"})
			return
		}

		if newBet.Bet == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"message": "Thanks for playing!"})
			return
		}

		slotMu.Lock()

		symbolArray := GenerateSymbolArray(symbols)

		balance -= newBet.Bet

		spin := GetSpin(symbolArray, 3, 3)

		// Check for winning combinations and update balance accordingly
		winningLines := CheckWin(spin, multipliers)

		winResults := []map[string]interface{}{}

		for _, multiplier := range winningLines {
			winAmount := newBet.Bet * multiplier
			balance += winAmount

			if multiplier > 0 {
				winResults = append(winResults, map[string]interface{}{
					"multiplier": multiplier,
					"winAmount":  winAmount,
				})
			}
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"balance": balance, "spin": spin, "winResults": winResults})
		slotMu.Unlock()

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/", handleState)
	http.HandleFunc("/name", handleName)
	http.HandleFunc("/bet", handleBet)

	fmt.Println("Server running on :8080...")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
