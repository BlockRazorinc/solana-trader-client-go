package examples

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	pb "github.com/BlockRazorinc/solana-trader-client-go/pb/serverpb"
	"github.com/gagliardetto/solana-go"
)

type tipAccountsResponse struct {
	TipAccounts []string `json:"tipAccounts"`
}

func HTTPTipAccount(ctx context.Context, client *http.Client, endpoint string) (solana.PublicKey, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return solana.PublicKey{}, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("get tip accounts: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return solana.PublicKey{}, fmt.Errorf("get tip accounts: unexpected HTTP status %s", resp.Status)
	}

	var result tipAccountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return solana.PublicKey{}, fmt.Errorf("decode tip accounts: %w", err)
	}
	return RandomTipAccount(result.TipAccounts)
}

func GRPCTipAccount(ctx context.Context, client pb.ServerClient) (solana.PublicKey, error) {
	resp, err := client.GetTipAccounts(ctx, &pb.GetTipAccountsRequest{})
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("get tip accounts: %w", err)
	}
	return RandomTipAccount(resp.GetTipAccounts())
}

func RandomTipAccount(accounts []string) (solana.PublicKey, error) {
	if len(accounts) == 0 {
		return solana.PublicKey{}, fmt.Errorf("get tip accounts: empty account list")
	}
	account, err := solana.PublicKeyFromBase58(accounts[rand.Intn(len(accounts))])
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("invalid tip account: %w", err)
	}
	return account, nil
}
