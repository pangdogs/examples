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
		client.L().Info("[echo]", zap.String("text", string(data)))
	}))
	if err != nil {
		client.L().Panic("listen data failed", zap.Error(err))
	}

	for {
		result := client.ProbeTime().Wait(client)
		if result.Error != nil {
			client.L().Panic("probe time failed", zap.Error(result.Error))
		}

		timeSample := result.Value.(*cli.TimeSample)

		client.L().Info("time sample",
			zap.Time("remote_time", timeSample.RemoteTime()),
			zap.Duration("rtt", timeSample.RTT()),
			zap.Duration("offset", timeSample.Offset()))

		var text string
		fmt.Scanln(&text)
		if err := client.DataIO().Send([]byte(text)); err != nil {
			client.L().Panic("send data failed", zap.Error(err))
		}
		client.L().Info("[send]", zap.String("text", text))
	}
}
