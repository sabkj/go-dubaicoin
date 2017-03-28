// Copyright 2015 The go-dubaicoin Authors
// This file is part of go-dubaicoin.
//
// go-dubaicoin is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-dubaicoin is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-dubaicoin. If not, see <http://www.gnu.org/licenses/>.

// bootnode runs a bootstrap node for the Dubaicoin Discovery Protocol.
package main

import (
	"crypto/ecdsa"
	"flag"
	"os"

	"github.com/dbix-project/go-dubaicoin/cmd/utils"
	"github.com/dbix-project/go-dubaicoin/crypto"
	"github.com/dbix-project/go-dubaicoin/logger/glog"
	"github.com/dbix-project/go-dubaicoin/p2p/discover"
	"github.com/dbix-project/go-dubaicoin/p2p/nat"

)

func main() {
	var (
		listenAddr  = flag.String("addr", ":57956", "listen address")
		genKey      = flag.String("genkey", "", "generate a node key and quit")
		nodeKeyFile = flag.String("nodekey", "", "private key filename")
		nodeKeyHex  = flag.String("nodekeyhex", "", "private key as hex (for testing)")
		natdesc     = flag.String("nat", "none", "port mapping mechanism (any|none|upnp|pmp|extip:<IP>)")

		nodeKey *ecdsa.PrivateKey
		err     error
	)
	flag.Var(glog.GetVerbosity(), "verbosity", "log verbosity (0-9)")
	flag.Var(glog.GetVModule(), "vmodule", "log verbosity pattern")
	glog.SetToStderr(true)
	flag.Parse()

	if *genKey != "" {
		key, err := crypto.GenerateKey()
		if err != nil {
			utils.Fatalf("could not generate key: %v", err)
		}
		if err := crypto.SaveECDSA(*genKey, key); err != nil {
			utils.Fatalf("%v", err)
		}
		os.Exit(0)
	}

	natm, err := nat.Parse(*natdesc)
	if err != nil {
		utils.Fatalf("-nat: %v", err)
	}
	switch {
	case *nodeKeyFile == "" && *nodeKeyHex == "":
		utils.Fatalf("Use -nodekey or -nodekeyhex to specify a private key")
	case *nodeKeyFile != "" && *nodeKeyHex != "":
		utils.Fatalf("Options -nodekey and -nodekeyhex are mutually exclusive")
	case *nodeKeyFile != "":
		if nodeKey, err = crypto.LoadECDSA(*nodeKeyFile); err != nil {
			utils.Fatalf("-nodekey: %v", err)
		}
	case *nodeKeyHex != "":
		if nodeKey, err = crypto.HexToECDSA(*nodeKeyHex); err != nil {
			utils.Fatalf("-nodekeyhex: %v", err)
		}
	}

	if _, err := discover.ListenUDP(nodeKey, *listenAddr, natm, ""); err != nil {
		utils.Fatalf("%v", err)
	}
	select {}
}
