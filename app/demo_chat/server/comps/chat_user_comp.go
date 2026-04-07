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
	"git.golaxy.org/examples/app/demo_chat/consts"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/addins/rpc"
	"go.uber.org/zap"
)

type ChatUserComp struct {
	framework.ComponentBehavior
}

func (c *ChatUserComp) C_InputText(channelName, text string) {
	if err := rpc.ResultVoid(c.RPC(consts.Gate, "GateChatChannelComp", "SendToChannel", channelName, text)).Error; err != nil {
		c.L().Error("send text to channel failed",
			zap.String("channel", channelName),
			zap.String("text", text),
			zap.Error(err))
		return
	}
	c.L().Info("[send]", zap.String("channel", channelName), zap.String("text", text))
}

func (c *ChatUserComp) Dispose() {
	c.L().Info("user disposed")
	c.Runtime().Terminate()
}
