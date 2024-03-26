package serv

import (
	"git.golaxy.org/core"
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/util/generic"
	"git.golaxy.org/examples/app/demo_cs/demo_server/comp"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/net/gtp"
	"git.golaxy.org/framework/plugins/gate"
	"git.golaxy.org/framework/plugins/router"
	"git.golaxy.org/framework/plugins/rpc"
	"git.golaxy.org/framework/plugins/rpc/processor"
	"git.golaxy.org/framework/plugins/rpc/rpcutil"
	"time"
)

type GateService struct {
	framework.ServiceBehavior
}

func (serv *GateService) Init(ctx service.Context) {
	core.CreateEntityPT(ctx).
		Prototype("user").
		AddComponent(&comp.UserComp{}).
		Declare()
}

func (serv *GateService) InstallRPC(ctx service.Context) {
	cliPubKey, err := misc.LoadPublicKey(serv.GetStartupConf().GetString("cli_pub_key"))
	if err != nil {
		panic(err)
	}
	_ = cliPubKey

	servPrivKey, err := misc.LoadPrivateKey(serv.GetStartupConf().GetString("serv_priv_key"))
	if err != nil {
		panic(err)
	}
	_ = servPrivKey

	gate.Install(ctx,
		gate.With.Endpoints("0.0.0.0:9090"),
		gate.With.IOTimeout(10*time.Second),
		gate.With.IOBufferCap(1024*1024*5),
		gate.With.AgreeClientEncryptionProposal(true),
		gate.With.AgreeClientCompressionProposal(true),
		gate.With.EncCipherSuite(gtp.CipherSuite{
			SecretKeyExchange:   gtp.SecretKeyExchange_ECDHE,
			SymmetricEncryption: gtp.SymmetricEncryption_AES,
			BlockCipherMode:     gtp.BlockCipherMode_GCM,
			PaddingMode:         gtp.PaddingMode_Pkcs7,
			MACHash:             gtp.Hash_Fnv1a64,
		}),
		gate.With.EncSignatureAlgorithm(gtp.SignatureAlgorithm{
			AsymmetricEncryption: gtp.AsymmetricEncryption_RSA_256,
			PaddingMode:          gtp.PaddingMode_Pkcs1v15,
			Hash:                 gtp.Hash_SHA256,
		}),
		gate.With.EncSignaturePrivateKey(servPrivKey),
		gate.With.EncVerifyClientSignature(true),
		gate.With.EncVerifySignaturePublicKey(cliPubKey),
		gate.With.CompressedSize(128),
		gate.With.SessionInactiveTimeout(time.Hour),
		gate.With.SessionStateChangedHandler(generic.CastDelegateAction3(func(sess gate.ISession, cur, old gate.SessionState) {
			if cur != gate.SessionState_Confirmed {
				return
			}

			rt := framework.CreateRuntime(ctx).Spawn()

			runtime.Concurrent(rt).CallVoid(func(...any) {
				entity, err := core.CreateEntity(rt).
					Prototype("user").
					Scope(ec.Scope_Global).
					Spawn()
				if err != nil {
					panic(err)
				}

				ret := <-rpcutil.ProxyService(ctx, Work).BalanceRPC(EntitiesPluginSelf.Name, "CreateEntity", entity.GetId())
				if !ret.OK() {
					panic(ret.Error)
				}

				err = router.Using(sess.GetContext()).Mapping(entity.GetId(), sess.GetId())
				if err != nil {
					panic(err)
				}

			}).Wait(ctx)
		})),
	)

	router.Install(ctx,
		router.With.CustomAddresses(serv.GetStartupConf().GetString("etcd.address")),
		router.With.CustomAuth(
			serv.GetStartupConf().GetString("etcd.username"),
			serv.GetStartupConf().GetString("etcd.password"),
		),
	)

	rpc.Install(ctx,
		rpc.With.Deliverers(
			processor.NewServiceDeliverer(),
		),
		rpc.With.Dispatchers(
			processor.NewServiceDispatcher(),
			processor.NewOutboundDispatcher(),
			processor.NewInboundDispatcher(gap.DefaultMsgCreator()),
		),
	)
}
