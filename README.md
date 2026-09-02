# BlockRazor Solana Transaction Sending — Go Examples

Go examples for sending signed Solana transactions through [BlockRazor](https://blockrazor.io/) over HTTP or gRPC. The repository demonstrates low-latency single-transaction submission, batch submission, and atomic bundles in both Fast and Sandwich Mitigation modes.

> These examples construct and sign transactions locally. Private keys are used only by the local Go process and must never be sent to BlockRazor or committed to source control.

## What this repository covers

| Capability | HTTP | gRPC | Examples in this repository |
| --- | --- | --- | --- |
| Send Transaction | `/sendTransaction`, `/sendBinaryTransaction` | `SendTransaction`, `SendBinaryTransaction` | Base64 and binary; Fast and Sandwich Mitigation |
| Send Batch | `/sendBatch`, `/v2/sendBinaryBatch` | `SendBatch` | Base64 over HTTP/gRPC and raw binary over HTTP |
| Send Bundle | `/sendBundle` | `SendBundle` | Base64 over HTTP/gRPC |
| Keep Alive | `/health` | `GetHealth` | Initial health check and a 30-second health interval |

BlockRazor also provides v2 Plain Text and raw Binary interfaces. This repository currently includes the raw HTTP binary batch example, but it does not include v2 Plain Text or raw Binary single-transaction examples. See the [official Solana Transaction Sending documentation](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana) for every supported interface.

## Fast and Sandwich Mitigation modes

BlockRazor supports two transaction-sending modes:

| Mode | Use it when | Important behavior |
| --- | --- | --- |
| `fast` | Inclusion speed is the main priority | Uses BlockRazor's globally distributed network and SWQoS routing to reach the upcoming Leader with low latency; this is the default mode |
| `sandwichMitigation` | You need additional protection against sandwich attacks | Routes through trusted SWQoS paths and avoids dynamically identified blacklisted Leaders; inclusion may take longer |

Sandwich Mitigation accepts two additional parameters:

- `safeWindow`: number of consecutive allowlisted-validator slots required before sending. The supported range is `3`–`13`; the default is `3`.
- `revertProtection`: when `true`, BlockRazor attempts to prevent the transaction from failing on-chain, which can reduce inclusion speed or prevent inclusion.

BlockRazor's official mode documentation states that durable nonce transactions must not be used with `sandwichMitigation`, because doing so makes the sandwich protection ineffective. This is a service constraint rather than a check implemented by these Go examples. See [Solana Send Transaction](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/send-transaction#request-parameter).

## Requirements

- Go 1.23.1 or later
- A BlockRazor auth key
- A Solana RPC endpoint for fetching a recent blockhash
- A funded Solana keypair for examples that transfer SOL and pay fees/tips

The generated Protobuf and gRPC Go files are committed under `pb/serverpb/`, so `protoc` is not required to build or run the existing examples.

Obtain an auth key by following [How to Obtain an Auth Token](https://docs.blockrazor.io/get-started/authentication). The default rate limit for transaction and batch sending is documented as 3 TPS; bundle access is not enabled by default. Contact BlockRazor if you need bundle access or a higher limit.

## Quick start

```bash
git clone https://github.com/BlockRazorinc/solana-trader-client-go.git
cd solana-trader-client-go
go mod download
```

Open the example you want to run and replace its configuration placeholders:

```go
const (
	mainNetRPC = "YOUR_SOLANA_RPC_ENDPOINT"
	authKey    = "YOUR_BLOCKRAZOR_AUTH_KEY"
	privateKey = "YOUR_BASE58_PRIVATE_KEY"
	publicKey  = "RECIPIENT_PUBLIC_KEY"
)
```

Each example contains a default Frankfurt endpoint. For production, select the region closest to your application from the current [Solana Transaction Sending endpoint list](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/endpoint).

> Running an example submits a real mainnet transaction. The examples include a tip transfer and may also include a recipient transfer. Review `tipAmount`, `amount`, the recipient, the selected tip account, and the transaction instructions before running them.

## Example matrix

### Send one transaction

| Transport | Mode | Payload | Command |
| --- | --- | --- | --- |
| gRPC | Fast | Base64 string | `go run ./example/grpc/mode_fast` |
| gRPC | Fast | Protobuf bytes | `go run ./example/grpc/mode_fast_binary` |
| gRPC | Sandwich Mitigation | Base64 string | `go run ./example/grpc/mode_sandwichMitigation` |
| gRPC | Sandwich Mitigation | Protobuf bytes | `go run ./example/grpc/mode_sandwichMitigation_binary` |
| HTTP | Fast | Base64 in JSON | `go run ./example/http/mode_fast` |
| HTTP | Fast | Base64 in the JSON `binaryTransaction` field | `go run ./example/http/mode_fast_binary` |
| HTTP | Sandwich Mitigation | Base64 in JSON | `go run ./example/http/mode_sandwichMitigation` |
| HTTP | Sandwich Mitigation | Base64 in the JSON `binaryTransaction` field | `go run ./example/http/mode_sandwichMitigation_binary` |

Transactions are constructed and signed locally with `gagliardetto/solana-go`. Base64 examples use `Transaction.ToBase64()`. Binary examples use `Transaction.MarshalBinary()`.

For HTTP `/sendBinaryTransaction`, the request body is still JSON. Go's `encoding/json` package serializes the `[]byte` field as a Base64 JSON string automatically. For gRPC `SendBinaryTransaction`, the serialized `[]byte` is sent directly as Protobuf bytes.

### Send a batch

Batch submission sends multiple signed transactions in one request. A batch can contain up to 25 fully signed transactions, and the response contains a signature or error for each item.

| Transport | Payload | Command |
| --- | --- | --- |
| gRPC | Base64 string array | `go run ./example/grpc/send_batch` |
| HTTP | Base64 string array in JSON | `go run ./example/http/send_batch` |
| HTTP v2 | Length-prefixed raw binary transactions | `go run ./example/http/send_batch_binary` |

The v2 binary batch body is a sequence of serialized transactions, each prefixed by a two-byte, big-endian length. The Go example builds this payload with `encoding/binary` and sends it as `application/octet-stream`. See [Solana Send Batch](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/send-batch) and [Send Batch in Binary](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/send-batch/send-in-binary) for the current contract.

### Send a bundle

A bundle contains up to four fully signed, Base64-encoded transactions. Transactions execute sequentially and atomically: if any transaction fails, the entire bundle fails.

| Transport | Payload | Command |
| --- | --- | --- |
| gRPC | Base64 string array | `go run ./example/grpc/send_bundle` |
| HTTP | Base64 string array in JSON | `go run ./example/http/send_bundle` |

Bundle access is not enabled by default. See [Solana Send Bundle](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/send-bundle) for activation and MEV-protection guidance.

## Authentication

Authentication depends on the interface:

| Interface | Auth location |
| --- | --- |
| HTTP JSON | `apikey` request header |
| gRPC | `apikey` per-RPC credentials metadata |
| HTTP v2 raw binary | `auth` URL query parameter |

The gRPC examples implement `credentials.PerRPCCredentials` through the local `Authentication` type. Never log auth keys or include them in committed files. Prefer environment variables or a secret manager when adapting these examples for production.

## Transaction encoding

| Interface | Field/body | Encoding |
| --- | --- | --- |
| gRPC `SendTransaction` | `transaction` | Base64 string |
| gRPC `SendBinaryTransaction` | `binaryTransaction` | Raw bytes in Protobuf |
| HTTP `/sendTransaction` | `transaction` | Base64 string in JSON |
| HTTP `/sendBinaryTransaction` | `binaryTransaction` | `[]byte`, encoded as Base64 by `encoding/json` |
| HTTP `/v2/sendTransaction` | Request body | Base64 text with `Content-Type: text/plain` |
| HTTP `/v2/sendBinaryTransaction` | Request body | Raw serialized bytes with `Content-Type: application/octet-stream` |

The last two v2 single-transaction interfaces are part of the BlockRazor service but do not currently have examples in this repository. Refer to [Send in Plain Text](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/send-transaction/send-in-plain-text) and [Send in Binary](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/send-transaction/send-in-binary).

## Priority fee and tip

The examples select a BlockRazor tip account at random and place a tip transfer instruction before the primary transfer instruction. Rotating tip accounts helps avoid performance degradation caused by account contention.

The current examples demonstrate the tip transfer but do not add a Compute Budget priority-fee instruction. For production transactions, configure both according to your strategy and current network conditions. Consult [Solana Transaction Sending Priority Fee & Tip](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/priority-fee-and-tip) for the current minimums, recommendations, and official tip accounts.

## Keep the connection warm

Every example performs a health request before sending and repeats it every 30 seconds in a background goroutine.

HTTP examples reuse a package-level `http.Client`. They read each response body to completion and close it, allowing Go's `http.Transport` to reuse the underlying connection. gRPC examples reuse a single `grpc.ClientConn` and generated `ServerClient`.

See [Solana Transaction Sending Keep Alive](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/keep-alive) for the current health endpoint contract.

## Go implementation

- `gagliardetto/solana-go` constructs, signs, serializes, and Base64-encodes transactions.
- `net/http` and `encoding/json` implement HTTP submission and response parsing.
- `google.golang.org/grpc` and the generated `pb/serverpb` package implement gRPC submission.
- `context` is passed through RPC, blockhash, and submission calls.
- `encoding/binary` builds the v2 length-prefixed binary batch payload.

The gRPC contract lives in [`server.proto`](server.proto):

```proto
service Server {
  rpc SendTransaction(SendRequest) returns (SendResponse);
  rpc SendBinaryTransaction(SendBinaryRequest) returns (SendResponse);
  rpc SendBatch(SendBatchRequest) returns (SendBatchResponse);
  rpc SendBundle(SendBundleRequest) returns (SendBundleResponse);
  rpc GetHealth(HealthRequest) returns (HealthResponse);
}
```

## Production checklist

- Load the auth key and private key from environment variables or a secret manager.
- Use a dedicated, low-balance wallet while adapting the examples.
- Validate the destination address, transaction instructions, amounts, priority fee, tip, and recent blockhash before signing.
- Select the nearest BlockRazor endpoint and test latency from your own deployment region.
- Reuse HTTP/gRPC connections and keep them warm before the critical send path.
- Add context deadlines, cancellation, retry behavior appropriate to transaction semantics, structured logging, metrics, and graceful shutdown.
- Check HTTP status codes and gRPC status codes before parsing successful responses.
- Do not blindly retry a request without checking its returned signature and on-chain status.
- Follow BlockRazor's documented Durable Nonce restriction when using Sandwich Mitigation mode.

## Documentation and support

- [Solana Transaction Sending](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana)
- [Endpoints](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/endpoint)
- [Send Transaction](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/send-transaction)
- [Send Batch](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/send-batch)
- [Send Bundle](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/send-bundle)
- [Priority Fee & Tip](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/priority-fee-and-tip)
- [Keep Alive](https://docs.blockrazor.io/transaction-submission/transaction-sending/solana/keep-alive)
- [BlockRazor Discord](https://discord.com/invite/qqJuwRb8Nh)
