/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package telnet

import (
	"fmt"
	"github.com/markel1974/goshell/shell/adaptiveticker"
	"github.com/markel1974/goshell/shell/cli"
	"github.com/markel1974/goshell/shell/context"
	"github.com/markel1974/goshell/shell/interfaces"
	"github.com/markel1974/goshell/shell/telnet/session"
	"github.com/markel1974/goshell/shell/terminal"
	"log"
	"net"
)

type Server struct {
	ticker   *adaptiveticker.AdaptiveTicker
	template *cli.Command
	prompt   string
	addr     string
	factory  *terminal.EquipmentFactory
	auth     interfaces.IAuthenticator
	autosave bool
}

func NewServer(ticker *adaptiveticker.AdaptiveTicker, auth interfaces.IAuthenticator, port int, autosave bool) *Server {
	return &Server{
		ticker:   ticker,
		addr:     fmt.Sprintf(":%d", port),
		factory:  terminal.NewEquipmentFactory(),
		auth:     auth,
		autosave: autosave,
	}
}

func (r *Server) SetPrompt(prompt string) {
	r.prompt = prompt
}

func (r *Server) handleConnection(c net.Conn) {
	//fmt.Printf("Serving %s\n", c.RemoteAddr().String())

	defer func() {
		if r := recover(); nil != r {
			log.Printf("Recovered from: (%T) %v\n"+"", r, r)
		}
	}()

	telnetSession := session.NewTelnet(c)

	ctx := context.NewContext(r.ticker, telnetSession, telnetSession, r.auth, r.factory, r.template, r.prompt, r.autosave)

	ctx.Setup()

	telnetSession.SetListenFunc(func(code session.IOCode, data []byte) {
		switch code {
		case session.WS:
			if len(data) != 4 {
				log.Println("Malformed window size data:", data)
				return
			}

			width := int(255*data[0]) + int(data[1])
			height := int(255*data[2]) + int(data[3])
			ctx.SetScreenSize(width, height)

		case session.TT:
			//c.terminal.SetTerminalType(string(data))

		default:
			log.Println("Unknown code", code, "data", data)
		}
	})
	telnetSession.WillEcho()
	telnetSession.WillSga()
	telnetSession.DoWindowSize()
	telnetSession.DoTerminalType()

	ctx.Exec()
	ctx.Close()

	_ = c.Close()
}

func (r *Server) SetTemplate(template *cli.Command) {
	r.template = template
}

func (r *Server) Start() {
	l, err := net.Listen("tcp4", r.addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go r.handleConnection(c)
	}
}

func (r *Server) AsyncStart() {
	go func() {
		r.Start()
	}()
}
