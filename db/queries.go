package db

import (
	"database/sql"
	"time"
	//"time"

	"github.com/google/uuid"
)

func InsertPeer (db *sql.DB, peerID, name, IP string) (uuid.UUID, error){
	var id uuid.UUID

	query := 
		`INSERT INTO peers(peer_id, name, ip_address, is_online)
		VALUES ($1, $2, $3, true)
		RETURNING id;
		`

	err := db.QueryRow(query, peerID, name, IP).Scan(&id)

	return id, err
}

func MarkPeerOffline(db *sql.DB, peerID string, lastSeen time.Time) error {
	query := `
		UPDATE peers
		SET is_online = false, 
		    last_seen = $2
		WHERE peer_id = $1;
	`
	_, err := db.Exec(query, peerID, lastSeen)
	return err
}
