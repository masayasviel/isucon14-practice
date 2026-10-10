package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

var erroredUpstream = errors.New("errored upstream")

// 決済APIは評価のたびに呼ばれ、失敗時は1リクエストで最大6往復する。
// http.DefaultClient のトランスポートは MaxIdleConnsPerHost が2しかないため、
// 同時リクエストが増えると接続が再利用されず毎回TCPハンドシェイクが走る。
// Timeout は ctx のキャンセルが効かない場合の保険で、通常はリクエストの
// ctx（クライアント切断で失効する）のほうが先に効く。
var paymentGatewayClient = newPaymentGatewayClient()

func newPaymentGatewayClient() *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 512
	transport.MaxIdleConnsPerHost = 256
	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}
}

// 本文を読み切ってから閉じないと、その接続はアイドルプールへ戻らず破棄される
func closeAndDrain(body io.ReadCloser) {
	io.Copy(io.Discard, body)
	body.Close()
}

type paymentGatewayPostPaymentRequest struct {
	Amount int `json:"amount"`
}

type paymentGatewayGetPaymentsResponseOne struct {
	Amount int    `json:"amount"`
	Status string `json:"status"`
}

func requestPaymentGatewayPostPayment(ctx context.Context, paymentGatewayURL string, token string, param *paymentGatewayPostPaymentRequest, retrieveRidesOrderByCreatedAtAsc func() ([]Ride, error)) error {
	b, err := json.Marshal(param)
	if err != nil {
		return err
	}

	// 失敗したらとりあえずリトライ
	// FIXME: 社内決済マイクロサービスのインフラに異常が発生していて、同時にたくさんリクエストすると変なことになる可能性あり
	retry := 0
	for {
		err := func() error {
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, paymentGatewayURL+"/payments", bytes.NewBuffer(b))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+token)

			res, err := paymentGatewayClient.Do(req)
			if err != nil {
				return err
			}
			defer closeAndDrain(res.Body)

			if res.StatusCode != http.StatusNoContent {
				// エラーが返ってきても成功している場合があるので、社内決済マイクロサービスに問い合わせ
				getReq, err := http.NewRequestWithContext(ctx, http.MethodGet, paymentGatewayURL+"/payments", bytes.NewBuffer([]byte{}))
				if err != nil {
					return err
				}
				getReq.Header.Set("Authorization", "Bearer "+token)

				getRes, err := paymentGatewayClient.Do(getReq)
				if err != nil {
					return err
				}
				defer closeAndDrain(getRes.Body)

				// GET /payments は障害と関係なく200が返るので、200以外は回復不能なエラーとする
				if getRes.StatusCode != http.StatusOK {
					return fmt.Errorf("[GET /payments] unexpected status code (%d)", getRes.StatusCode)
				}
				var payments []paymentGatewayGetPaymentsResponseOne
				if err := json.NewDecoder(getRes.Body).Decode(&payments); err != nil {
					return err
				}

				rides, err := retrieveRidesOrderByCreatedAtAsc()
				if err != nil {
					return err
				}

				if len(rides) != len(payments) {
					return fmt.Errorf("unexpected number of payments: %d != %d. %w", len(rides), len(payments), erroredUpstream)
				}

				return nil
			}
			return nil
		}()
		if err != nil {
			if retry < 5 {
				retry++
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				return err
			}
		}
		break
	}

	return nil
}
