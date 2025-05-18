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
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/examples/app/demo_chat/consts"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/addins/gate"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/addins/router"
	"git.golaxy.org/framework/addins/rpc"
	"git.golaxy.org/framework/addins/rpc/rpcpcsr"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/net/gtp"
	"git.golaxy.org/framework/net/gtp/transport"
	"time"
)

// GateService 网关服务
type GateService struct {
	framework.Service
}

func (s *GateService) Built(svc framework.IService) {
	// 定义用户实体原型
	s.BuildEntityPT(consts.User).
		SetScope(ec.Scope_Global).
		AddComponent(&GateUserComp{}).
		AddComponent(&GateChatChannelComp{}).
		Declare()
}

func (s *GateService) Started(svc framework.IService) {
	if _, err := router.Using(s).AddGroup(s, consts.GlobalChannel); err != nil {
		log.Panicf(s, "create channel %s failed, %s", consts.GlobalChannel, err)
	}
}

func (s *GateService) InstallRPC(svc framework.IService) {
	// 加载客户端签名公钥
	cliPubKey, err := gtp.LoadPublicKeyFile(s.GetAppConf().GetString("cli_pub_key"))
	if err != nil {
		panic(err)
	}

	// 加载服务器签名私钥
	servPrivKey, err := gtp.LoadPrivateKeyFile(s.GetAppConf().GetString("serv_priv_key"))
	if err != nil {
		panic(err)
	}

	// 安装网关插件
	gate.Install(s,
		gate.With.TCPAddress("0.0.0.0:9090"),
		gate.With.WebSocketURL("ws://0.0.0.0:8080"),
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
		gate.With.SessionStateChangedHandler(generic.CastAction3(s.onSessionStateChanged).ToDelegate()),
	)

	// 安装路由插件
	router.Install(s,
		router.With.CustomAddresses(s.GetAppConf().GetString("etcd.address")),
		router.With.CustomAuth(
			s.GetAppConf().GetString("etcd.username"),
			s.GetAppConf().GetString("etcd.password"),
		),
	)

	// 安装RPC插件
	rpc.Install(s,
		rpc.With.Processors(
			rpcpcsr.NewServiceProcessor(nil, true),
			rpcpcsr.NewGateProcessor(gap.DefaultMsgCreator()),
			rpcpcsr.NewForwardProcessor(consts.Gate, gap.DefaultMsgCreator(), generic.CastDelegate2(rpcpcsr.DefaultValidateCliPermission), true),
		),
	)
}

func (s *GateService) onSessionStateChanged(session gate.ISession, curState, lastState gate.SessionState) {
	if curState != gate.SessionState_Confirmed {
		return
	}

	// 创建用户实体
	user, err := s.BuildEntityAsync(consts.User).
		SetPersistId(session.GetId()).
		SetMeta(map[string]any{"session": session}).
		New()
	if err != nil {
		log.Errorf(s, "create gate user %s failed, %s", session.GetId(), err)
		session.Close(&transport.RstError{
			Code:    gtp.Code_Reject,
			Message: err.Error(),
		})
		return
	}

	log.Infof(s, "create gate user %s ok", user.GetId())
}
