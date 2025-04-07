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
	"context"
	"fmt"
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/examples/app/demo_chat/misc"
	"git.golaxy.org/framework/addins/gate/cli"
	"git.golaxy.org/framework/addins/rpc"
	"git.golaxy.org/framework/addins/rpc/rpcli"
	"git.golaxy.org/framework/net/gtp"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lucasb-eyer/go-colorful"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"hash/fnv"
	"strings"
	"time"
)

func main() {
	pflag.String("cli_priv_key", "cli.pem", "client private key for sign")
	pflag.String("serv_pub_key", "serv.pub", "service public key for verify sign")
	pflag.String("endpoint", "localhost:9090", "connect endpoint")
	pflag.Bool("ws", false, "use websocket")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	cliPrivKey, err := gtp.LoadPrivateKeyFile(viper.GetString("cli_priv_key"))
	if err != nil {
		panic(err)
	}

	servPubKey, err := gtp.LoadPublicKeyFile(viper.GetString("serv_pub_key"))
	if err != nil {
		panic(err)
	}

	np := cli.TCP
	if viper.GetBool("ws") {
		np = cli.WebSocket
	}

	logger := zap.NewNop()
	proc := &MainProc{}

	rpcli, err := rpcli.BuildRPCli().
		SetNetProtocol(np).
		SetIOTimeout(10*time.Second).
		SetGTPAutoReconnect(true).
		SetGTPEncCipherSuite(gtp.CipherSuite{
			SecretKeyExchange:   gtp.SecretKeyExchange_ECDHE,
			SymmetricEncryption: gtp.SymmetricEncryption_AES,
			BlockCipherMode:     gtp.BlockCipherMode_GCM,
			PaddingMode:         gtp.PaddingMode_Pkcs7,
			MACHash:             gtp.Hash_Fnv1a64,
		}).
		SetGTPEncSignatureAlgorithm(gtp.SignatureAlgorithm{
			AsymmetricEncryption: gtp.AsymmetricEncryption_RSA256,
			PaddingMode:          gtp.PaddingMode_Pkcs1v15,
			Hash:                 gtp.Hash_SHA256,
		}).
		SetGTPEncSignaturePrivateKey(cliPrivKey).
		SetGTPEncVerifyServerSignature(true).
		SetGTPEncVerifySignaturePublicKey(servPubKey).
		SetGTPCompression(gtp.Compression_Brotli).
		SetGTPCompressedSize(0).
		SetGTPAutoReconnectRetryTimes(0).
		SetZapLogger(logger).
		SetMainProcedure(proc).
		Connect(context.Background(), viper.GetString("endpoint"))
	if err != nil {
		panic(err)
	}

	go proc.Console()

	<-rpcli.Done()

	if err := context.Cause(rpcli); err != nil {
		rpcli.GetLogger().Infof("close cause:%s", err)
	}
}

const (
	gap = "\n\n"
)

type MainProc struct {
	rpcli.Procedure
	viewport viewport.Model
	textarea textarea.Model
	messages []string
}

func (m *MainProc) Init() tea.Cmd {
	return textarea.Blink
}

func (m *MainProc) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)
		m.viewport.GotoBottom()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			m.GetCli().Close(fmt.Errorf("console: %s", msg.Type))
			return m, tea.Quit
		case tea.KeyEnter:
			if err := rpc.ResultVoid(<-m.GetCli().RPC(misc.Chat, "ChatUserComp", "C_InputText", m.textarea.Value())).Extract(); err != nil {
				m.GetCli().GetLogger().Errorf("console: input %s failed, %s", m.textarea.Value(), err)
				break
			}
			m.textarea.Reset()
		}

	case error:
		m.GetCli().GetLogger().Errorf("console: %s", msg)
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m *MainProc) View() string {
	return fmt.Sprintf(
		"%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
	)
}

func (m *MainProc) Console() {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(150)
	ta.SetHeight(3)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(150, 5)
	vp.SetContent(`Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	m.textarea = ta
	m.viewport = vp

	if _, err := tea.NewProgram(m, tea.WithContext(m.GetCli())).Run(); err != nil {
		m.GetCli().GetLogger().Errorf("console: %s", err)
		return
	}
}

func (m *MainProc) OutputText(ts int64, channel, userId, text string) {
	strToColor := func(str string) colorful.Color {
		hash := fnv.New32()
		hash.Write([]byte(str))
		hv := hash.Sum32()

		hue := float64(hv%360) / 360.0
		saturation := 0.7 + 0.2*float64((hv/360)%3)
		value := 0.8 + 0.1*float64((hv/1080)%2)

		return colorful.Hsv(hue*360, saturation, value)
	}

	channelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(strToColor(channel).Hex()))
	userStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(strToColor(userId).Hex()))

	var you string

	if uid.Id(userId) == m.GetCli().GetSessionId() {
		you = "(YOU)"
	}

	msg := fmt.Sprintf("[%s][%s]%s: %s",
		time.Unix(ts, 0).Format(time.TimeOnly),
		channelStyle.Render(channel),
		userStyle.Render(userId+you),
		text,
	)

	if len(m.messages) > 100 {
		m.messages = m.messages[1:]
	}
	m.messages = append(m.messages, msg)

	m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
	m.viewport.GotoBottom()
}
