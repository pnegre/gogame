package sprite

import "github.com/pnegre/gogame"

type Simple struct {
	drw  gogame.Drawable
	Rect gogame.Rect
}

func NewSimple(drw gogame.Drawable) *Simple {
	w, h := drw.GetDimensions()
	return &Simple{drw, gogame.Rect{0, 0, w, h}}
}

func (self *Simple) Draw() {
	self.drw.BlitRect(&self.Rect)
}

func (self *Simple) GetRect() *gogame.Rect {
	return &self.Rect
}

func (self *Simple) Update() {

}

type frame struct {
	drw   gogame.Drawable
	ticks int
}

type Animation struct {
	frames      []frame
	activeframe int
	ticks       int
	Rect        gogame.Rect
	finished    bool
	eternal     bool
	sequence    []int
	idxseq      int
}

func NewAnimation() *Animation {
	anim := new(Animation)
	return anim
}

func (self *Animation) AddFrame(drw gogame.Drawable, ticks int) {
	self.frames = append(self.frames, frame{drw, ticks})
	self.updateRect()
}

func (self *Animation) SetRepeat(etn bool) {
	self.eternal = etn
}

func (self *Animation) SetSequence(seq []int) {
	self.sequence = seq
	self.Reset()
}

func (self *Animation) Reset() {
	self.ticks = 0
	self.finished = false
	if self.sequence != nil {
		self.idxseq = 0
		self.activeframe = self.sequence[0]
	} else {
		self.activeframe = 0
	}
	self.updateRect()
}

func (self *Animation) updateRect() {
	frame := self.frames[self.activeframe]
	w, h := frame.drw.GetDimensions()
	self.Rect.W, self.Rect.H = w, h
}

func (self *Animation) Update() {
	if self.finished {
		return
	}
	self.ticks += 1

	if self.sequence == nil {
		// Normal sequence
		if self.ticks > self.frames[self.activeframe].ticks {
			if self.activeframe < len(self.frames)-1 {
				self.activeframe++
				self.ticks = 0
				self.updateRect()
			} else {
				if !self.eternal {
					self.finished = true
				} else {
					self.Reset()
				}
			}
		}

	} else {
		// Custom sequence
		if self.ticks > self.frames[self.activeframe].ticks {
			if self.idxseq < len(self.sequence)-1 {
				self.idxseq++
				self.activeframe = self.sequence[self.idxseq]
				self.ticks = 0
				self.updateRect()
			} else {
				if !self.eternal {
					self.finished = true
				} else {
					self.Reset()
				}
			}
		}

	}
}

func (self *Animation) Draw() {
	self.frames[self.activeframe].drw.BlitRect(&self.Rect)
}

func (self *Animation) GetRect() *gogame.Rect {
	return &self.Rect
}

func (self *Animation) IsFinished() bool {
	return self.finished
}
