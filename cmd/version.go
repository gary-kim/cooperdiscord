//    Copyright (C) 2020 Gary Kim <gary@garykim.dev>, All Rights Reserved
//
//    This program is free software: you can redistribute it and/or modify
//    it under the terms of the GNU Affero General Public License as published
//    by the Free Software Foundation, either version 3 of the License, or
//    (at your option) any later version.
//
//    This program is distributed in the hope that it will be useful,
//    but WITHOUT ANY WARRANTY; without even the implied warranty of
//    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
//    GNU Affero General Public License for more details.
//
//    You should have received a copy of the GNU Affero General Public License
//    along with this program.  If not, see <https://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "v0.0.1"

func init () {
	versionCmd := &cobra.Command{
		Use: "version",
		Short: "Print version then exit",
		Run: func(command *cobra.Command, args []string) {
			fmt.Printf("dcli - Version %s\n", Version)
		},
	}
	Root.AddCommand(versionCmd)
}