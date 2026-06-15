# Solana-trader-client-go

Example for solana-trader-client in Go.

# Document

See [document](https://blockrazor.gitbook.io/blockrazor/solana/send-transaction/go).

# Transaction Encoding

Two transaction encoding methods are supported:

| Encoding | gRPC Method | HTTP Endpoint | Request Field |
|---|---|---|---|
| Base64 | `SendTransaction` | `/sendTransaction` | `transaction` |
| Binary | `SendBinaryTransaction` | `/sendBinaryTransaction` | `binaryTransaction` |

Binary examples use `tx.MarshalBinary()`.

For HTTP Binary requests, the request body is still JSON. Go's `encoding/json`
package encodes the `[]byte` value as a Base64 JSON string automatically.

# Quickstart

1. **Download git repository**

   ```bash
   git clone https://github.com/BlockRazorinc/solana-trader-client-go.git
   ```

2. **Change directory**

   ```bash
   cd solana-trader-client-go
   ```

3. **Download dependencies**

   ```bash
   go mod tidy
   ```

4. **Edit `example/grpc/mode_fast/main.go`**

   ```go
   // BlockRazor relay endpoint address
   blzRelayEndpoint = "frankfurt.solana-grpc.blockrazor.xyz:80"
   // replace your solana rpc endpoint
   mainNetRPC = ""
   // replace your authKey
   authKey = ""
   // replace your private key (base58)
   privateKey = ""
   // publicKey (base58)
   publicKey = ""
   // transfer amount
   amount = 200_000
   // send mode
   mode = "fast"
   // tip amount
   tipAmount = 1_000_000
   ```

5. **Run the gRPC fast mode Base64 example**

   ```bash
   go run example/grpc/mode_fast/main.go
   ```

For Binary transaction examples, see the Binary sections below.

# gRPC

## Fast Mode

### Base64

Uses `SendTransaction` with a Base64-encoded transaction.

Edit:

```text
example/grpc/mode_fast/main.go
```

Run:

```bash
go run example/grpc/mode_fast/main.go
```

### Binary

Uses `SendBinaryTransaction` with transaction bytes.

Edit:

```text
example/grpc/mode_fast_binary/main.go
```

Run:

```bash
go run example/grpc/mode_fast_binary/main.go
```

## Sandwich Mitigation Mode

Set `safeWindow` and `revertProtection` according to your requirements.

### Base64

Uses `SendTransaction` with a Base64-encoded transaction.

Edit:

```text
example/grpc/mode_sandwichMitigation/main.go
```

Run:

```bash
go run example/grpc/mode_sandwichMitigation/main.go
```

### Binary

Uses `SendBinaryTransaction` with transaction bytes.

Edit:

```text
example/grpc/mode_sandwichMitigation_binary/main.go
```

Run:

```bash
go run example/grpc/mode_sandwichMitigation_binary/main.go
```

# HTTP

## Fast Mode

### Base64

Uses `/sendTransaction` with the `transaction` request field.

Edit:

```text
example/http/mode_fast/main.go
```

Run:

```bash
go run example/http/mode_fast/main.go
```

### Binary

Uses `/sendBinaryTransaction` with the `binaryTransaction` request field.

Edit:

```text
example/http/mode_fast_binary/main.go
```

Run:

```bash
go run example/http/mode_fast_binary/main.go
```

## Sandwich Mitigation Mode

Set `safeWindow` and `revertProtection` according to your requirements.

### Base64

Uses `/sendTransaction` with the `transaction` request field.

Edit:

```text
example/http/mode_sandwichMitigation/main.go
```

Run:

```bash
go run example/http/mode_sandwichMitigation/main.go
```

### Binary

Uses `/sendBinaryTransaction` with the `binaryTransaction` request field.

Edit:

```text
example/http/mode_sandwichMitigation_binary/main.go
```

Run:

```bash
go run example/http/mode_sandwichMitigation_binary/main.go
```
