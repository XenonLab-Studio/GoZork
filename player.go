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
	"fmt"
	"sort"
	"strings"
)

type Player struct {
	console   *bufio.ReadWriter
	room      *Room
	maxPoints byte
	points    byte
	trollai   TrollAI
	dead      bool
	win       bool
	// a player is a object container (inventory)
	ObjectContainer
}

// verbs are mapped to Player methods:

func (p *Player) Go(args []string) bool {
	if newRoom := p.room.ExitDirection(args[0]); newRoom != nil {
		p.room.Leave()
		p.room = newRoom
		p.Look(!newRoom.visited)
		newRoom.Enter()
		p.trollai.PlayerMove()
	} else {
		p.Println("You can't go in that direction.")
	}
	return true
}

func (p *Player) Look(printDesc bool) bool {
	//p.room.Look(true)
	p.Println(p.room.name)
	if printDesc && p.room.desc != "" {
		p.Println(p.room.desc)
	}
	if objstr, err := p.room.ObjectNames(); err == nil {
		p.Println("There is " + objstr + " here.")
	}
	return true
}
func (p *Player) LookAt(args []string) bool {
	if len(args) == 0 {
		return p.Look(true)
	}
	if obj := p.FindNearObject(args); obj != nil {
		p.Println(obj.GetDesc())
	} else {
		p.Printf("I don't see any %v here.\n", strings.Join(args, " "))
	}
	return true
}

func (p *Player) Take(args []string) bool {
	if obj := p.room.FindObject(args); obj != nil {
		if obj.carryable && !obj.fixture {
			p.room.RemoveObject(obj)
			p.AddObject(obj)
			p.Println("Taken.")
		} else {
			p.Println("This can't be taken.")
		}

	} else {
		p.Printf("I don't see any %v here.\n", strings.Join(args, " "))
	}
	return true
}

func (p *Player) Push(obj *Object) bool {
	if obj != nil {
		if cb := obj.verbs["PUSH"]; cb != nil {
			cb(obj, p)
		} else {
			p.Println("You can't push this.")
		}
	} else {
		p.Println("I don't know what you are referring to.")
	}
	return true
}

func (p *Player) Pull(obj *Object) bool {
	if obj != nil {
		if cb := obj.verbs["PULL"]; cb != nil {
			cb(obj, p)
		} else {
			p.Println("You can't pull this.")
		}
	} else {
		p.Println("I don't know what you are referring to.")
	}
	return true
}

func (p *Player) LookUnder(obj *Object) bool {
	if obj != nil {
		if cb := obj.verbs["LOOK UNDER"]; cb != nil {
			cb(obj, p)
		} else {
			p.Println("You don't see anything out of the ordinary.")
		}
	} else {
		p.Println("I don't know what you are referring to.")
	}
	return true
}

func (p *Player) Drop(args []string) bool {
	if obj := p.FindObject(args); obj != nil {
		p.RemoveObject(obj)
		p.room.AddObject(obj)
		p.Println("Dropped.")
	} else {
		p.Printf("I don't see any %v here.\n", strings.Join(args, " "))
	}
	return true
}

func (p *Player) Inventory(args []string) bool {
	if objstr, err := p.ObjectNames(); err == nil {
		p.Println("You are carrying " + objstr + ".")
	} else {
		p.Println("You are empty handed.")
	}
	return true
}

func (p *Player) Open(args []string) bool {
	if obj := p.FindNearObject(args); obj != nil {
		if obj.openable {
			if !obj.open {
				obj.open = true
				p.Println("Opened.")
			} else {
				p.Println("Already open.")
			}
		} else {
			p.Println("I can't open that.")
		}
	} else {
		p.Printf("I don't see any %v here.\n", strings.Join(args, " "))
	}
	return true
}

func (p *Player) Close(args []string) bool {
	if obj := p.FindNearObject(args); obj != nil {
		if obj.openable {
			if obj.open {
				obj.open = false
				p.Println("Closed.")
			} else {
				p.Println("Already closed.")
			}
		} else {
			p.Println("I can't close that.")
		}
	} else {
		p.Printf("I don't see any %v here.\n", strings.Join(args, " "))
	}
	return true
}

func (p *Player) Wait() bool {
	p.Println("Time passes.")
	return true
}

func (p *Player) Help(args []string) bool {
	p.Println("This is a text adventure game, the goal is to find and kill the troll.")
	p.Println("The game only understands very simple single-verb, single-object sentences, for instance: PICK UP HAT, or OPEN DOOR etc.")
	p.Println("The Verbs this game understands are: LOOK, LOOK AT, LOOK UNDER, PUSH, PULL, TAKE, DROP, WAIT, OPEN, CLOSE and INVENTORY.")
	p.Println("Directions are: NORTH, SOUTH, EAST, WEST, UP, DOWN, IN and OUT.")
	p.Println("There are also many aliases for verbs and directions.")
	return true
}

