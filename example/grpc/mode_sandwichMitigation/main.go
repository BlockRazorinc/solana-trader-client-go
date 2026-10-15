package main

import (
	"context"
	"fmt"
	"time"

	examples "github.com/BlockRazorinc/solana-trader-client-go/example"
	pb "github.com/BlockRazorinc/solana-trader-client-go/pb/serverpb"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// This example endpoint may change. Always use the current endpoint published on the BlockRazor website.
	blzRelayEndpoint = "frankfurt.solana-grpc.blockrazor.xyz:80"
	// replace your solana rpc endpoint
	mainNetRPC = ""
	// replace your authKey
	authKey = ""
	// relace your private key(base58)
	privateKey = ""
	// publicKey(base58)
	publicKey = ""
	// transfer amount
	amount = 200_000
	// send mode
	mode = "sandwichMitigation"
	// safeWindow
	safeWindow = 5
	// revertProtection
	revertProtection = false
	// tip amount
	tipAmount = 100_000
)

func main() {
	var err error
	account, err := solana.WalletFromPrivateKeyBase58(privateKey)
	receivePub := solana.MustPublicKeyFromBase58(publicKey)

	// setup grpc connect
	conn, err := grpc.NewClient(blzRelayEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&Authentication{authKey}),
	)
	if err != nil {
		panic(fmt.Sprintf("connect error: %v", err))
	}

	// use the Gateway client connection interface
	client := pb.NewServerClient(conn)

	// Pre-warm: perform an initial health check to establish the grpc connection
	err = pingHealth(client)
	if err != nil {
		fmt.Printf("health check failed: %v\n", err)
	}
	// Start a background goroutine to periodically send /health requests
	// For low-frequency users, this keeps the grpc connection alive (warm)
	go func() {
		for {
			err := pingHealth(client)
			if err != nil {
				fmt.Printf("health check failed: %v\n", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()

	// new rpc client and get latest block hash
	rpcClient := rpc.New(mainNetRPC)
	blockhash, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(fmt.Sprintf("[get latest block hash] error: %v", err))
	}

	tipAccount, err := examples.GRPCTipAccount(context.Background(), client)
	if err != nil {
		panic(fmt.Sprintf("[get tip accounts] error: %v", err))
	}

	// construct instruction
	transferIx := system.NewTransferInstruction(amount, account.PublicKey(), receivePub).Build()
	tipIx := system.NewTransferInstruction(tipAmount, account.PublicKey(), tipAccount).Build()

	// construct transation, replace your transation
	tx, err := solana.NewTransaction(
		[]solana.Instruction{tipIx, transferIx},
		blockhash.Value.Blockhash,
		solana.TransactionPayer(account.PublicKey()),
	)
	if err != nil {
		panic(fmt.Sprintf("new tx error: %v", err))
	}

	// transaction sign
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if account.PublicKey().Equals(key) {
				return &account.PrivateKey
			}
			return nil
		},
	)
	if err != nil {
		panic(fmt.Sprintf("sign tx error: %v", err))
	}

	txBase64, _ := tx.ToBase64()
	sendRes, err := client.SendTransaction(context.TODO(), &pb.SendRequest{
		Transaction:      txBase64,
		Mode:             mode,
		SafeWindow:       safeWindow,
		RevertProtection: revertProtection,
	})
	if err != nil {
		panic(fmt.Sprintf("[send tx] error: %v", err))
	}

	fmt.Printf("[send tx] response: %+v \n", sendRes)
	return
}

type Authentication struct {
	apiKey string
}

func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"apikey": a.apiKey}, nil
}

func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

func pingHealth(client pb.ServerClient) error {
	// grpc request warmup
	healthRes, err := client.GetHealth(context.Background(), &pb.HealthRequest{})
	if err != nil {
		return err
	}
	fmt.Printf("[health] response: %+v \n", healthRes.Status)
	return err
}
