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
	"os"
	"os/exec"
	"runtime"
)

// Create a map for storing clear funcs
var clear map[string]func()

func init() {
	// Initialize it
	clear = make(map[string]func())
	clear["linux"] = func() {
		// Linux example
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		// Windows
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	// runtime.GOOS -> linux, windows, darwin etc.
	value, ok := clear[runtime.GOOS]
	// if we defined a clear func for that platform:
	if ok {
		// we execute it
		value()
	} else {
		// unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
