// main.go - Crypto currency transaction submition Kaetzchen service plugin program.
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

package main

import (
	"bytes"
	"errors"
	"flag"
	"net/http"

	"github.com/hashicorp/go-plugin"
	"github.com/katzenpost/server_plugins/currency-go/common"
	"github.com/katzenpost/server_plugins/currency-go/config"
	"github.com/ugorji/go/codec"
	"gopkg.in/op/go-logging.v1"
)

func main() {
	cfgFile := flag.String("f", "currency.toml", "Path to the currency config file.")
	flag.Parse()

	// Set the umask to something "paranoid".
	syscall.Umask(0077)

	cfg, err := LoadFile(*cfgFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config file '%v': %v\n", *cfgFile, err)
		os.Exit(-1)
	}

	currency := proxy.New(cfg)
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: common.Handshake,
		Plugins: map[string]plugin.Plugin{
			currencyCapability: &common.KaetzchenPlugin{Impl: currency},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