func (p *Player) GivePoints(off byte) {
	p.points += off
	p.Printf("(Your score increased by %d points, you now have %d/%d points.)\n", off, p.points, p.maxPoints)
}

func (p *Player) Die() {
	p.dead = true
	p.Println(" **** GAME OVER! You are dead.")
	p.Printf("You managed to score %d out of %d possible points.\n", p.points, p.maxPoints)
}

func (p *Player) Win() {
	p.win = true
	p.Println(" **** CONGRATULATIONS! YOU WON THE GAME!")
	p.Printf("You managed to score %d out of %d possible points.\n", p.points, p.maxPoints)
}

func (p *Player) FindNearObject(args []string) *Object {
	// objects in the current room:
	if obj := p.room.FindObject(args); obj != nil {
		return obj
	}
	// objects/items in the player inventory:
	if obj := p.FindObject(args); obj != nil {
		return obj
	}
	return nil
}

func (p *Player) ExecuteCommand(command string) bool {
	verbMap := map[string]func([]string) bool{
		"GO":         func(args []string) bool { return p.Go(args) },
		"LOOK":       func(args []string) bool { return p.Look(true) },
		"LOOK AT":    func(args []string) bool { return p.LookAt(args) },
		"TAKE":       func(args []string) bool { return p.Take(args) },
		"PUSH":       func(args []string) bool { return p.Push(p.FindNearObject(args)) },
		"PULL":       func(args []string) bool { return p.Pull(p.FindNearObject(args)) },
		"LOOK UNDER": func(args []string) bool { return p.LookUnder(p.room.FindObject(args)) },
		"DROP":       func(args []string) bool { return p.Drop(args) },
		"OPEN":       func(args []string) bool { return p.Open(args) },
		"WAIT":       func(args []string) bool { return p.Wait() },
		"CLOSE":      func(args []string) bool { return p.Close(args) },
		"INVENTORY":  func(args []string) bool { return p.Inventory(args) },
		"XYZZY":      func(args []string) bool { return true },
		"HELP":       func(args []string) bool { return p.Help(args) },
	}
	delegated := false
	// we need to make sure to sort the verbs by length first:
	verbs := []string{} // make([]string, len(verbMap))
	for verb := range verbMap {
		verbs = append(verbs, verb)
	}
	sort.Sort(ByLength(verbs))
	sort.Sort(sort.Reverse(sort.StringSlice(verbs)))
	for _, verb := range verbs {
		fn := verbMap[verb]
		if verb == command {
			delegated = fn(make([]string, 0))
		} else if i := strings.Index(command, verb+" "); i == 0 {
			command = command[len(verb)+1:]
			delegated = fn(strings.Split(command, " "))
		}
	}
	if !delegated {
		p.Println("Sorry, what?")
	} else {
		p.trollai.Turn()
	}
	return delegated
}

func (p *Player) VerbAliasReplace(cmd string) string {
	for verb, aliases := range verbAliasMap {
		for _, alias := range aliases {
			if alias == cmd {
				return verb
			} else if i := strings.Index(cmd, verb+" "); i == 0 {
				return cmd
			} else if i := strings.Index(cmd, alias+" "); i == 0 {
				return strings.Replace(cmd, alias+" ", verb+" ", 1)
			}
		}
	}
	return cmd
}

func (p *Player) Println(line string) {
	p.Printf(line + "\n")
}

func (p *Player) Printf(format string, args ...interface{}) {
	line := fmt.Sprintf(format, args...)
	p.console.WriteString(line)
	p.console.Flush()
}

func (p *Player) Run() {
	p.Println("Welcome to GOZORK! Type HELP for help.")

	// start the player west of the house
	p.room.Enter()
	p.Look(true)

	for {
		p.Printf(">")

		cmd, err := p.console.ReadString('\n')
		if err != nil {
			p.Println("error reading console")
			break
		}
		cmd = strings.ToUpper(strings.Trim(cmd, "\n"))

		// replace alias mapping
		cmd = p.VerbAliasReplace(cmd)
		//fmt.Printf("[command read as: %v]\n", cmd)

		if cmd == "QUIT" {
			p.Println("Thanks for playing!")
			break
		}

		p.ExecuteCommand(cmd)

		if p.dead || p.win {
			break
		}
	}
}
