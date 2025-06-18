package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
)

func CreateSession(db *sql.DB, userID int) (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		fmt.Println("Random generation failed:", err)
		return "", err
	}
	token := hex.EncodeToString(bytes)

	_, err = db.Exec(`
		INSERT INTO sessions (user_id, token, created_at, deleted_at)
		VALUES ($1, $2, $3, $4)`,
		userID, token, time.Now(), time.Now().Add(7*24*time.Hour),
	)
	if err != nil {
		fmt.Println("Session insert failed:", err)
		return "", err
	}

	return token, nil
}
