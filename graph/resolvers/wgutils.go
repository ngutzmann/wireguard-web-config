package resolvers

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
)

// wgPeer - Type that corresponds to the wg show <interface> dump command
type wgPeer struct {
	PublicKey       string
	Endpoint        *string
	AllowedIP       string
	LatestHandshake *int
	TransferRx      *int
	TransferTx      *int
}

// getWGInterface - Get the Wireguard Interface used right now
func getWGInterface() string {
	return os.Getenv("WG_INTERFACE")
}

// readWgConfig - Read Wireguard config and stats
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
