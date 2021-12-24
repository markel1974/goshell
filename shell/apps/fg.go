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

package apps

import (
	"github.com/markel1974/goshell/shell/apps/commandcreator"
	"github.com/markel1974/goshell/shell/cli"
	"strconv"
)

func CreateFg(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "fg"
	root.Short = "Foreground"
	root.Long = "Foreground"
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()

		r.WriteLn("")

		if len(args) <= 0 {
			r.WriteLn("Empty argument")
			return
		}
		pid, err := strconv.Atoi(args[0])
		if err != nil {
			r.WriteLn("Invalid argument: " + args[0])
			return
		}

		if !r.SetFg(pid) {
			r.WriteLn("Unknown task: " + args[0])
		}
	}

	return root
}
