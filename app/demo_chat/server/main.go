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
	"git.golaxy.org/examples/app/demo_chat/misc"
	"git.golaxy.org/framework"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

/*
 * 基于framework层提供的支持，演示一个简单的匿名聊天系统，支持加入多个群组。
 */
func main() {
	framework.NewApp().
		Setup(misc.Gate, &GateService{}).
		Setup(misc.Chat, &ChatService{}).
		InitCB(func(*cobra.Command) {
			pflag.String("cli_pub_key", "cli.pub", "client public key for verify sign")
			pflag.String("serv_priv_key", "serv.pem", "service private key for sign")
		}).
		Run()
}
