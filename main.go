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
	"errors"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strings"
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

type Object struct {
	name       string
	desc       string
	adjectives []string
	aliases    []string
	verbs      map[string]func(*Object, *Player)
	openable   bool
	open       bool
	// whether or not this object should be mentioned below the room description
	fixture bool
	// if this object can be picked up by the player
	carryable bool
}

// Return true if the string matches the object.
func (o *Object) RespondTo(args []string) bool {
	str := strings.Join(args, " ")
	// match all adjectives
	for _, adj := range o.adjectives {
		if i := strings.Index(str, strings.ToUpper(adj)+" "); i == 0 {
			str = str[len(adj)+1:]
		}
	}
	if str == strings.ToUpper(o.name) {
		return true
	}
	for _, alias := range o.aliases {
		if str == strings.ToUpper(alias) {
			return true
		}
	}
	return false
}

func (o *Object) GetName() string {
	res := "a " + o.name
	if o.openable {
		if o.open {
			res += " (open)"
		} else {
			res += " (closed)"
		}
	}
	return res
}

func (o *Object) GetDesc() string {
	res := ""
	if len(o.desc) > 0 {
		res = o.desc
		if o.openable {
			res += "\n"
		}
	}
	if o.openable {
		if o.open {
			res += fmt.Sprintf("The %v is open.", o.name)
		} else {
			res += fmt.Sprintf("The %v is closed.", o.name)
		}
	}
	return res
}

type ObjectContainer struct {
	objects []*Object
}

func (c *ObjectContainer) AddObject(objs ...*Object) {
	c.objects = append(c.objects, objs...)
}

func (c *ObjectContainer) RemoveObject(obj *Object) {
	loc := -1
	for i, val := range c.objects {
		if val == obj {
			loc = i
			break
		}
	}
	c.objects = append(c.objects[:loc], c.objects[loc+1:]...)
}

func (c *ObjectContainer) ObjectNames() (res string, err error) {
	count := 0
	res = ""
	for i, obj := range c.objects {
		if obj.fixture {
			continue
		}
		res += obj.GetName()
		count++
		if i < len(c.objects)-2 {
			res += ", "
		} else if i == len(c.objects)-2 {
			res += " and "
		}
	}
	if count == 0 {
		err = errors.New("no objects found")
	}
	return
}

func (c *ObjectContainer) FindObject(args []string) *Object {
	for _, obj := range c.objects {
		if obj.RespondTo(args) {
			return obj
		}
	}
	return nil
}

type Room struct {
	name    string
	desc    string
	visited bool
	// called when a player enters this room
	enterFunc func(*Player)
	// a function that can block exits if it returns false for the dir
	exitFunc func(dir string) bool
	// room exits:
	n, s, w, e, up, down, in, out *Room
	// a Room is a ObjectContainer (objects laying on the floor, or fixtures)
	ObjectContainer
}

func (r *Room) ExitDirection(dir string) *Room {
	if r.exitFunc != nil && !r.exitFunc(dir) {
		return nil
	}
	switch dir {
	case "NORTH":
		return r.n
	case "SOUTH":
		return r.s
	case "WEST":
		return r.w
	case "EAST":
		return r.e
	case "UP":
		return r.up
	case "DOWN":
		return r.down
	case "IN":
		return r.in
	case "OUT":
		return r.out
	default:
		return nil
	}
}

func (r *Room) Enter() {
	r.visited = true
}

func (r *Room) Leave() {
}

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

