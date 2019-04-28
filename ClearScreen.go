/*
 *****************************************************
 * © 2019 Stefano Peris <xenonlab.develop@gmail.com> *
 *****************************************************
 *
 * Released under the GNU/GPL 3.0 license
 *
 * github: <https://github.com/XenonLab-Studio/The-Great-Adventure>
 *
 * ▄▄▄▄▄ ▄ .▄▄▄▄ .
 * •██  ██▪▐█▀▄.▀·
 *  ▐█.▪██▀▐█▐▀▀▪▄
 *  ▐█▌·██▌▐▀▐█▄▄▌
 *  ▀▀▀ ▀▀▀ · ▀▀▀
 *  ▄▄ • ▄▄▄  ▄▄▄ . ▄▄▄· ▄▄▄▄▄
 * ▐█ ▀ ▪▀▄ █·▀▄.▀·▐█ ▀█ •██
 * ▄█ ▀█▄▐▀▀▄ ▐▀▀▪▄▄█▀▀█  ▐█.▪
 * ▐█▄▪▐█▐█•█▌▐█▄▄▌▐█ ▪▐▌ ▐█▌·
 * ·▀▀▀▀ .▀  ▀ ▀▀▀  ▀  ▀  ▀▀▀
 *  ▄▄▄· ·▄▄▄▄   ▌ ▐·▄▄▄ . ▐ ▄ ▄▄▄▄▄▄• ▄▌▄▄▄  ▄▄▄ .
 * ▐█ ▀█ ██▪ ██ ▪█·█▌▀▄.▀·•█▌▐█•██  █▪██▌▀▄ █·▀▄.▀·
 * ▄█▀▀█ ▐█· ▐█▌▐█▐█•▐▀▀▪▄▐█▐▐▌ ▐█.▪█▌▐█▌▐▀▀▄ ▐▀▀▪▄
 * ▐█ ▪▐▌██. ██  ███ ▐█▄▄▌██▐█▌ ▐█▌·▐█▄█▌▐█•█▌▐█▄▄▌
 *  ▀  ▀ ▀▀▀▀▀• . ▀   ▀▀▀ ▀▀ █▪ ▀▀▀  ▀▀▀ .▀  ▀ ▀▀▀
 * ._.                  .     ._  .___.      .    ._.
 *  | ._  __._ *._. _  _|   _ |,    _/  _ ._.;_/   |
 * _|_[ )_) [_)|[  (/,(_]  (_)|   ./__.(_)[  | \  _|_
 *          |
 */

package main

import (
	"os"
	"os/exec"
	"runtime"
)

// clear screen
var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
