package main

import (
	"context"
	"fmt"
	"git.golaxy.org/examples/app/demo_cs/demo_server/comp"
	"git.golaxy.org/examples/app/demo_cs/demo_server/misc"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework/net/gtp"
	"git.golaxy.org/framework/plugins/gate/cli"
	"git.golaxy.org/framework/plugins/rpc"
	"git.golaxy.org/framework/plugins/rpc/rpcli"
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

	cliPrivKey, err := misc.LoadPrivateKey(viper.GetString("cli_priv_key"))
	if err != nil {
		panic(err)
	}
	_ = cliPrivKey

	servPubKey, err := misc.LoadPublicKey(viper.GetString("serv_pub_key"))
	if err != nil {
		panic(err)
	}
	_ = servPubKey

	logger, _ := zap.NewDevelopment()

	proc := &MainProc{}

	np := cli.TCP
	if viper.GetBool("ws") {
		np = cli.WebSocket
	}

	rpcli, err := rpcli.CreateRPCli().
		NetProtocol(np).
		IOTimeout(7*time.Second).
		GTPAutoReconnect(true).
		GTPEncCipherSuite(gtp.CipherSuite{
			SecretKeyExchange:   gtp.SecretKeyExchange_ECDHE,
			SymmetricEncryption: gtp.SymmetricEncryption_AES,
			BlockCipherMode:     gtp.BlockCipherMode_GCM,
			PaddingMode:         gtp.PaddingMode_Pkcs7,
			MACHash:             gtp.Hash_Fnv1a64,
		}).
		GTPEncSignatureAlgorithm(gtp.SignatureAlgorithm{
			AsymmetricEncryption: gtp.AsymmetricEncryption_RSA_256,
			PaddingMode:          gtp.PaddingMode_Pkcs1v15,
			Hash:                 gtp.Hash_SHA256,
		}).
		GTPEncSignaturePrivateKey(cliPrivKey).
		GTPEncVerifyServerSignature(true).
		GTPEncVerifySignaturePublicKey(servPubKey).
		GTPCompression(gtp.Compression_Brotli).
		GTPCompressedSize(0).
		GTPAutoReconnectRetryTimes(0).
		FutureTimeout(10*time.Second).
		ZapLogger(logger).
		MainProcedure(proc).
		Connect(context.Background(), viper.GetString("endpoint"))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			var txt string
			fmt.Scanln(&txt)

			send := time.Now()
			ret, err := rpc.Result1[string](<-proc.RPC(misc.Work, comp.CmdCompSelf.Name, "Echo", txt))
			if err != nil {
				fmt.Printf("<= delay:%dms, error: %s\n", time.Now().Sub(send).Milliseconds(), err)
				continue
			}

			fmt.Printf("<= delay:%dms, echo: %s\n", time.Now().Sub(send).Milliseconds(), ret)
		}
	}()

	<-rpcli.Done()

	if err := context.Cause(rpcli); err != nil {
		fmt.Println("cause:", err)
	}
}

type MainProc struct {
	rpcli.Procedure
}
