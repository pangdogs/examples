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
	"log"

	"git.golaxy.org/core/ec"
	"git.golaxy.org/core/runtime"
)

// HelloWorldComp HelloWorld组件
type HelloWorldComp struct {
	ec.ComponentBehavior
}

// Awake 组件唤醒
func (comp *HelloWorldComp) Awake() {
	log.Printf("[%s] Awake", comp.Entity().Id())
}

// OnEnable 组件启用
func (comp *HelloWorldComp) OnEnable() {
	log.Printf("[%s] OnEnable", comp.Entity().Id())
}

// Start 组件开始
func (comp *HelloWorldComp) Start() {
	log.Printf("[%s] Start", comp.Entity().Id())
}

// Update 组件更新
func (comp *HelloWorldComp) Update() {
	frame := runtime.Current(comp).Frame()
	log.Printf("[%s] Update, frame %d, last loop elapse %fs", comp.Entity().Id(), frame.CurFrames(), frame.LastLoopElapseTime().Seconds())
}

// LateUpdate 组件滞后更新
func (comp *HelloWorldComp) LateUpdate() {
	frame := runtime.Current(comp).Frame()
	log.Printf("[%s] Late Update, frame %d, last loop elapse %fs", comp.Entity().Id(), frame.CurFrames(), frame.LastLoopElapseTime().Seconds())
}

// Shut 组件停止
func (comp *HelloWorldComp) Shut() {
	log.Printf("[%s] Shut", comp.Entity().Id())
}

// OnDisable 组件关闭
func (comp *HelloWorldComp) OnDisable() {
	log.Printf("[%s] OnDisable", comp.Entity().Id())
}

// Dispose 组件销毁
func (comp *HelloWorldComp) Dispose() {
	log.Printf("[%s] Dispose", comp.Entity().Id())
}
