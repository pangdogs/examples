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

package comps

import (
	"context"
	"fmt"
	"time"

	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/examples/app/demo_chat/consts"
	"git.golaxy.org/framework"
	. "git.golaxy.org/framework/addins"
	"git.golaxy.org/framework/addins/rpc"
	"go.uber.org/zap"
)

type GateChatChannelComp struct {
	framework.ComponentBehavior
	createdGroups generic.SliceMap[string, time.Time]
}

func (c *GateChatChannelComp) Start() {
	c.JoinChannel(consts.GlobalChannel)
}

func (c *GateChatChannelComp) Shut() {
	c.createdGroups.Each(func(channelName string, createdTime time.Time) {
		c.C_RemoveChannel(channelName)
	})
	for _, group := range Router.Require(c.Service()).GetGroupsByEntity(context.Background(), c.Id()) {
		c.C_LeaveChannel(group.Name())
	}
}

func (c *GateChatChannelComp) C_CreateChannel(channelName string) {
	if channelName == consts.GlobalChannel {
		return
	}

	group, err := Router.Require(c.Service()).AddGroup(c.Entity(), channelName, nil, 15*time.Second)
	if err != nil {
		c.L().Panic("create channel failed", zap.String("channel", channelName), zap.Error(err))
	}
	if _, err = group.KeepAliveContinuous(c.Entity()); err != nil {
		c.L().Panic("keep alive channel failed", zap.String("channel", channelName), zap.Error(err))
	}
	c.L().Info("channel created", zap.String("channel", channelName))

	c.createdGroups.Add(channelName, time.Now())

	c.SendToChannel(consts.GlobalChannel, fmt.Sprintf("channel %s created", channelName))
	c.C_JoinChannel(channelName)
}

func (c *GateChatChannelComp) C_RemoveChannel(channelName string) {
	if channelName == consts.GlobalChannel {
		return
	}

	if !c.createdGroups.Exist(channelName) {
		return
	}

	if _, ok := Router.Require(c.Service()).GetGroupByName(context.Background(), channelName); !ok {
		c.L().Error("channel not found", zap.String("channel", channelName))
		return
	}

	c.SendToChannel(consts.GlobalChannel, fmt.Sprintf("channel %s removed", channelName))
	rpc.ProxyGroup(c, channelName).CliOnewayRPC("", "ChannelKickOut", channelName)

	Router.Require(c.Service()).DeleteGroup(context.Background(), channelName)
	c.createdGroups.Delete(channelName)
}

func (c *GateChatChannelComp) C_JoinChannel(channelName string) {
	if channelName == consts.GlobalChannel {
		return
	}
	c.JoinChannel(channelName)
}

func (c *GateChatChannelComp) C_LeaveChannel(channelName string) {
	if channelName == consts.GlobalChannel {
		return
	}

	group, ok := Router.Require(c.Service()).GetGroupByName(context.Background(), channelName)
	if !ok {
		c.L().Error("channel not found", zap.String("channel", channelName))
		return
	}

	c.SendToChannel(channelName, "leaved")
	c.CliOnewayRPC("", "ChannelKickOut", channelName)

	if err := group.Remove(context.Background(), []uid.Id{c.Id()}); err != nil {
		c.L().Error("leave channel failed", zap.String("channel", channelName), zap.Error(err))
		return
	}

	c.L().Info("channel leaved", zap.String("channel", channelName))
}

func (c *GateChatChannelComp) C_InChannel(channelName string) bool {
	for _, group := range Router.Require(c.Service()).GetGroupsByEntity(c.Entity(), c.Id()) {
		if group.Name() == channelName {
			return true
		}
	}
	return false
}

func (c *GateChatChannelComp) SendToChannel(channelName, text string) {
	err := rpc.ProxyGroup(c, channelName).CliOnewayRPC("", "OutputText", time.Now().Unix(), channelName, c.Id(), text)
	if err != nil {
		c.L().Error("send to channel failed", zap.String("channel", channelName), zap.String("text", text), zap.Error(err))
		return
	}
	c.L().Info("send to channel ok", zap.String("channel", channelName), zap.String("text", text))
}

func (c *GateChatChannelComp) JoinChannel(channelName string) {
	if c.C_InChannel(channelName) {
		return
	}

	channel, ok := Router.Require(c.Service()).GetGroupByName(c.Entity(), channelName)
	if !ok {
		c.L().Error("channel not found", zap.String("channel", channelName))
		return
	}

	if err := channel.Add(c.Entity(), []uid.Id{c.Id()}); err != nil {
		c.L().Error("join channel failed", zap.String("channel", channelName), zap.Error(err))
		return
	}

	c.L().Info("join channel ok", zap.String("channel", channelName))

	c.Await(c.TimeAfterAsync(time.Second)).AnyVoid(func(framework.IRuntime, async.Result, ...any) {
		c.SendToChannel(channelName, "joined")
	})
}
