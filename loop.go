package gogame

import "time"

type Runnable interface {
	IsFinished() bool
	Events()
	Update()
	Render()
}

// Update the game using a fixed time step (makes everything simpler).
// Pass a runnable object (implement methods IsFinished(), Draw(), Update() and Events()).
//
// See this link: http://gameprogrammingpatterns.com/game-loop.html
//
// It goes like this: A certain amount of time has elapsed since the last
// return of the game loop. This is how much time we need to simulate for
// the game to catch up. We do that using a series of fixed time steps.
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
