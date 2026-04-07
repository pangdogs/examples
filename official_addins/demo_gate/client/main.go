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
	"os"
	"time"

	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/framework/addins/gate/cli"
	"git.golaxy.org/framework/net/gtp"
	"go.uber.org/zap"
)

func main() {
	endpoint := "localhost:9090"

	if len(os.Args) > 1 {
		endpoint = os.Args[1]
	}

	logger, _ := zap.NewDevelopment()

	client, err := cli.Connect(context.Background(), endpoint,
		cli.With.EncCipherSuite(gtp.CipherSuite{
			SecretKeyExchange:   gtp.SecretKeyExchange_ECDHE,
			SymmetricEncryption: gtp.SymmetricEncryption_ChaCha20_Poly1305,
			BlockCipherMode:     gtp.BlockCipherMode_None,
			PaddingMode:         gtp.PaddingMode_None,
			HMAC:                gtp.Hash_None,
		}),
		cli.With.CompressionThreshold(128),
		cli.With.IOTimeout(3*time.Second),
		cli.With.IOBufferCap(1024*1024*5),
		cli.With.AutoReconnect(true),
		cli.With.Logger(logger),
	)
	if err != nil {
		logger.Panic("connect failed", zap.Error(err))
	}
	defer client.Close(nil)

	err = client.DataIO().Listen(nil, generic.CastDelegateVoid1(func(data []byte) {
		client.Logger().Info("[echo]", zap.String("text", string(data)))
	}))
	if err != nil {
		client.Logger().Panic("listen data failed", zap.Error(err))
	}

	for {
		future := <-client.RequestTime().Chan()
		if future.Error != nil {
			client.Logger().Panic("sync time failed", zap.Error(future.Error))
		}

		respTime := future.Value.(*cli.ResponseTime)

		client.Logger().Info("sync time",
			zap.Time("time", respTime.NowTime()),
			zap.Duration("rtt", respTime.RTT()))

		var text string
		fmt.Scanln(&text)
		if err := client.DataIO().Send([]byte(text)); err != nil {
			client.Logger().Panic("send data failed", zap.Error(err))
		}
		client.Logger().Info("[send]", zap.String("text", text))
	}
}
