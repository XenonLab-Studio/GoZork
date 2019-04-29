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
