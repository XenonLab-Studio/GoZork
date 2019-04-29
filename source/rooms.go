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

package source

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
