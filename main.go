/*
 *****************************************************
 * © 2019 Stefano Peris <xenonlab.develop@gmail.com> *
 *****************************************************
 *
 * Released under the GNU/GPL 3.0 license
 *
 * github: <https://github.com/XenonLab-Studio/GoZork>
 *
 *
 * :'######::::'#######::'########::'#######::'########::'##:::'##:
 *  ##... ##::'##.... ##:..... ##::'##.... ##: ##.... ##: ##::'##::
 *  ##:::..::: ##:::: ##::::: ##::: ##:::: ##: ##:::: ##: ##:'##:::
 *  ##::'####: ##:::: ##:::: ##:::: ##:::: ##: ########:: #####::::
 *  ##::: ##:: ##:::: ##::: ##::::: ##:::: ##: ##.. ##::: ##. ##:::
 *  ##::: ##:: ##:::: ##:: ##:::::: ##:::: ##: ##::. ##:: ##:. ##::
 * . ######:::. #######:: ########:. #######:: ##:::. ##: ##::. ##:
 * :......:::::.......:::........:::.......:::..:::::..::..::::..::
 *
 *     Textual adventure written in golang inspired by "Zork I"
 */

package main

import (
	"bufio"
	"os"

	"github.com/XenonLab-Studio/GoZork/clearscr"
	"github.com/XenonLab-Studio/GoZork/objects"
	"github.com/XenonLab-Studio/GoZork/player"
	"github.com/XenonLab-Studio/GoZork/rooms"
	"github.com/XenonLab-Studio/GoZork/trollai"
)

// sort by string length:
// https://mmcgrana.github.io/2012/09/go-by-example-sort-by-function.html
type ByLength []string

func (s ByLength) Len() int {
	return len(s)
}
func (s ByLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

var verbAliasMap = map[string][]string{
	"GO NORTH":   {"N", "NORTH"},
	"GO SOUTH":   {"S", "SOUTH"},
	"GO WEST":    {"W", "WEST"},
	"GO EAST":    {"E", "EAST"},
	"GO IN":      {"IN", "INSIDE", "ENTER"},
	"GO OUT":     {"OUT", "OUTSIDE", "LEAVE"},
	"GO UP":      {"UP"},
	"GO DOWN":    {"DOWN"},
	"LOOK AT":    {"EXAMINE", "INSPECT", "X"},
	"LOOK UNDER": {"LOOK BENEATH", "LOOK BELOW"},
	"TAKE":       {"PICK UP", "GET"},
	"DROP":       {"THROW"},
	"INVENTORY":  {"I"},
	"WAIT":       {"Z"},
}

func main() {
	// everything goes through this buffered read/writer which
	// would make it very easy to plug this onto a telnet server.
	console := bufio.NewReadWriter(
		bufio.NewReader(os.Stdin),
		bufio.NewWriter(os.Stdout))
	player := Player{console: console, maxPoints: 11}
	start, trollroom := NewGameWorld()
	player.room = start
	player.trollai.Init(trollroom, &player)
	player.Run()
}
