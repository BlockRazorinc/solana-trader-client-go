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

4. **Edit example/grpc/mode_fast/main.go**

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

5. **Run grpc/mode_fast example**
   
   `go run example/grpc/mode_fast/main.go`

# GRPC

## fast mode

1. **Edit example/grpc/mode_fast/main.go**
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
 
2. **Run grpc/mode_fast example**
   
   `go run example/grpc/mode_fast/main.go`

## sandwichMitigation mode


1. **Edit example/grpc/mode_sandwichMitigation/main.go**
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
	mode = "sandwichMitigation"
	// safeWindow
	safeWindow = 5
	// revertProtection
	revertProtection = false
	// tip amount
	tipAmount = 1_000_000
	```
 
2. **Run grpc/mode_sandwichMitigation example**
   
   `go run example/grpc/mode_sandwichMitigation/main.go`

# HTTP


## fast mode

1. **Edit example/http/mode_fast/main.go**
    ```
	httpEndpoint   = "http://frankfurt.solana.blockrazor.xyz:443/sendTransaction"
	healthEndpoint = "http://frankfurt.solana.blockrazor.xyz:443/health"
	mainNetRPC     = ""
	authKey        = ""
	privateKey     = ""
	publicKey      = ""
	amount         = 200_000
	tipAmount      = 1_000_000
	mode           = "fast"
	```
 
2. **Run http/mode_fast example**
   
   `go run example/http/mode_fast/main.go`

## sandwichMitigation mode


1. **Edit example/http/mode_sandwichMitigation/main.go**
    ```
	httpEndpoint     = "http://frankfurt.solana.blockrazor.xyz:443/sendTransaction"
	healthEndpoint   = "http://frankfurt.solana.blockrazor.xyz:443/health"
	mainNetRPC       = ""
	authKey          = ""
	privateKey       = ""
	publicKey        = ""
	amount           = 200_000
	tipAmount        = 1_000_000
	mode             = "sandwichMitigation"
	safeWindow       = 5
	revertProtection = false
	```
 
2. **Run http/mode_sandwichMitigation example**
   
   `go run example/http/mode_sandwichMitigation/main.go`
