package db

import (
	"database/sql"
	//"time"

	"github.com/google/uuid"
)

func newPeer(db *sql.DB, peerID, name, IP string) (uuid.UUID, error){
	var id uuid.UUID

	query := 
		`INSERT INTO peers(peer_id, name, ip_address, is_online)
		VALUES ($1, $2, $3, true)
		RETURNING id;
		`

	err := db.QueryRow(query, peerID, name, IP).Scan(&id)

	return id, err
}