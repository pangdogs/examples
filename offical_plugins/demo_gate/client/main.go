package main

import (
	"context"
	"fmt"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/framework/net/gtp"
	"git.golaxy.org/framework/plugins/gate/cli"
	"go.uber.org/zap"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		panic("missing endpoint")
	}

	zaplogger, _ := zap.NewProduction()
	log := zaplogger.Sugar()

	cli, err := cli.Connect(context.Background(), os.Args[1],
		cli.With.RecvDataHandler(generic.MakeDelegateFunc1(func(data []byte) error {
			log.Infoln(string(data))
			return nil
		})),
		cli.With.EncCipherSuite(gtp.CipherSuite{
			SecretKeyExchange:   gtp.SecretKeyExchange_ECDHE,
			SymmetricEncryption: gtp.SymmetricEncryption_XChaCha20,
			BlockCipherMode:     gtp.BlockCipherMode_None,
			PaddingMode:         gtp.PaddingMode_None,
			MACHash:             gtp.Hash_Fnv1a64,
		}),
		cli.With.CompressedSize(128),
		cli.With.IOTimeout(3*time.Second),
		cli.With.IOBufferCap(1024*1024*5),
		cli.With.AutoReconnect(true),
		cli.With.ZapLogger(zaplogger),
	)
	if err != nil {
		panic(err)
	}
	defer cli.Close(nil)

	log.Infoln("this console is", cli.GetSessionId())

	for {
		respTime := <-cli.RequestTime(context.Background())
		if respTime.Error != nil {
			log.Infof("sync time: %s", respTime.Error)
		} else {
			log.Infof("sync time: %s, rtt: %s, raw: %+v\n", respTime.Value.SyncTime(), respTime.Value.RTT(), respTime.Value)
		}
		var text string
		fmt.Scanln(&text)
		if err := cli.SendData([]byte(text)); err != nil {
			log.Infoln("send data err:", err)
		}
	}
}
