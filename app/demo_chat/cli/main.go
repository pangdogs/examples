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
	"github.com/peterh/liner"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
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

	logger, _ := zap.NewDevelopment(zap.IncreaseLevel(zap.InfoLevel))
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

	go proc.Console()

	<-rpcli.Done()

	if err := context.Cause(rpcli); err != nil {
		rpcli.GetLogger().Infof("close cause:%s", err)
	}
}

type MainProc struct {
	rpcli.Procedure
}

func (p *MainProc) Console() {
	line := liner.NewLiner()
	defer line.Close()

	curChannel := misc.GlobalChannel

	for {
		text, err := line.Prompt(fmt.Sprintf("%s > ", curChannel))
		if err != nil {
			return
		}

		fields := strings.Fields(text)
		if len(fields) < 1 {
			continue
		}
		line.AppendHistory(text)

		switch strings.ToLower(fields[0]) {
		case "create":
			if len(fields) < 2 {
				continue
			}
			channel := fields[1]
			if err := rpc.ResultVoid(<-p.GetCli().RPC(misc.Gate, "ChatChannelComp", "C_CreateChannel", channel)).Extract(); err != nil {
				p.GetCli().GetLogger().Debugf("create channel %s failed, %s", channel, err)
				continue
			}
			p.GetCli().GetLogger().Debugf("create channel %s ok", channel)
		case "remove":
			if len(fields) < 2 {
				continue
			}
			channel := fields[1]
			if err := rpc.ResultVoid(<-p.GetCli().RPC(misc.Gate, "ChatChannelComp", "C_RemoveChannel", channel)).Extract(); err != nil {
				p.GetCli().GetLogger().Debugf("remove channel %s failed, %s", channel, err)
				continue
			}
			p.GetCli().GetLogger().Debugf("remove channel %s ok", channel)
		case "join":
			if len(fields) < 2 {
				continue
			}
			channel := fields[1]
			if err := rpc.ResultVoid(<-p.GetCli().RPC(misc.Gate, "ChatChannelComp", "C_JoinChannel", channel)).Extract(); err != nil {
				p.GetCli().GetLogger().Debugf("join channel %s failed, %s", channel, err)
				continue
			}
			p.GetCli().GetLogger().Debugf("join channel %s ok", channel)
		case "leave":
			if len(fields) < 2 {
				continue
			}
			channel := fields[1]
			if err := rpc.ResultVoid(<-p.GetCli().RPC(misc.Gate, "ChatChannelComp", "C_LeaveChannel", channel)).Extract(); err != nil {
				p.GetCli().GetLogger().Debugf("leave channel %s failed, %s", channel, err)
				continue
			}
			p.GetCli().GetLogger().Debugf("leave channel %s ok", channel)
		case "switch":
			if len(fields) < 2 {
				continue
			}
			channel := fields[1]
			if err := rpc.ResultVoid(<-p.GetCli().RPC(misc.Chat, "ChatUserComp", "C_SwitchChannel", channel)).Extract(); err != nil {
				p.GetCli().GetLogger().Debugf("switch channel %s failed, %s", channel, err)
				continue
			}
			p.GetCli().GetLogger().Debugf("switch channel %s ok", channel)
			curChannel = channel
		default:
			if err := rpc.ResultVoid(<-p.GetCli().RPC(misc.Chat, "ChatUserComp", "C_InputText", text)).Extract(); err != nil {
				p.GetCli().GetLogger().Debugf("input %s failed, %s", text, err)
				continue
			}
			p.GetCli().GetLogger().Debugf("input %s ok", text)
		}
	}
}

func (p *MainProc) OutputText(ts int64, channel, userId, text string) {
	fmt.Printf("[%s][%s] %s: %s\n", time.Unix(ts, 0).Format(time.RFC3339), channel, userId, text)
}
