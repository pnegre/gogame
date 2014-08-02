package sprite

import "github.com/pnegre/gogame"

// Simple sprite. Intended to be used as a base type for game objects
type Simple struct {
	drw  gogame.Drawable
	Rect gogame.Rect
}

// Creates new Simple sprite
func NewSimple(drw gogame.Drawable) *Simple {
	w, h := drw.GetDimensions()
	return &Simple{drw, gogame.Rect{0, 0, w, h}}
}

// Draw to screen, using internal Rect
func (self *Simple) Draw() {
	self.drw.BlitRect(&self.Rect)
}

// Gets internal Rect
func (self *Simple) GetRect() *gogame.Rect {
	return &self.Rect
}

// Update function for sprite
func (self *Simple) Update() {

}

type frame struct {
	drw   gogame.Drawable
	ticks int
}

// Animation. Simple framework for animating 2D sprites
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

// Creates new Animation object
func NewAnimation() *Animation {
	anim := new(Animation)
	return anim
}

// Add new frame (can be a texture or subtexture). Specify number of ticks for this frame.
func (self *Animation) AddFrame(drw gogame.Drawable, ticks int) {
	self.frames = append(self.frames, frame{drw, ticks})
	self.updateRect()
}

// Set animation in continous mode or single.
func (self *Animation) SetRepeat(etn bool) {
	self.eternal = etn
}

// Set sequence for animation. If this function is not called, the animation
// will use a standard sequence
func (self *Animation) SetSequence(seq []int) {
	self.sequence = seq
	self.Reset()
}

// Reset animation
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

// Update animation. Advance frame if necessary
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

// Draw animation to screen (current frame).
func (self *Animation) Draw() {
	self.frames[self.activeframe].drw.BlitRect(&self.Rect)
}

// Get animation Rect, corresponding to current frame
func (self *Animation) GetRect() *gogame.Rect {
	return &self.Rect
}

// For animations not in continuous mode, call this function to check if they have finished
func (self *Animation) IsFinished() bool {
	return self.finished
}
