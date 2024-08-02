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

package comp

import (
	"git.golaxy.org/core/define"
	"git.golaxy.org/examples/app/demo_cs/misc"
	"git.golaxy.org/framework"
	"git.golaxy.org/framework/plugins/rpc/rpcutil"
)

var UserCompSelf = define.Component[UserComp]()

type UserComp struct {
	framework.ComponentBehavior
}

func (c *UserComp) Dispose() {
	if c.GetService().GetName() == misc.Gate {
		c.RPC(misc.Work, rpcutil.NoComp, "DestroySelf")
	}
}
