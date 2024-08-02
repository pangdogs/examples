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
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/examples/app/demo_cs/demo_server/serv"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework"
	"github.com/spf13/pflag"
)

func main() {
	framework.NewApp().
		Setup(misc.Gate, framework.ServiceGenericT[serv.GateService]{}).
		Setup(misc.Work, framework.ServiceGenericT[serv.WorkService]{}).
		InitCB(generic.MakeDelegateAction1(func(*framework.App) {
			pflag.String("cli_pub_key", "cli.pub", "client public key for verify sign")
			pflag.String("serv_priv_key", "serv.pem", "service private key for sign")
		})).
		Run()
}
