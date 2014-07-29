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

func KillFromAllGroups(s Sprite) {
	if gl, ok := infoSprites[s]; ok {
		for _, g := range gl {
			g.Remove(s)
		}
		delete(infoSprites, s)
	}

}

// Check if two sprites are colliding
func Collide(c1, c2 Sprite) bool {
	r1 := c1.GetRect()
	r2 := c2.GetRect()
	left1, left2 := r1.X, r2.X
	right1, right2 := r1.X+r1.W, r2.X+r2.W
	top1, top2 := r1.Y, r2.Y
	bottom1, bottom2 := r1.Y+r1.H, r2.Y+r2.H

	if left1 <= right2 && left2 <= right1 && top1 <= bottom2 && top2 <= bottom1 {
		// TODO: check for pixel perfect collision. Check alpha components of the two textures
		return true
	}

	return false
}

type Sprite interface {
	GetRect() *gogame.Rect
	Update()
	Draw()
}

type Group struct {
	sprList *list.List
}

func NewGroup() *Group {
	g := new(Group)
	g.sprList = list.New()
	return g
}

func (self *Group) Add(s Sprite) {
	self.sprList.PushFront(s)
	registerSprite(s, self)
}

func (self *Group) Len() int {
	return self.sprList.Len()
}

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

func (self *Group) Update() {
	e := self.sprList.Front()
	for e != nil {
		s := e.Value.(Sprite)
		e = e.Next()
		s.Update()
	}
}

func (self *Group) Draw() {
	for e := self.sprList.Front(); e != nil; e = e.Next() {
		s := e.Value.(Sprite)
		s.Draw()
	}
}

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
// returns the sprites colliding (group1, group2)
func (self *Group) CollideGroup(g *Group) (Sprite, Sprite, bool) {
	for e := g.sprList.Front(); e != nil; e = e.Next() {
		s := e.Value.(Sprite)
		if ss, ok := self.CollideSpr(s); ok {
			return ss, s, true
		}
	}
	return nil, nil, false
}
