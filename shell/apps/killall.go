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

func CreateKillAll(t commandcreator.ICreator) *cli.Command {
	root := t.CreateCommand()
	root.Use = "killall"
	root.Short = "Kill All"
	root.Long = "Kill All"
	root.Run = func(cmd *cli.Command, pid int, args []string) {
		r := cmd.GetRootContext()
		r.WriteLn("")
		var arg string
		if len(args) > 0 {
			arg = args[0]
		}

		count := r.DeactivateAll(arg)

		r.WriteLn("Task deactivated: " + strconv.Itoa(count))
	}
	return root
}