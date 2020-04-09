package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/ngutzmann/wireguard-web-config/graph/generated"
	"github.com/ngutzmann/wireguard-web-config/graph/model"
)

type wgPeer struct {
	PublicKey       string
	Endpoint        *string
	AllowedIP       string
	LatestHandshake *int
	TransferRx      *int
	TransferTx      *int
}

func getWGInterface() string {
	return os.Getenv("WG_INTERFACE")
}

func readWgConfig(ctx context.Context) (map[string]wgPeer, error) {

	cmdText := strings.Split(fmt.Sprintf("wg show %s dump", getWGInterface()), " ")
	cmd := exec.CommandContext(ctx, "sudo", cmdText...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()
	if err != nil {
		log.Println("Failed to read wg config.", err)
		return nil, err
	}
	outStr := string(stdout.Bytes())

	scanner := bufio.NewScanner(strings.NewReader(outStr))
	count := 0
	peerMap := make(map[string]wgPeer)
	for scanner.Scan() {
		line := scanner.Text()
		count++
		if count == 1 {
			// Skip the first line
			continue
		}
		tokens := strings.Split(line, "\t")
		if len(tokens) != 8 {
			log.Printf("Only found %v tokens.\n", len(tokens))
			continue
		}

		transferRx, _ := strconv.Atoi(tokens[5])
		transferTx, _ := strconv.Atoi(tokens[6])
		latestHandshake, _ := strconv.Atoi(tokens[4])
		peerMap[tokens[0]] = wgPeer{
			PublicKey:       tokens[0],
			Endpoint:        &tokens[2],
			AllowedIP:       tokens[3],
			LatestHandshake: &latestHandshake,
			TransferRx:      &transferRx,
			TransferTx:      &transferTx,
		}

	}
	err = scanner.Err()
	return peerMap, err
}

func (r *mutationResolver) CreatePeer(ctx context.Context, input model.NewPeer) (*model.Peer, error) {
	db := r.Resolver.DB
	id := uuid.New()
	_, err := db.ExecContext(ctx, "INSERT INTO peers (id, name, public_key, allowed_ip) VALUES ($1, $2, $3, $4)", id, input.Name, input.PublicKey, input.AllowedIP)

	if err != nil {
		log.Println("Failed to create peer entry in DB:", err)
		return nil, err
	}

	peer := &model.Peer{
		ID:        id.String(),
		Name:      input.Name,
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

func (r *queryResolver) Peers(ctx context.Context) ([]*model.Peer, error) {
	var peers []*model.Peer
	db := r.Resolver.DB
	rows, err := db.QueryContext(ctx, "SELECT id, name, public_key, allowed_ip FROM peers")
	if err != nil {
		log.Println("Error querying peers in DB:", err)
	}
	defer rows.Close()
	var (
		id        string
		name      string
		publicKey string
		allowedIP string
	)

	peerMap, err := readWgConfig(ctx)
	for rows.Next() {
		err := rows.Scan(&id, &name, &publicKey, &allowedIP)
		if err != nil {
			log.Println("Error reading row from `peers` table: ", err)
		}

		wgPeer := peerMap[publicKey]
		peer := &model.Peer{
			ID:              id,
			Name:            name,
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

func (r *queryResolver) Peer(ctx context.Context, id string) (*model.Peer, error) {
	var peer *model.Peer
	db := r.Resolver.DB
	rows, err := db.QueryContext(ctx, "SELECT id, name, public_key, allowed_ip FROM peers WHERE id=$1", id)
	if err != nil {
		log.Println("Error querying peers in DB:", err)
	}
	defer rows.Close()
	var (
		name      string
		publicKey string
		allowedIP string
	)

	peerMap, err := readWgConfig(ctx)
	for rows.Next() {
		err := rows.Scan(&id, &name, &publicKey, &allowedIP)
		if err != nil {
			log.Println("Error reading row from `peers` table: ", err)
		}

		wgPeer := peerMap[publicKey]
		peer = &model.Peer{
			ID:              id,
			Name:            name,
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

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
