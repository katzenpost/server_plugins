// proxy.go - Crypto currency transaction proxy.
// Copyright (C) 2018  David Stainton.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package proxy

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/btcsuite/btcd/btcjson"
	"github.com/katzenpost/server_plugins/currency-go/common"
	"github.com/katzenpost/server_plugins/currency-go/config"
	"github.com/ugorji/go/codec"
	"gopkg.in/op/go-logging.v1"
)

var errInvalidCurrencyRequest = errors.New("kaetzchen/currency: invalid request")

type Currency struct {
	log        *logging.Logger
	jsonHandle codec.JsonHandle

	ticker  string
	rpcUser string
	rpcPass string
	rpcUrl  string
}

func (k *Currency) OnRequest(payload []byte, hasSURB bool) ([]byte, error) {
	if hasSURB {
		k.log.Debug("Ignoring request: erroneously contains a SURB.")
		return nil, errInvalidCurrencyRequest
	}

	k.log.Debug("Handling request")

	// Parse out the request payload.
	var req common.CurrencyRequest
	dec := codec.NewDecoderBytes(bytes.TrimRight(payload, "\x00"), &k.jsonHandle)
	if err := dec.Decode(&req); err != nil {
		k.log.Debugf("Failed to decode request: (%v)", err)
		return nil, errInvalidCurrencyRequest
	}

	// Sanity check the request.
	if req.Version != common.CurrencyVersion {
		k.log.Debugf("Failed to parse request: (invalid version: %v)", req.Version)
		return nil, errInvalidCurrencyRequest
	}
	if req.Ticker != k.ticker {
		k.log.Debugf("Failed to parse request: (currency ticker mismatch: %v)", req.Ticker)
		return nil, errInvalidCurrencyRequest
	}

	// Send request to HTTP RPC.
	err := k.sendTransaction(req.Tx)
	if err != nil {
		k.log.Debug("Failed to send currency transaction request: (%v)", err)
	}
	return nil, errInvalidCurrencyRequest
}

func (k *Currency) sendTransaction(txHex string) error {
	k.log.Debug("sendTransaction")

	// marshall new transaction blob
	allowHighFees := true
	cmd := btcjson.NewSendRawTransactionCmd(txHex, &allowHighFees)
	txId := 0 // this txId is not important
	marshalledJSON, err := btcjson.MarshalCmd(txId, cmd)
	bodyReader := bytes.NewReader(marshalledJSON)

	// create an http request
	httpReq, err := http.NewRequest("POST", k.rpcUrl, bodyReader)
	if err != nil {
		return err
	}
	httpReq.Close = true
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.SetBasicAuth(k.rpcUser, k.rpcPass)

	// send http request
	client := http.Client{}
	httpResponse, err := client.Do(httpReq)
	k.log.Debugf("currency RPC response status: %s", httpResponse.Status)

	return err
}

func New(config *config.Config) *Currency {
	currency := &Currency{
		ticker:  config.Ticker,
		rpcUser: config.RPCUser,
		rpcPass: config.RPCPass,
		rpcUrl:  config.RPCURL,
	}
	currency.jsonHandle.Canonical = true
	currency.jsonHandle.ErrorIfNoField = true
	return currency
}
