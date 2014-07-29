package gogame

import "time"
import "log"

type Clock struct {
	frameDuration int64
	lastTime      int64
	ticks         int64
	totalDuration int64
}

func NewClock(fps int) *Clock {
	frameDuration := int64(1e9 / fps)
	log.Println(frameDuration)
	return &Clock{frameDuration: frameDuration}
}

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
		fps := float64(1.0) / ((float64(self.totalDuration) / float64(1e9)) / float64(self.ticks))
		log.Printf("Clock. FPS: %.2f", fps)
		self.totalDuration = 0
		self.ticks = 0
	}
}
