package gogame

import "time"
import "log"

type Clock struct {
	frameDuration int64
	lastTime      int64
	ticks         int64
	totalDuration int64
	avgFPS        float64
}

// Constructs new Clock object. Can be used to track time and control a game's framerate.
func NewClock(fps int) *Clock {
	frameDuration := int64(1e9 / fps)
	log.Println(frameDuration)
	return &Clock{frameDuration: frameDuration}
}

// Must be called once per frame. Limit framerate to internal value (provided on the object creation)
func (self *Clock) Wait() {
	self.ticks++
	if self.lastTime == 0 {
		self.lastTime = time.Now().UnixNano()
		return
	}

	delta := int64(time.Now().UnixNano() - self.lastTime)
	if delta < self.frameDuration {
		time.Sleep(time.Duration(self.frameDuration - delta))
	}

	t := time.Now().UnixNano()
	self.totalDuration += t - self.lastTime
	self.lastTime = t

	if self.ticks > 500 {
		self.avgFPS = float64(1.0) / ((float64(self.totalDuration) / float64(1e9)) / float64(self.ticks))
		self.totalDuration = 0
		self.ticks = 0
	}
}

// Gets averaged Frames Per Second. Updated every 500 iterations
func (self *Clock) GetFPS() float64 {
	return self.avgFPS
}
