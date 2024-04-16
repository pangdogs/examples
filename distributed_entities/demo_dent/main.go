package main

import (
	"fmt"
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/plugin"
	"git.golaxy.org/core/pt"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/core/util/uid"
	"git.golaxy.org/framework/plugins/broker/nats_broker"
	"git.golaxy.org/framework/plugins/dent"
	"git.golaxy.org/framework/plugins/dentq"
	"git.golaxy.org/framework/plugins/discovery/etcd_discovery"
	"git.golaxy.org/framework/plugins/dserv"
	"git.golaxy.org/framework/plugins/dsync/etcd_dsync"
	"git.golaxy.org/framework/plugins/log"
	"git.golaxy.org/framework/plugins/log/console_log"
	"git.golaxy.org/framework/plugins/rpc"
	"github.com/segmentio/ksuid"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	entIdList := []ksuid.KSUID{ksuid.New()}
	var servChanList []<-chan struct{}
	total := 6

	for i := 0; i < total; i++ {
		// 创建实体库，注册实体原型
		entityLib := pt.NewEntityLib(pt.DefaultComponentLib())
		entityLib.Register("demo", pt.CompAlias(DemoComp{}, "DemoComp"))

		// 创建插件包，安装插件
		pluginBundle := plugin.NewPluginBundle()
		console_log.Install(pluginBundle, console_log.With.Level(log.DebugLevel), console_log.With.ServiceInfo(true))
		nats_broker.Install(pluginBundle, nats_broker.With.CustomAddresses("127.0.0.1:4222"))
		etcd_discovery.Install(pluginBundle, etcd_discovery.With.CustomAddresses("127.0.0.1:12379"))
		etcd_dsync.Install(pluginBundle, etcd_dsync.With.CustomAddresses("127.0.0.1:12379"))
		dserv.Install(pluginBundle, dserv.With.FutureTimeout(time.Minute))
		rpc.Install(pluginBundle)
		dentq.Install(pluginBundle, dentq.With.CustomAddresses("127.0.0.1:12379"))

		// 创建服务上下文与服务，并开始运行
		serv := core.NewService(service.NewContext(
			service.With.EntityLib(entityLib),
			service.With.PluginBundle(pluginBundle),
			service.With.Name(fmt.Sprintf("demo_dent_%d", i%(total/2))),
			service.With.RunningHandler(generic.MakeDelegateAction2(func(ctx service.Context, state service.RunningState) {
				if state != service.RunningState_Started {
					return
				}

				// 监听退出信号
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

				go func() {
					<-sigChan
					ctx.GetCancelFunc()()
				}()

				// 创建运行时上下文与运行时，安装插件并开始运行
				rtCtx := runtime.NewContext(ctx,
					runtime.With.Context.RunningHandler(generic.MakeDelegateAction2(func(_ runtime.Context, state runtime.RunningState) {
						if state != runtime.RunningState_Terminated {
							return
						}
						ctx.GetCancelFunc()()
					})),
				)
				console_log.Install(rtCtx, console_log.With.Level(log.DebugLevel), console_log.With.ServiceInfo(true))
				dent.Install(rtCtx, dent.With.CustomAddresses("127.0.0.1:12379"))

				rt := core.NewRuntime(rtCtx, core.With.Runtime.AutoRun(true))

				// 在运行时线程环境中，创建实体
				for i := range entIdList {
					entId := entIdList[i]

					core.AsyncVoid(rt, func(ctx runtime.Context, _ ...any) {
						entity, err := core.CreateEntity(ctx).
							Prototype("demo").
							Scope(ec.Scope_Global).
							PersistId(uid.Id(entId.String())).
							Spawn()
						if err != nil {
							log.Panic(service.Current(ctx), err)
						}
						log.Infof(service.Current(ctx), "create entity %q finish", entity)
					})
				}
			})),
		))

		servChanList = append(servChanList, serv.Run())
	}

	for _, servChan := range servChanList {
		<-servChan
	}
}