// create a game world "instance" returning the starting room
func NewGameWorld() (*Room, *Room) {
	const nothingSpecialDesc = "You don't see anything special about this."
	// objects are items, furniture etc.
	window := Object{
		name:       "Window",
		desc:       "A small window, it is too dirty to look inside the house.",
		fixture:    true,
		openable:   true,
		adjectives: []string{"small"},
	}
	can := Object{
		name:      "Can",
		carryable: true,
		desc:      "This is a unlabled can.",
	}
	trapdoor := Object{
		name:      "Trapdoor",
		carryable: false,
		openable:  true,
		aliases:   []string{"door"},
	}
	rug := Object{
		name:    "Rug",
		desc:    "A large oriental rug is covering the floor, it looks very dusty and pale.",
		fixture: true,
		// its not openable, but we use the open flag to indicate if it was shoved aside
		openable:   false,
		adjectives: []string{"large", "huge", "oriental", "dusty", "pale"},
		aliases:    []string{"floor"},
		verbs: map[string]func(*Object, *Player){
			"PUSH": func(object *Object, player *Player) {
				player.Println("Pushing the rug won't do anything, instead you should try to pull it.")
			},
			"PULL": func(object *Object, player *Player) {
				if !object.open {
					player.Println("Pulling the rug aside, revealed a trapdoor.")
					object.desc = "A large oriental rug lies rolled up on the floor, there is a trapdoor the rug was covering."
					player.room.desc = "Even in the day the room is sparsly lit. A large rug lies rolled up on the floor. The front door is boarded shut."
					player.room.AddObject(&trapdoor)
					object.open = true
					player.GivePoints(1)
				} else {
					player.Println("Pulling the rug further won't accomplish anything.")
				}
			},
			"LOOK UNDER": func(object *Object, player *Player) {
				player.Println("Be more specific, how do you look under a rug exactly?")
			},
		},
	}
	fish := Object{
		name:       "Trout",
		desc:       "The smell of this rotten fish gives you a headache.",
		adjectives: []string{"large", "smelly", "rotten"},
		aliases:    []string{"fish"},
	}
	bed := Object{
		name:    "Bed",
		fixture: true,
		desc:    "You can't find anything interesting in the bed, but the smell gets worse near it.",
		verbs: map[string]func(*Object, *Player){
			"LOOK UNDER": func(object *Object, player *Player) {
				if !object.open {
					player.Println("Under the bed is a large smelly trout.\nTaken.")
					// make the smell disappear:
					object.desc = "You can't find anything interesting in the bed."
					player.room.desc = "There is only a bed and a wooden cabinet in this plain bedroom."
					player.AddObject(&fish)
					object.open = true
					player.GivePoints(3)
				} else {
					player.Println("There is nothing under the bed.")
				}
			},
		},
	}
	cabinet := Object{
		name:    "Cabinet",
		fixture: true,
		desc:    nothingSpecialDesc,
	}
	/*                           +----------------+
	                             |                |
	        +--------------------+ North of House +-------------+
	        |                    |                |             |
	        |                    +----------------+             |
	        v                                                   |
	        +                                                   |
	        |                                                  \|/
	+-------+-------+      +------------+  +---------+  +-------+------+
	|               |  +   |            |  |         |  |              |
	| West of House +--+   | Livingroom +--+ Kitchen +--+ Behind House |
	|               |  +   |            |  |         |  |              |
	+-------+-------+      +------------+  +---------+  +-------+------+
	        |                                                  /|\
	        +                                                   |
	        ^                                                   |
	        |                   +----------------+              |
	        |                   |                |              |
	        +-------------------+ South of House +--------------+
	                            |                |
	                            +----------------+*/
	// rooms:
	nhouse := Room{name: "North of House", desc: "The path leads around the house to the east."}
	whouse := Room{name: "West of House", desc: "You are standing in an open field west of a white house with a boarded front door. Pathways lead north and south around the house."}
	shouse := Room{name: "South of House", desc: "The pathway extends to the east behind the white house."}
	bhouse := Room{name: "Behind House", desc: "To your west is a white house with a small window. Pathways lead north and south around the house.",
		exitFunc: func(dir string) bool {
			if dir == "WEST" || dir == "IN" {
				return window.open
			}
			return true
		}}
	kitchen := Room{name: "Kitchen", desc: "The kitchen looks as if it weren't used for many years. The room opens to the west into the livingroom and a staircase leads upwards.",
		exitFunc: func(dir string) bool {
			if dir == "EAST" || dir == "OUT" {
				return window.open
			}
			return true
		}}
	lroom := Room{name: "Living Room", desc: "Even in the day the room is sparsly lit. A huge rug is covering the floor. The front door is boarded shut.",
		exitFunc: func(dir string) bool {
			if dir == "DOWN" {
				return trapdoor.open
			}
			return true
		}}
	bedroom := Room{name: "Bedroom", desc: "There is only a bed and a wooden cabinet in this plain bedroom. Something smells terrible, giving you a light headache."}
	passage := Room{name: "Passage", desc: "You are standing in a damp narrow passageway. There is a ladder leading up, the passage continues to the north."}
	troom := Room{name: "Troll Room", desc: "Light shines from the high ceiling to the remains of unlucky, half-eaten adventurers."}
	// connections:
	nhouse.w, nhouse.e = &whouse, &bhouse
	whouse.n, whouse.s = &nhouse, &shouse
	shouse.w, shouse.e = &whouse, &bhouse
	bhouse.n, bhouse.s, bhouse.w = &nhouse, &shouse, &kitchen
	kitchen.e, kitchen.w = &bhouse, &lroom
	lroom.e = &kitchen
	bhouse.in = &kitchen
	kitchen.up, kitchen.out = &bedroom, &bhouse
	bedroom.down = &kitchen
	lroom.down = &passage
	passage.n, passage.up = &troom, &lroom
	troom.s = &passage
	// place objects in rooms:
	bhouse.AddObject(&window)
	kitchen.AddObject(&window)
	kitchen.AddObject(&can)
	lroom.AddObject(&rug)
	bedroom.AddObject(&bed, &cabinet)

	// return starting room and troll room:
	return &whouse, &troom
}

