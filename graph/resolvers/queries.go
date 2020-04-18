package resolvers

import (
	"context"
	"database/sql"
	"log"

	"github.com/ngutzmann/wireguard-web-config/graph/model"
)

// GetPeers - get all the peers
func GetPeers(ctx context.Context, db *sql.DB) ([]*model.Peer, error) {
	var peers []*model.Peer
	rows, err := db.QueryContext(ctx, "SELECT id, user_f_name, user_l_name, hostname, public_key, allowed_ip FROM peers")
	if err != nil {
		log.Println("Error querying peers in DB:", err)
		return nil, err
	}
	defer rows.Close()
	var (
		id        string
		userFName *string
		userLName *string
		hostname  string
		publicKey string
		allowedIP string
	)

	peerMap, err := readWgConfig(ctx)
	for rows.Next() {
		err := rows.Scan(&id, &userFName, &userLName, &hostname, &publicKey, &allowedIP)
		if err != nil {
			log.Println("Error reading row from `peers` table: ", err)
			return nil, err
		}

		wgPeer := peerMap[publicKey]
		peer := &model.Peer{
			ID:              id,
			UserFName:       userFName,
			UserLName:       userLName,
			Hostname:        hostname,
			PublicKey:       publicKey,
			AllowedIP:       allowedIP,
			Endpoint:        wgPeer.Endpoint,
			LatestHandshake: wgPeer.LatestHandshake,
			TransferRxBytes: wgPeer.TransferRx,
			TransferTxBytes: wgPeer.TransferTx,
		}
		peers = append(peers, peer)
	}
	return peers, nil
}

// GetPeer - Get a single peer based on the ID
func GetPeer(ctx context.Context, db *sql.DB, id string) (*model.Peer, error) {
	var peer *model.Peer
	rows, err := db.QueryContext(ctx, "SELECT id, user_f_name, user_l_name, hostname, public_key, allowed_ip FROM peers WHERE id=$1", id)
	if err != nil {
		log.Println("Error querying peers in DB:", err)
		return nil, err
	}
	defer rows.Close()
	var (
		userFName *string
		userLName *string
		hostname  string
		publicKey string
		allowedIP string
	)

	peerMap, err := readWgConfig(ctx)
	for rows.Next() {
		err := rows.Scan(&id, &userFName, &userLName, &hostname, &publicKey, &allowedIP)
		if err != nil {
			log.Println("Error reading row from `peers` table: ", err)
			return nil, err
		}

		wgPeer := peerMap[publicKey]
		peer = &model.Peer{
			ID:              id,
			UserFName:       userFName,
			UserLName:       userLName,
			Hostname:        hostname,
			PublicKey:       publicKey,
			AllowedIP:       allowedIP,
			Endpoint:        wgPeer.Endpoint,
			LatestHandshake: wgPeer.LatestHandshake,
			TransferRxBytes: wgPeer.TransferRx,
			TransferTxBytes: wgPeer.TransferTx,
		}
	}
	return peer, nil
}
