package main

import (
	"context"
	"fmt"
	"git.golaxy.org/examples/app/demo_cs/demo_server/comp"
	"git.golaxy.org/examples/app/demo_cs/demo_server/serv"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework/net/gtp"
	"git.golaxy.org/framework/plugins/rpc"
	"git.golaxy.org/framework/plugins/rpc/rpcli"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"time"
)

func main() {
	pflag.String("cli_priv_key", "cli.pem", "")
	pflag.String("serv_pub_key", "serv.pub", "")
	pflag.String("endpoint", "localhost:9090", "")

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

	proc := &MainProc{}

	cli, err := rpcli.CreateRPCli().
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
		FutureTimeout(time.Hour).
		MainProc(proc).
		Connect(context.Background(), viper.GetString("endpoint"))
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			var txt string
			fmt.Scanln(&txt)

			ret, err := rpc.Result1[string](<-proc.RPC(serv.Work, comp.CmdCompSelf.Name, "Echo", txt))
			if err != nil {
				fmt.Println("error:", err)
				continue
			}

			fmt.Println("echo:", ret)
		}
	}()

	<-cli.Done()

	if err := context.Cause(cli); err != nil {
		fmt.Println("cause:", err)
	}
}

type MainProc struct {
	rpcli.Proc
}