// how many turns before the troll will kill the player
const trollDifficulty = 5

/**
 * The troll AI is very simple, it will follow the player around the house and
 * kill him, if the player drops the fish the troll will die on food poisoning
 * and the player will win.
 */
type TrollAI struct {
	troll  Object
	player *Player
	// the room the troll is currently in:
	room *Room
	// indicates if player movements should "teleport" the troll too
	follow bool
	// turns until the troll will kill the player
	aggro uint
}

func (ai *TrollAI) Turn() {
	// win the game by giving him the fish
	if fish := ai.room.FindObject([]string{"TROUT"}); fish != nil {
		ai.player.Println("The troll sees the fish on the floor, immediately picks it up and eats it without chewing in a single gulp.")
		ai.player.Println("The troll looks ill, slowly, the huge creature sinks onto the floor.")
		ai.player.Println("The rotten fish killed the troll, by giving him food poisoning!")
		ai.player.GivePoints(5)
		ai.player.Win()
		return
	}
	// the can is not the solution but close:
	if can := ai.room.FindObject([]string{"CAN"}); can != nil {
		ai.player.Println("The troll sees the can on the floor, immediately picks it up and eats it without chewing in a single gulp.")
		ai.player.Println("Still the beast looks hungry at you.")
		ai.room.RemoveObject(can)
		ai.player.GivePoints(2)
		// this also resets the aggro counter
		if ai.aggro < trollDifficulty/2 {
			ai.aggro = trollDifficulty / 2
		}
		return
	}
	if ai.room == ai.player.room {
		if ai.aggro <= 0 {
			// there is a chance you survive:
			if rand.Intn(100) > 70 {
				ai.player.Println("The troll strikes at you with his club, hitting you on the head.")
				ai.player.Die()
			} else {
				ai.player.Println("The troll leaps forward and tries to hit you with his club, but misses!")
				// this resets the aggro counter
				ai.aggro = trollDifficulty / 2
			}
		} else if ai.aggro == 3 {
			ai.player.Println("The troll looks at you threateningly.")
		}
		ai.aggro--
	}
}

func (ai *TrollAI) PlayerMove() {
	if ai.follow {
		// reset kill turn counter
		ai.aggro = trollDifficulty
		// do not follow the player outside (trolls cant fit through windows?
		// also they would turn to stone in the sunlight!):
		if ai.player.room.name == "Behind House" {
			ai.follow = false
			return
		}

		// move troll to where the player is:
		ai.room.RemoveObject(&ai.troll)
		ai.player.room.AddObject(&ai.troll)
		ai.room = ai.player.room
		ai.player.Println("The monstrous creature follows you into the room!")
	} else if ai.room == ai.player.room {
		// player moved into the room where the troll is, follow him!
		ai.follow = true
	}
}

func (ai *TrollAI) Init(room *Room, player *Player) {
	ai.aggro = trollDifficulty
	ai.troll = Object{
		name:       "Troll",
		carryable:  false,
		desc:       "A huge, dangerous creature with sharp fanged teeth and a big broad nose. The monster is holding a heavy looking club in one of its enourmous hands.",
		adjectives: []string{"huge", "dangerous"},
		aliases:    []string{"creature", "monster"},
	}
	room.AddObject(&ai.troll)
	ai.room = room
	ai.player = player
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
