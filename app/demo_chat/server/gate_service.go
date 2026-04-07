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
	"time"

	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/examples/app/demo_chat/consts"
	"git.golaxy.org/examples/app/demo_chat/server/comps"
	"git.golaxy.org/framework"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/gate"
	"git.golaxy.org/framework/addins/rpc/rpcpcsr"
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/net/gtp"
	"git.golaxy.org/framework/net/gtp/transport"
	"go.uber.org/zap"
)

// GateService 网关服务
type GateService struct {
	framework.ServiceBehavior
}

func (s *GateService) OnBuilt(svc framework.IService) {
	// 定义用户实体原型
	s.BuildEntityPT(consts.User).
		SetScope(ec.Scope_Global).
		AddComponent(&comps.GateUserComp{}).
		AddComponent(&comps.GateChatChannelComp{}).
		Declare()
}

func (s *GateService) OnStarted(svc framework.IService) {
	_, err := Gate.Require(s).Watch(s, generic.CastDelegateVoid1(s.handleSessionEstablished))
	if err != nil {
		s.L().Panic("watch session established failed", zap.Error(err))
	}

	group, err := Router.Require(s).AddGroup(s, consts.GlobalChannel, nil, 15*time.Second)
	if err != nil {
		s.L().Panic("create channel failed", zap.String("channel", consts.GlobalChannel), zap.Error(err))
	}
	_, err = group.KeepAliveContinuous(s.Terminated().Context(nil))
	if err != nil {
		s.L().Panic("keep alive channel failed", zap.String("channel", consts.GlobalChannel), zap.Error(err))
	}
}

func (s *GateService) InstallRPC(svc framework.IService) {
	// 加载客户端签名公钥
	cliPubKey, err := gtp.LoadPublicKeyFile(s.AppConf().GetString("cli_pub_key"))
	if err != nil {
		s.L().Panic("load cli public key failed", zap.Error(err))
	}

	// 加载服务器签名私钥
	servPrivKey, err := gtp.LoadPrivateKeyFile(s.AppConf().GetString("serv_priv_key"))
	if err != nil {
		s.L().Panic("load serv private key failed", zap.Error(err))
	}

	// 安装网关插件
	Gate.Install(s,
		GateWith.TCPAddress("0.0.0.0:9090"),
		GateWith.WebSocketURL("ws://0.0.0.0:8080"),
		GateWith.IOTimeout(10*time.Second),
		GateWith.IOBufferCap(1024*1024*5),
		GateWith.AgreeClientEncryptionProposal(true),
		GateWith.AgreeClientCompressionProposal(true),
		GateWith.EncCipherSuite(gtp.CipherSuite{
			SecretKeyExchange:   gtp.SecretKeyExchange_ECDHE,
			SymmetricEncryption: gtp.SymmetricEncryption_XChaCha20_Poly1305,
		}),
		GateWith.EncSignatureAlgorithm(gtp.SignatureAlgorithm{
			AsymmetricEncryption: gtp.AsymmetricEncryption_RSA,
			PaddingMode:          gtp.PaddingMode_Pkcs1v15,
			Hash:                 gtp.Hash_SHA256,
		}),
		GateWith.EncSignaturePrivateKey(servPrivKey),
		GateWith.EncVerifyClientSignature(true),
		GateWith.EncVerifySignaturePublicKey(cliPubKey),
		GateWith.CompressionThreshold(128),
		GateWith.SessionInactiveTimeout(15*time.Second),
	)

	// 安装路由插件
	Router.Install(s,
		RouterWith.CustomAddresses(s.AppConf().GetString("etcd.address")),
		RouterWith.CustomAuth(
			s.AppConf().GetString("etcd.username"),
			s.AppConf().GetString("etcd.password"),
		),
	)

	// 安装RPC插件
	RPC.Install(s,
		RPCWith.Processors(
			rpcpcsr.NewServiceProcessor(nil, true),
			rpcpcsr.NewGateProcessor(gap.DefaultMsgCreator()),
			rpcpcsr.NewForwardProcessor(consts.Gate, gap.DefaultMsgCreator(), generic.CastDelegate2(rpcpcsr.DefaultValidateCliPermission), true),
		),
	)
}

func (s *GateService) handleSessionEstablished(session gate.ISession) {
	// 创建用户实体
	user, err := s.BuildEntity(consts.User).
		SetPersistId(session.Id()).
		New()
	if err != nil {
		s.L().Panic("create user failed", zap.Any("session", session), zap.Error(err))
		session.Close(&transport.RstError{
			Code:    gtp.Code_Reject,
			Message: err.Error(),
		})
		return
	}
	s.L().Info("user created", zap.Any("session", session), zap.Any("user", user))
}
