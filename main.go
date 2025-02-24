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

package main

import (
	"fmt"
	"github.com/markel1974/goshell/shell"
	"github.com/markel1974/goshell/shell/authenticator"
	"github.com/markel1974/goshell/shell/cli"
	"github.com/markel1974/goshell/shell/interfaces"
	"log"
	"math/rand"
)

func fooCommand() *cli.Command {
	foo := cli.NewCommand()
	foo.Run = func(cmd *cli.Command, pid int, args []string) {
		cmd.WriteLn([]byte{})
		cmd.WriteLn([]byte("Foo"))
	}
	foo.Use = "foo"
	foo.Short = "Command foo"
	foo.Long = "This is a command"

	bar := cli.NewCommand()
	bar.Activate = true
	bar.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn("")
		r.Write(fmt.Sprint("Bar is opened, pid ", pid))
	}
	bar.ReadEvent = func(cmd *cli.Command, pid int, ctx interface{}, code int, key rune) {
	}
	bar.Use = "bar"
	bar.Short = "Command bar"
	bar.Long = "This is a nested command"

	var data []float64

	plot := cli.NewCommand()
	plot.Use = "plot"
	plot.Short = "Plot"
	plot.Long = "Plot"
	plot.Activate = true
	plot.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.CreateTimer(pid, 0, 300, -1)
	}
	plot.TimerEvent = func(cmd *cli.Command, pid int, tid int, ctx interface{}, interval int) {
		cmd.GetRootContext().PaintRequest(pid)
	}
	plot.PaintEvent = func(cmd *cli.Command, pid int, ctx interface{}, surface interfaces.ISurface) {
		max := 100
		min := 0
		cpuUsage := float64(rand.Intn(max-min) + min)
		data = append(data, cpuUsage)
		if len(data) > 10 {
			data = data[1:]
		}
		surface.DrawSeries(data, -1, -1, -1, -1)
	}

	_ = foo.AddCommand(bar)
	_ = foo.AddCommand(plot)

	return foo
}

func main() {
	//reader := bufio.NewReader(os.Stdin)
	//writer := bufio.NewWriter(os.Stdout)

	/*
		max := 100
		min := 0
		var data[]float64
		for ;; {
			cpuUsage := float64(rand.Intn(max-min) + min)
			data = append(data, cpuUsage)
			if len(data) > 10 {
				data = data[1:]
			}

			out := graph.Plot(data, "\n", -1, -1, graph.Width(30), graph.Height(10))
			fmt.Println(string(out))
			fmt.Println("")
			time.Sleep(1 * time.Second)
		}
	*/

	const appName = "foo"
	const appVersion = "1.1.1"
	const prompt = appName + " " + appVersion + "> "
	const port = 1234
	const user = "u"
	const secure = true
	const pass = "p"

	t := cli.NewCommand()
	_ = t.AddCommand(fooCommand())

	auth := authenticator.NewSimpleAuthenticator()
	if err := auth.Setup(user, pass); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting shell")
	fmt.Println("port", port)
	fmt.Println("secure", secure)
	fmt.Println("user", user)
	k := shell.New(secure, auth, port, false)
	k.SetPrompt(prompt)
	k.SetTemplate(t)

	k.Start()
}
