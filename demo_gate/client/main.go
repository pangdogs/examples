package main

import (
	"context"
	"fmt"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/plugins/gtp"
	"git.golaxy.org/plugins/gtp_cli"
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

	cli, err := gtp_cli.Connect(context.Background(), os.Args[1],
		gtp_cli.Option{}.RecvDataHandler(generic.CastDelegateFunc1(func(data []byte) error {
			log.Infoln(string(data))
			return nil
		})),
		gtp_cli.Option{}.EncCipherSuite(gtp.CipherSuite{
			SecretKeyExchange:   gtp.SecretKeyExchange_ECDHE,
			SymmetricEncryption: gtp.SymmetricEncryption_XChaCha20,
			BlockCipherMode:     gtp.BlockCipherMode_None,
			PaddingMode:         gtp.PaddingMode_None,
			MACHash:             gtp.Hash_Fnv1a64,
		}),
		gtp_cli.Option{}.CompressedSize(128),
		gtp_cli.Option{}.IOTimeout(3*time.Second),
		gtp_cli.Option{}.IOBufferCap(1024*1024*5),
		gtp_cli.Option{}.AutoReconnect(true),
		gtp_cli.Option{}.ZapLogger(zaplogger),
	)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

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
