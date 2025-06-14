package utils

import (
	"fmt"
	"math/rand"
	"transaction-service/data"
)

// simulate transaction execution that return boolean for success or failed
func  ExecuteTransaction(transaction *data.Transaction) error {
	transaction.Status = rand.Intn(2)
	
	if transaction.Status == 0 {
		return fmt.Errorf("transaction failed: unsufficient balance")
	}

	return nil
}