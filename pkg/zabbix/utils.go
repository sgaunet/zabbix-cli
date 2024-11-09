package zabbix

import (
	"crypto/rand"
	"math/big"
)

const maxUniqueID = 10000

// generateUniqueID generates a unique ID for the JSON-RPC request
// It returns a random integer
func generateUniqueID() int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(maxUniqueID))
	if err != nil {
		panic(err)
	}
	return int(nBig.Uint64()) //nolint:gosec
}
