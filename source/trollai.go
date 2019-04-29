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

import (
	"math/rand"
)

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
