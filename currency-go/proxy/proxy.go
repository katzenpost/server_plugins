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
	"github.com/btcsuite/btcd/btcjson"
	"github.com/katzenpost/server_plugins/currency-go/common"
	"github.com/katzenpost/server_plugins/currency-go/config"
)

type Currency struct {
	jsonHandle codec.JsonHandle

	ticker  string
	rpcUser string
	rpcPass string
	rpcUrl  string
}

func (Currency) OnRequest(payload string) (string, error) {

	return payload, nil // XXX fix me
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
