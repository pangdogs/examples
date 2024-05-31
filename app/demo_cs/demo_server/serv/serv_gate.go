package serv

import (
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
	"git.golaxy.org/core/service"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/examples/app/demo_cs/demo_server/comp"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/net/gtp"
	"git.golaxy.org/framework/plugins/gate"
	"git.golaxy.org/framework/plugins/log"
	"git.golaxy.org/framework/plugins/router"
	"git.golaxy.org/framework/plugins/rpc"
	"git.golaxy.org/framework/plugins/rpc/rpcpcsr"
	"git.golaxy.org/framework/plugins/rpc/rpcutil"
	"time"
)

// GateService 网关服务
type GateService struct {
	framework.ServiceInstance
}

func (serv *GateService) Built(ctx service.Context) {
	// 定义User实体原型
	serv.CreateEntityPT(misc.User).
		AddComponent(&comp.UserComp{}).
		Scope(ec.Scope_Global).
		Declare()
}

func (serv *GateService) InstallRPC(ctx service.Context) {
	// 加载客户端签名公钥
	cliPubKey, err := gtp.LoadPublicKeyFile(serv.GetStartupConf().GetString("cli_pub_key"))
	if err != nil {
		panic(err)
	}
	_ = cliPubKey

	// 加载服务器签名私钥
	servPrivKey, err := gtp.LoadPrivateKeyFile(serv.GetStartupConf().GetString("serv_priv_key"))
	if err != nil {
		panic(err)
	}
	_ = servPrivKey

	// 安装网关插件
	gate.Install(ctx,
		gate.With.TCPAddress("0.0.0.0:9090"),
		gate.With.WebSocketURL("http://0.0.0.0:8080"),
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
			AsymmetricEncryption: gtp.AsymmetricEncryption_RSA256,
			PaddingMode:          gtp.PaddingMode_Pkcs1v15,
			Hash:                 gtp.Hash_SHA256,
		}),
		gate.With.EncSignaturePrivateKey(servPrivKey),
		gate.With.EncVerifyClientSignature(true),
		gate.With.EncVerifySignaturePublicKey(cliPubKey),
		gate.With.CompressedSize(128),
		gate.With.SessionInactiveTimeout(time.Minute),
		gate.With.SessionStateChangedHandler(generic.MakeAction3(serv.sessionChanged).CastDelegate()),
	)

	// 安装路由插件
	router.Install(ctx,
		router.With.CustomAddresses(serv.GetStartupConf().GetString("etcd.address")),
		router.With.CustomAuth(
			serv.GetStartupConf().GetString("etcd.username"),
			serv.GetStartupConf().GetString("etcd.password"),
		),
	)

	// 安装RPC插件
	rpc.Install(ctx,
		rpc.With.Processors(
			rpcpcsr.NewServiceProcessor(nil),
			rpcpcsr.NewGateProcessor(gap.DefaultMsgCreator()),
		),
	)
}

func (serv *GateService) sessionChanged(sess gate.ISession, cur, old gate.SessionState) {
	if cur != gate.SessionState_Confirmed {
		return
	}

	// 创建用户实体
	user, err := serv.CreateConcurrentEntity(misc.User).Spawn()
	if err != nil {
		log.Panicln(serv, err)
	}

	// 调用工作服唤醒用户
	err = rpcutil.ProxyService(serv, misc.Work).
		BalanceRPC(rpcutil.NoComp, "WakeUpUser", user.GetId()).
		Wait(serv).Error
	if err != nil {
		log.Panicln(serv, err)
	}

	// 映射路由
	mapping, err := router.Using(sess.GetContext()).Mapping(user.GetId(), sess.GetId())
	if err != nil {
		log.Panicln(serv, err)
	}

	log.Infof(serv, "create user %q, sessionId:%q", user.GetId(), sess.GetId())

	go func() {
		<-mapping.Done()
		<-runtime.Concurrent(user).Terminate()
		log.Infof(serv, "destroy user %q, sessionId:%q", user.GetId(), sess.GetId())
	}()
}
