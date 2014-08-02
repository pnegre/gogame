package sprite

import "container/list"
import "github.com/pnegre/gogame"

type glist []*Group

var infoSprites = make(map[Sprite]glist)

func registerSprite(sp Sprite, g *Group) {
	if gl, ok := infoSprites[sp]; ok {
		gl = append(gl, g)
		infoSprites[sp] = gl
		return
	}

	var gl glist
	gl = append(gl, g)
	infoSprites[sp] = gl
}

// Remove the sprite from all groups
func KillFromAllGroups(s Sprite) {
	if gl, ok := infoSprites[s]; ok {
		for _, g := range gl {
			g.Remove(s)
		}
		delete(infoSprites, s)
	}

}

// Check if two sprites are colliding, using rects
func Collide(c1, c2 Sprite) bool {
	r1 := c1.GetRect()
	r2 := c2.GetRect()
	// TODO: check for pixel perfect collision. Check alpha components of the two textures
	return r1.Intersects(r2)
}

type Sprite interface {
	GetRect() *gogame.Rect
	Update()
	Draw()
}

// Container to hold and manage multiple sprite objects
type Group struct {
	sprList *list.List
}

// Creates new container
func NewGroup() *Group {
	g := new(Group)
	g.sprList = list.New()
	return g
}

// Add new sprite to group
func (self *Group) Add(s Sprite) {
	self.sprList.PushFront(s)
	registerSprite(s, self)
}

// Get number of sprites of this group
func (self *Group) Len() int {
	return self.sprList.Len()
}

// Get sprite number "n" from this group
func (self *Group) GetElement(n int) Sprite {
	e := self.sprList.Front()
	i := 0
	for e != nil {
		if i == n {
			return e.Value.(Sprite)
		}
		e = e.Next()
		i++
	}
	return nil
}

// Calls the Update() method on all sprites in this group
func (self *Group) Update() {
	e := self.sprList.Front()
	for e != nil {
		s := e.Value.(Sprite)
		e = e.Next()
		s.Update()
	}
}

// Clears group list
func (self *Group) Clear() {
	self.sprList = list.New()
}

// Calls the Draw() method on all sprites in this group
func (self *Group) Draw() {
	for e := self.sprList.Front(); e != nil; e = e.Next() {
		s := e.Value.(Sprite)
		s.Draw()
	}
}

// Remove sprite from group
func (self *Group) Remove(sp Sprite) {
	var le *list.Element
	for e := self.sprList.Front(); e != nil; e = e.Next() {
		s := e.Value.(Sprite)
		if s == sp {
			le = e
			break
		}
	}
	if le != nil {
		self.sprList.Remove(le)

	}

}

// Find the sprite in this group that collides with the provided sprite.
func (self *Group) CollideSpr(sp Sprite) (Sprite, bool) {
	for e := self.sprList.Front(); e != nil; e = e.Next() {
		s := e.Value.(Sprite)
		if Collide(s, sp) {
			return s, true
		}
	}
	return nil, false
}

// Check if any two of sprites of the two groups collide
// returns the colliding sprites (group1, group2)
func (self *Group) CollideGroup(g *Group) (Sprite, Sprite, bool) {
	for e := g.sprList.Front(); e != nil; e = e.Next() {
		s := e.Value.(Sprite)
		if ss, ok := self.CollideSpr(s); ok {
			return ss, s, true
		}
	}
	return nil, nil, false
}
