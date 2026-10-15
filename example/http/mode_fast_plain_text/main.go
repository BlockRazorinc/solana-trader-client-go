package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	examples "github.com/BlockRazorinc/solana-trader-client-go/example"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	// Example endpoints only. Always use the current endpoints published on the BlockRazor website.
	httpEndpoint        = "http://frankfurt.solana.blockrazor.xyz:443/v2/sendTransaction"
	healthEndpoint      = "http://frankfurt.solana.blockrazor.xyz:443/health"
	tipAccountsEndpoint = "http://frankfurt.solana.blockrazor.xyz:443/getTipAccounts"
	mainNetRPC          = ""
	authKey             = ""
	privateKey          = ""
	publicKey           = ""
	amount              = 200_000
	tipAmount           = 100_000
	mode                = "fast"
)

var httpClient = &http.Client{Timeout: 10 * time.Second}

func main() {
	if err := pingHealth(); err != nil {
		fmt.Printf("health check failed: %v\n", err)
	}
	go func() {
		for {
			if err := pingHealth(); err != nil {
				fmt.Printf("health check failed: %v\n", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()
	if err := sendTx(); err != nil {
		fmt.Printf("send tx failed: %v\n", err)
	}
}

func pingHealth() error {
	req, err := http.NewRequest(http.MethodPost, healthEndpoint, bytes.NewBufferString("{}"))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", authKey)
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(io.Discard, resp.Body)
	return err
}

func sendTx() error {
	account, err := solana.WalletFromPrivateKeyBase58(privateKey)
	if err != nil {
		return err
	}
	receiver := solana.MustPublicKeyFromBase58(publicKey)
	tipAccount, err := examples.HTTPTipAccount(context.Background(), httpClient, tipAccountsEndpoint)
	if err != nil {
		return err
	}
	blockhash, err := rpc.New(mainNetRPC).GetLatestBlockhash(context.Background(), rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf("get blockhash: %w", err)
	}
	tx, err := solana.NewTransaction(
		[]solana.Instruction{
			system.NewTransferInstruction(tipAmount, account.PublicKey(), tipAccount).Build(),
			system.NewTransferInstruction(amount, account.PublicKey(), receiver).Build(),
		},
		blockhash.Value.Blockhash,
		solana.TransactionPayer(account.PublicKey()),
	)
	if err != nil {
		return err
	}
	if _, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if account.PublicKey().Equals(key) {
			return &account.PrivateKey
		}
		return nil
	}); err != nil {
		return err
	}
	encoded, err := tx.ToBase64()
	if err != nil {
		return err
	}
	u, err := url.Parse(httpEndpoint)
	if err != nil {
		return err
	}
	query := u.Query()
	query.Set("auth", authKey)
	query.Set("mode", mode)
	u.RawQuery = query.Encode()
	req, err := http.NewRequest(http.MethodPost, u.String(), bytes.NewBufferString(encoded))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("send transaction failed: status=%s body=%s", resp.Status, strings.TrimSpace(string(bodyBytes)))
	}
	signature := strings.TrimSpace(string(bodyBytes))
	if signature == "" {
		return fmt.Errorf("send transaction returned an empty signature")
	}
	fmt.Printf("[send tx] signature: %s\n", signature)
	return nil
}
