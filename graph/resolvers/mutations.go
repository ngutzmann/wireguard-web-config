package resolvers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/google/uuid"
	"github.com/ngutzmann/wireguard-web-config/graph/model"
)

// CreatePeer - Add a peer to the server
func CreatePeer(ctx context.Context, db *sql.DB, input model.NewPeer) (*model.Peer, error) {
	id := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO peers (id, user_f_name, user_l_name, hostname, public_key, allowed_ip) VALUES ($1, $2, $3, $4, $5, $6)", id, input.UserFName, input.UserLName, input.Hostname, input.PublicKey, input.AllowedIP)

	if err != nil {
		log.Println("Failed to create peer entry in DB:", err)
		return nil, err
	}

	peer := &model.Peer{
		ID:        id.String(),
		UserFName: input.UserFName,
		UserLName: input.UserLName,
		Hostname:  input.Hostname,
		PublicKey: input.PublicKey,
		AllowedIP: input.AllowedIP,
	}

	cmd := fmt.Sprintf("wg set %s peer %s allowed-ips %s/32", getWGInterface(), input.PublicKey, input.AllowedIP)

	cmdCtx := exec.CommandContext(ctx, "sudo", strings.Split(cmd, " ")...)
	err = cmdCtx.Run()
	if err != nil {
		log.Printf("Failed to add WG peer: key: %s IP: %s\n", input.PublicKey, input.AllowedIP)
		return nil, err
	}

	cmdText := strings.Split("wg-quick save %s", " ")
	cmdCtx = exec.CommandContext(ctx, "sudo", cmdText...)
	err = cmdCtx.Run()

	if err != nil {
		log.Printf("Failed to save WG config after adding peer: key: %s IP: %s\n", input.PublicKey, input.AllowedIP)
		return nil, err
	}
	return peer, nil
}
