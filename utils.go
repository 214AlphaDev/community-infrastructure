package community_infrastructure

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/jinzhu/gorm"
)

var randomTestDB = func() (*gorm.DB, error) {

	pathBytes := make([]byte, 30)
	_, err := rand.Read(pathBytes)
	if err != nil {
		panic(err)
	}

	return gorm.Open("sqlite3", "/tmp/"+hex.EncodeToString(pathBytes))

}
