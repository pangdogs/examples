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
	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
)

// DemoComp Demo组件实现
type DemoComp struct {
	ec.ComponentBehavior
}

// Awake 组件唤醒
func (comp *DemoComp) Awake() {
	fmt.Printf("I'm entity %q, comp %q Awake.\n", comp.GetEntity(), comp)
}

// Start 组件开始
func (comp *DemoComp) Start() {
	fmt.Printf("I'm entity %q, comp %q Start.\n", comp.GetEntity(), comp)
}

// Update 组件更新
func (comp *DemoComp) Update() {
	frame := runtime.Current(comp).GetFrame()

	if frame.GetCurFrames()%int64(frame.GetTargetFPS()) == 0 {
		fmt.Printf("I'm entity %q, comp %q Update(%s).\n", comp.GetEntity(), comp, frame.GetRunningElapseTime())
	}
}

// LateUpdate 组件滞后更新
func (comp *DemoComp) LateUpdate() {
	frame := runtime.Current(comp).GetFrame()

	if frame.GetCurFrames()%int64(frame.GetTargetFPS()) == 0 {
		fmt.Printf("I'm entity %q, comp %q LateUpdate(%s).\n", comp.GetEntity(), comp, frame.GetRunningElapseTime())
	}
}

// Shut 组件停止
func (comp *DemoComp) Shut() {
	fmt.Printf("I'm entity %q, comp %q Shut.\n", comp.GetEntity(), comp)
}

// Dispose 组件销毁
func (comp *DemoComp) Dispose() {
	fmt.Printf("I'm entity %q, comp %q Dispose.\n", comp.GetEntity(), comp)
}
