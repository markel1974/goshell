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

package games

import (
	"github.com/markel1974/goshell/shell/apps/commandcreator"
	"github.com/markel1974/goshell/shell/apps/games/invaders"
	"github.com/markel1974/goshell/shell/cli"
	"github.com/markel1974/goshell/shell/interfaces"
)

func CreateInvaders(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "invaders"
	root.Short = "Invaders"
	root.Long = "Invaders"
	root.Activate = true
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		w, h := r.GetScreenSize()
		g := invaders.NewGame(w, h)
		g.SetMenuState()
		r.SetContext(pid, g)
		r.CreateTimer(pid, 0, 100, -1)
	}
	root.ReadEvent = func(cmd *cli.Command, pid int, ctx interface{}, code int, key rune) {
		g := ctx.(*invaders.Invaders)
		g.HandleKey(key)
	}
	root.TimerEvent = func(cmd *cli.Command, pid int, tid int, ctx interface{}, interval int) {
		r := cmd.GetRootContext()
		g := ctx.(*invaders.Invaders)
		g.Update()
		r.PaintRequest(pid)
	}
	root.PaintEvent = func(cmd *cli.Command, pid int, ctx interface{}, surface interfaces.ISurface) {
		g := ctx.(*invaders.Invaders)

		rows, columns := surface.GetSize()
		w, h := g.GetSize()
		if h != rows || w != columns {
			g.SetSize(columns, rows)
			g.SetMenuState()
		}

		g.Draw(surface)
	}

	return root
}
