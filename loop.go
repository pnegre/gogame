package gogame

import "time"

type Runnable interface {
	IsFinished() bool
	Events()
	Update()
	Render()
}

// See this: http://gameprogrammingpatterns.com/game-loop.html
func Loop(runnable Runnable, updatesPerSecond int) {
	var currentTime, elapsedTime, lag int64
	updateInterval := int64(time.Second) / int64(updatesPerSecond)
	previousTime := time.Now().UnixNano()

	runnable.Update()
	for !runnable.IsFinished() {
		currentTime = time.Now().UnixNano()
		elapsedTime = currentTime - previousTime
		previousTime = currentTime
		lag += elapsedTime

		runnable.Events()

		for lag >= updateInterval {
			runnable.Update()
			lag -= updateInterval
		}

		runnable.Render()
	}
}
