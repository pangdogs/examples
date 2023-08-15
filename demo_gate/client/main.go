package main

import (
	"context"
	"fmt"
	"kit.golaxy.org/plugins/gate/gtp_client"
	"kit.golaxy.org/plugins/transport"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		panic("missing endpoint")
	}

	cli, err := gtp_client.Connect(context.Background(), os.Args[1],
		gtp_client.Option{}.RecvDataHandlers(func(client *gtp_client.Client, data []byte, sequenced bool) error {
			fmt.Println(string(data))
			return nil
		}),
		gtp_client.Option{}.EncCipherSuite(transport.CipherSuite{
			SecretKeyExchange:   transport.SecretKeyExchange_ECDHE,
			SymmetricEncryption: transport.SymmetricEncryption_XChaCha20,
			BlockCipherMode:     transport.BlockCipherMode_None,
			PaddingMode:         transport.PaddingMode_None,
			MACHash:             transport.Hash_Fnv1a64,
		}))
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	fmt.Println("this console is", cli.GetSessionId())

	for {
		var text string
		fmt.Scanln(&text)
		if err := cli.SendData([]byte(text), true); err != nil {
			panic(err)
		}
	}
}
