# Solana-trader-client-go
example for solana-trader-client in Go

# Document
see [document](https://blockrazor.gitbook.io/blockrazor/solana/send-transaction/go)

# Quickstart

1. **Download git repository**
   
   `git clone https://github.com/BlockRazorinc/solana-trader-client-go.git`

2. **Change directory**
   
   `cd solana-trader-client-go`

3. **Download dependencies**
   
   `go mod tidy`

4. **Edit example/main.go**

	```
	// BlockRazor relay endpoint address
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
	mode = "fast"

	// tip amount
	tipAmount = 1_000_000
	```

5.1 **Run mode_fast example**
   
   `cd example/mode_fast && go run main.go`

5.2 **Run mode_sandwichMitigation example**
   
   `cd example/mode_sandwichMitigation && go run main.go`
