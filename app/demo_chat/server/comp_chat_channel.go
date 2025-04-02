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
	"fmt"
	"git.golaxy.org/examples/app/demo_chat/misc"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/addins/router"
	"git.golaxy.org/framework/addins/rpc/rpcli"
	"git.golaxy.org/framework/addins/rpc/rpcutil"
	"time"
)

type ChatChannelComp struct {
	framework.ComponentBehavior
}

func (c *ChatChannelComp) Start() {
	channel, ok := router.Using(c.GetService()).GetGroup(c, misc.GlobalChannel)
	if !ok {
		log.Panicf(c, "gate user %s get channel %s failed", c.GetId(), misc.GlobalChannel)
	}

	if err := channel.Add(c, c.GetId()); err != nil {
		log.Panicf(c, "gate user %s join channel %s failed, %s", c.GetId(), misc.GlobalChannel, err)
	}
}

func (c *ChatChannelComp) Shut() {
	router.Using(c.GetService()).EachGroups(nil, c.GetId(), func(channel router.IGroup) {
		channel.Remove(nil, c.GetId())
	})
}

func (c *ChatChannelComp) C_CreateChannel(channelName string) {
	if channelName == misc.GlobalChannel {
		return
	}

	_, err := router.Using(c.GetService()).AddGroup(c, channelName)
	if err != nil {
		log.Errorf(c, "gate user %s create channel %s failed, %s", c.GetId(), channelName, err)
		return
	}

	c.SendToChannel(misc.GlobalChannel, fmt.Sprintf("create channel %s", channelName))
	c.C_JoinChannel(channelName)
}

func (c *ChatChannelComp) C_RemoveChannel(channelName string) {
	if channelName == misc.GlobalChannel {
		return
	}

	_, ok := router.Using(c.GetService()).GetGroup(c, channelName)
	if !ok {
		log.Errorf(c, "gate user %s get channel %s failed", c.GetId(), channelName)
		return
	}
	router.Using(c.GetService()).DeleteGroup(c, channelName)

	c.SendToChannel(misc.GlobalChannel, fmt.Sprintf("remove channel %s", channelName))
}

func (c *ChatChannelComp) C_JoinChannel(channelName string) {
	if channelName == misc.GlobalChannel {
		return
	}

	channel, ok := router.Using(c.GetService()).GetGroup(c, channelName)
	if !ok {
		log.Errorf(c, "gate user %s get channel %s failed", c.GetId(), channelName)
		return
	}
	if err := channel.Add(c, c.GetId()); err != nil {
		log.Errorf(c, "gate user %s join channel %s failed, %s", c.GetId(), channelName, err)
		return
	}

	c.SendToChannel(channelName, "joined")
}

func (c *ChatChannelComp) C_LeaveChannel(channelName string) {
	if channelName == misc.GlobalChannel {
		return
	}

	channel, ok := router.Using(c.GetService()).GetGroup(c, channelName)
	if !ok {
		log.Errorf(c, "gate user %s get channel %s failed", c.GetId(), channelName)
		return
	}
	if err := channel.Remove(c, c.GetId()); err != nil {
		log.Errorf(c, "gate user %s join channel %s failed, %s", c.GetId(), channelName, err)
		return
	}

	c.SendToChannel(channelName, "leaved")
}

func (c *ChatChannelComp) SendToChannel(channelName, text string) {
	err := rpcutil.ProxyGroup(c, channelName).CliOnewayRPC(rpcli.Main, "OutputText", time.Now().Unix(), channelName, c.GetId(), text)
	if err != nil {
		log.Errorf(c, "gate user %s send %q to channel %s failed, %s", c.GetId(), text, channelName, err)
		return
	}
	log.Infof(c, "gate user %s send %q to channel %s ok", c.GetId(), text, channelName)
}
