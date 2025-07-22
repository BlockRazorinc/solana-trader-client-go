package main

import (
	"context"
	"fmt"
	"math/rand"

	pb "github.com/BlockRazorinc/solana-trader-client-go/pb/serverpb"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	// BlockRazor relay endpoint address
	blzRelayEndpoint = "frankfurt.solana-grpc.blockrazor.xyz:80"
	// replace your solana rpc endpoint
	mainNetRPC = ""
	// replace your authKey
	authKey = ""
	// relace your private key
	privateKey = ""
	// publicKey
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
	tipAmount = 1_000_000
)

var tipAccounts = []string{
	"Gywj98ophM7GmkDdaWs4isqZnDdFCW7B46TXmKfvyqSm",
	"FjmZZrFvhnqqb9ThCuMVnENaM3JGVuGWNyCAxRJcFpg9",
	"6No2i3aawzHsjtThw81iq1EXPJN6rh8eSJCLaYZfKDTG",
	"A9cWowVAiHe9pJfKAj3TJiN9VpbzMUq6E4kEvf5mUT22",
	"68Pwb4jS7eZATjDfhmTXgRJjCiZmw1L7Huy4HNpnxJ3o",
	"4ABhJh5rZPjv63RBJBuyWzBK3g9gWMUQdTZP2kiW31V9",
	"B2M4NG5eyZp5SBQrSdtemzk5TqVuaWGQnowGaCBt8GyM",
	"5jA59cXMKQqZAVdtopv8q3yyw9SYfiE3vUCbt7p8MfVf",
	"5YktoWygr1Bp9wiS1xtMtUki1PeYuuzuCF98tqwYxf61",
	"295Avbam4qGShBYK7E9H5Ldew4B3WyJGmgmXfiWdeeyV",
	"EDi4rSy2LZgKJX74mbLTFk4mxoTgT6F7HxxzG2HBAFyK",
	"BnGKHAC386n4Qmv9xtpBVbRaUTKixjBe3oagkPFKtoy6",
	"Dd7K2Fp7AtoN8xCghKDRmyqr5U169t48Tw5fEd3wT9mq",
	"AP6qExwrbRgBAVaehg4b5xHENX815sMabtBzUzVB4v8S",
}

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

	// grpc request warmup
	client.GetHealth(context.Background(), &pb.HealthRequest{})

	// new rpc client and get latest block hash
	rpcClient := rpc.New(mainNetRPC)
	blockhash, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		panic(fmt.Sprintf("[get latest block hash] error: %v", err))
	}

	tipAccount := tipAccounts[rand.Intn(len(tipAccounts))]

	// construct instruction
	transferIx := system.NewTransferInstruction(amount, account.PublicKey(), receivePub).Build()
	tipIx := system.NewTransferInstruction(tipAmount, account.PublicKey(), solana.MustPublicKeyFromBase58(tipAccount)).Build()

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
