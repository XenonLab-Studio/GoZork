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

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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

package main

import (
	"os"
)

func main() {

}
