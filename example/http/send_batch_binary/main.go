package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/rpc"
)

const (
	httpBinaryEndpoint = "http://frankfurt.solana.blockrazor.xyz:443/v2/sendBinaryBatch"
	healthEndpoint     = "http://frankfurt.solana.blockrazor.xyz:443/health"
	mainNetRPC         = ""
	authKey            = ""
	privateKey         = ""
	publicKey          = ""
	amount             = 200_000
	tipAmount          = 1_000_000
	mode               = "fast"
	safeWindow         = 5
	revertProtection   = false
)

var tipAccounts = []string{
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
	"8EpzWhsr7uyZkfVD7MTHoWqupM14xngzMeUcCkWbY1UG",
	"56XRik27PddNnR7JkHKqd5pWdMHeqje7TCm1gsdmvnTP",
	"2DJNtR1tT7vfub7sDJ6v9NGmrzGVFUer5FvjhEcWEpvv",
	"3AveZTS3PJo5qKEqHJci5bjW8neufriirr9mJXfsv4iH",
	"9DC798ezxU2f2i8ggQ4vsqrsNxdDWJWAFg4cT1xLDR1X",
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// BatchResponse .
type BatchResponse struct {
	Result []BatchResponseItem `json:"result"`
	Error  string              `json:"error"`
}

// BatchResponseItem .
type BatchResponseItem struct {
	Signature string `json:"signature"`
	Error     string `json:"error"`
}

type HealthResponse struct {
	Result string `json:"result"`
}

func main() {
	// Pre-warm: perform an initial health check to establish the HTTP connection
	err := pingHealth()
	if err != nil {
		fmt.Printf("health check failed: %v\n", err)
	}
	// Start a background goroutine to periodically send /health requests
	// For low-frequency users, this keeps the HTTP connection alive (warm)
	go func() {
		for {
			err := pingHealth()
			if err != nil {
				fmt.Printf("health check failed: %v\n", err)
			}
			time.Sleep(30 * time.Second)
		}
	}()
	// send binary batch
	if err := sendBinaryBatch(); err != nil {
		fmt.Printf("send binary batch failed: %v\n", err)
	}
}

func pingHealth() error {
	req, err := http.NewRequest("GET", healthEndpoint, nil)
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
	// ⚠️ Important Note:
	// According to the Go net/http documentation, in order for the underlying TCP connection
	// to be reused (i.e. kept alive), the response body must be fully read and closed.
	// Otherwise, the Transport may not reuse the connection for future requests.
	// Reference: https://pkg.go.dev/net/http#Response
	// > "The default HTTP client's Transport may not reuse HTTP/1.x 'keep-alive' TCP connections
	//    if the Body is not read to completion and closed."

	// Read the full response body to enable connection reuse
	bodyBytes, _ := io.ReadAll(resp.Body)
	var healthRes HealthResponse
	if err := json.Unmarshal(bodyBytes, &healthRes); err != nil {
		return fmt.Errorf("decode error: %v", err)
	}

	return nil
}

func sendBinaryBatch() error {
	account, err := solana.WalletFromPrivateKeyBase58(privateKey)
	if err != nil {
		return err
	}
	receivePub := solana.MustPublicKeyFromBase58(publicKey)
	tipPub := solana.MustPublicKeyFromBase58(tipAccounts[rand.Intn(len(tipAccounts))])

	rpcClient := rpc.New(mainNetRPC)
	blockhash, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return fmt.Errorf("[get blockhash] %v", err)
	}

	transferIx := system.NewTransferInstruction(amount, account.PublicKey(), receivePub).Build()
	tipIx := system.NewTransferInstruction(tipAmount, account.PublicKey(), tipPub).Build()

	tx, err := solana.NewTransaction(
		[]solana.Instruction{tipIx, transferIx},
		blockhash.Value.Blockhash,
		solana.TransactionPayer(account.PublicKey()),
	)
	if err != nil {
		return fmt.Errorf("build tx error: %v", err)
	}

	_, err = tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if account.PublicKey().Equals(key) {
			return &account.PrivateKey
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("sign tx error: %v", err)
	}

	txBin, err := tx.MarshalBinary()
	if err != nil {
		return err
	}

	txBinBody, err := buildBinaryBody([][]byte{txBin})
	if err != nil {
		return err
	}

	u, err := url.Parse(httpBinaryEndpoint)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("mode", mode)
	q.Set("safeWindow", strconv.Itoa(safeWindow))
	q.Set("revertProtection", strconv.FormatBool(revertProtection))
	q.Set("auth", authKey)
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequest("POST", u.String(), bytes.NewBuffer(txBinBody))
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/octet-stream")
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("send http error: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	var sendRes BatchResponse
	if err := json.Unmarshal(bodyBytes, &sendRes); err != nil {
		return fmt.Errorf("decode error: %v", err)
	}
	fmt.Printf("[send batch] response: %+v\n", sendRes)
	return nil
}

func buildBinaryBody(binaryTransactions [][]byte) ([]byte, error) {
	var body bytes.Buffer

	for i, txBytes := range binaryTransactions {
		if len(txBytes) > 65535 {
			return nil, fmt.Errorf("tx[%d] too large: %d bytes", i, len(txBytes))
		}

		txLen := uint16(len(txBytes))

		if err := binary.Write(&body, binary.BigEndian, txLen); err != nil {
			return nil, fmt.Errorf("write tx[%d] length failed: %w", i, err)
		}

		if _, err := body.Write(txBytes); err != nil {
			return nil, fmt.Errorf("write tx[%d] bytes failed: %w", i, err)
		}
	}

	return body.Bytes(), nil
}
