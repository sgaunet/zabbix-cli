package zabbix

import (
	"math/rand/v2"
)

const maxUniqueID = 10000

// generateUniqueID generates a unique ID for the JSON-RPC request
// It returns a random integer
func generateUniqueID() int {
	return rand.IntN(maxUniqueID)
}
