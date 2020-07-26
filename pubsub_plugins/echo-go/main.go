// main.go - echo service using cbor plugin system
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
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/katzenpost/server/pubsubplugin/common"
	"github.com/katzenpost/server/pubsubplugin/server"
)

type EchoSpool struct {
	appMessageCh   chan *common.AppMessages
	subscriptionID *common.SubscriptionID
}

func New() *EchoSpool {
	return &EchoSpool{
		appMessageCh:   make(chan *common.AppMessages),
		subscriptionID: nil,
	}
}
func (e *EchoSpool) Subscribe(subscriptionID *common.SubscriptionID, spoolID *common.SpoolID, lastSpoolIndex uint64) error {
	e.subscriptionID = subscriptionID
	e.appMessageCh <- &common.AppMessages{
		SubscriptionID: *e.subscriptionID,
		Messages: []common.SpoolMessage{
			common.SpoolMessage{
				Index:   0,
				Payload: []byte("echo"),
			},
		},
	}
	return nil
}

func (e *EchoSpool) Unsubscribe(subscriptionID *common.SubscriptionID) error {
	e.subscriptionID = nil
	return nil
}

func main() {
	var logLevel string
	var logDir string
	flag.StringVar(&logDir, "log_dir", "", "logging directory")
	flag.StringVar(&logLevel, "log_level", "DEBUG", "logging level could be set to: DEBUG, INFO, NOTICE, WARNING, ERROR, CRITICAL")
	flag.Parse()

	e := New()

	config := &server.Config{
		Name:          "echo",
		Parameters:    new(common.Parameters),
		LogDir:        logDir,
		LogLevel:      logLevel,
		Spool:         e,
		AppMessagesCh: e.appMessageCh,
	}

	_, err := server.New(config)
	if err != nil {
		panic(err)
	}

	// Wait for a control-c and then exit.
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
