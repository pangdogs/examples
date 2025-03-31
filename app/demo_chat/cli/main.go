/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package main

import (
	"context"
	"fmt"
	"git.golaxy.org/examples/app/demo_chat/misc"
	"git.golaxy.org/framework/addins/gate/cli"
	"git.golaxy.org/framework/addins/rpc"
	"git.golaxy.org/framework/addins/rpc/rpcli"
	"git.golaxy.org/framework/net/gtp"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"time"
)

func main() {
	pflag.String("cli_priv_key", "cli.pem", "client private key for sign")
	pflag.String("serv_pub_key", "serv.pub", "service public key for verify sign")
	pflag.String("endpoint", "localhost:9090", "connect endpoint")
	pflag.Bool("ws", false, "use websocket")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	cliPrivKey, err := gtp.LoadPrivateKeyFile(viper.GetString("cli_priv_key"))
	if err != nil {
		panic(err)
	}

	servPubKey, err := gtp.LoadPublicKeyFile(viper.GetString("serv_pub_key"))
	if err != nil {
		panic(err)
	}

	np := cli.TCP
	if viper.GetBool("ws") {
		np = cli.WebSocket
	}

	logger, _ := zap.NewDevelopment()
	proc := &MainProc{}

	rpcli, err := rpcli.BuildRPCli().
		SetNetProtocol(np).
		SetIOTimeout(10*time.Second).
		SetGTPAutoReconnect(true).
		SetGTPEncCipherSuite(gtp.CipherSuite{
			SecretKeyExchange:   gtp.SecretKeyExchange_ECDHE,
			SymmetricEncryption: gtp.SymmetricEncryption_AES,
			BlockCipherMode:     gtp.BlockCipherMode_GCM,
			PaddingMode:         gtp.PaddingMode_Pkcs7,
			MACHash:             gtp.Hash_Fnv1a64,
		}).
		SetGTPEncSignatureAlgorithm(gtp.SignatureAlgorithm{
			AsymmetricEncryption: gtp.AsymmetricEncryption_RSA256,
			PaddingMode:          gtp.PaddingMode_Pkcs1v15,
			Hash:                 gtp.Hash_SHA256,
		}).
		SetGTPEncSignaturePrivateKey(cliPrivKey).
		SetGTPEncVerifyServerSignature(true).
		SetGTPEncVerifySignaturePublicKey(servPubKey).
		SetGTPCompression(gtp.Compression_Brotli).
		SetGTPCompressedSize(0).
		SetGTPAutoReconnectRetryTimes(0).
		SetZapLogger(logger).
		SetMainProcedure(proc).
		Connect(context.Background(), viper.GetString("endpoint"))
	if err != nil {
		panic(err)
	}

	go proc.MainLoop()

	<-rpcli.Done()

	if err := context.Cause(rpcli); err != nil {
		rpcli.GetLogger().Infof("close cause:%s", err)
	}
}

type MainProc struct {
	rpcli.Procedure
}

func (p *MainProc) MainLoop() {
	go func() {
		for {
			var txt string
			fmt.Scanln(&txt)

			sendTs := time.Now()
			ret := rpc.ResultVoid(<-p.GetCli().RPC(misc.Chat, "ChatUserComp", "C_Input", txt)).Extract()
			if ret != nil {
				p.GetCli().GetLogger().Infof("input failed, delay:%dms, %s", time.Now().Sub(sendTs).Milliseconds(), ret.Error)
				continue
			}

			p.GetCli().GetLogger().Infof("input ok, delay:%dms", time.Now().Sub(sendTs).Milliseconds())
		}
	}()

}
