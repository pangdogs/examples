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
	"git.golaxy.org/framework/net/gap"
	"git.golaxy.org/framework/net/gap/variant"
	"git.golaxy.org/framework/utils/binaryutil"
)

func init() {
	gap.DefaultMsgCreator().Declare(&MsgDemo{})
}

const (
	MsgId_Demo = gap.MsgId_Customize + iota
)

type MsgDemo struct {
	Int    int
	Double float64
	Str    string
	Map    variant.Map
	Array  variant.Array
}

// Read implements io.Reader
func (m MsgDemo) Read(p []byte) (int, error) {
	bs := binaryutil.NewBigEndianStream(p)
	if err := bs.WriteVarint(int64(m.Int)); err != nil {
		return bs.BytesWritten(), err
	}
	if err := bs.WriteDouble(m.Double); err != nil {
		return bs.BytesWritten(), err
	}
	if err := bs.WriteString(m.Str); err != nil {
		return bs.BytesWritten(), err
	}

	rn := bs.BytesWritten()

	n, err := m.Map.Read(p[rn:])
	rn += n
	if err != nil {
		return rn, nil
	}

	n, err = m.Array.Read(p[rn:])
	rn += n
	if err != nil {
		return rn, nil
	}

	return rn, nil
}

// Write implements io.Writer
func (m *MsgDemo) Write(p []byte) (int, error) {
	bs := binaryutil.NewBigEndianStream(p)

	i, err := bs.ReadVarint()
	if err != nil {
		return bs.BytesRead(), err
	}

	d, err := bs.ReadDouble()
	if err != nil {
		return bs.BytesRead(), err
	}

	str, err := bs.ReadString()
	if err != nil {
		return bs.BytesRead(), err
	}

	wn := bs.BytesRead()

	var vmap variant.Map
	n, err := vmap.Write(p[wn:])
	wn += n
	if err != nil {
		return wn, err
	}

	var arr variant.Array
	n, err = arr.Write(p[wn:])
	wn += n
	if err != nil {
		return wn, err
	}

	m.Int = int(i)
	m.Double = d
	m.Str = str
	m.Map = vmap
	m.Array = arr

	return wn, nil
}

// Size 大小
func (m MsgDemo) Size() int {
	return binaryutil.SizeofVarint(int64(m.Int)) + binaryutil.SizeofDouble() + binaryutil.SizeofString(m.Str) + m.Map.Size() + m.Array.Size()
}

// MsgId 消息Id
func (MsgDemo) MsgId() gap.MsgId {
	return MsgId_Demo
}
