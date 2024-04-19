package main

import (
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/examples/app/demo_cs/demo_server/serv"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework"
	"github.com/spf13/pflag"
)

func main() {
	framework.NewApp().
		Setup(misc.Gate, &serv.GateService{}).
		Setup(misc.Work, &serv.WorkService{}).
		InitCB(generic.MakeDelegateAction1(func(*framework.App) {
			pflag.String("cli_pub_key", "cli.pub", "client public key for verify sign")
			pflag.String("serv_priv_key", "serv.pem", "service private key for sign")
		})).
		Run()
}
