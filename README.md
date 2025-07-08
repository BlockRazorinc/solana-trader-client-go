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

4. **Edit example/mode_fast/main.go or example/mode_sandwichMitigation/main.go**

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
	mode = ""

	// tip amount
	tipAmount = 1_000_000
	```

5. **Run mode_fast example**
   
   `go run example/mode_fast/main.go`

6. **Run mode_sandwichMitigation example**
   
   `go run example/mode_sandwichMitigation/main.go`
