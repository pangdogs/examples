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

package behavior

import (
	"git.golaxy.org/examples/app/demo_chat/consts"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/addins/log"
	"git.golaxy.org/framework/addins/rpc"
)

type ChatUserComp struct {
	framework.ComponentBehavior
}

func (c *ChatUserComp) C_InputText(channelName, text string) {
	if err := rpc.ResultVoid(<-c.RPC(consts.Gate, "GateChatChannelComp", "SendToChannel", channelName, text)).Extract(); err != nil {
		log.Errorf(c, "chat user %s send %q to channel %s failed, %s", c.GetId(), text, channelName, err)
		return
	}
	log.Infof(c, "chat user %s send %q to channel %s ok", c.GetId(), text, channelName)
}

func (c *ChatUserComp) Dispose() {
	log.Infof(c, "chat user %s destroyed", c.GetId())
}
