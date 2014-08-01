package gogame

/*
#cgo pkg-config: sdl2
#include "SDL.h"

int getEventType(SDL_Event *e) {
    return e->type;
}

int getKeyCode(SDL_Event *e) {
    return e->key.keysym.sym;
}

int isKeyRepeat(SDL_Event *e) {
    return e->key.repeat;
}

int isKeyPressed(int kcode) {
    const Uint8 *state = SDL_GetKeyboardState(NULL);
    int sc = SDL_GetScancodeFromKey(kcode);
    return state[sc];
}

*/
import "C"

const (
	K_LEFT   = C.SDLK_LEFT
	K_RIGHT  = C.SDLK_RIGHT
	K_SPACE  = C.SDLK_SPACE
	K_ESC    = C.SDLK_ESCAPE
	K_RETURN = C.SDLK_RETURN
	K_P      = C.SDLK_p
	K_I      = C.SDLK_i
)

type Event interface{}

type EventQuit interface{}

type EventUnknown interface{}

type EventKey struct {
	Code int
	Down bool
}

// Poll for pending envents. Return nil if there is no event available
func PollEvent() Event {
	var cev C.SDL_Event

	for {
		if 0 == C.SDL_PollEvent(&cev) {
			return nil
		}

		switch C.getEventType(&cev) {

		case C.SDL_QUIT:
			return new(EventQuit)

		case C.SDL_KEYDOWN:
			// Ignore repeat key events
			if C.isKeyRepeat(&cev) != 0 {
				break
			}
			kde := new(EventKey)
			kde.Code = int(C.getKeyCode(&cev))
			kde.Down = true
			return kde

		case C.SDL_KEYUP:
			kde := new(EventKey)
			kde.Code = int(C.getKeyCode(&cev))
			kde.Down = false
			return kde

		default:
			return new(EventUnknown)
		}
	}

}

// Process events. Returns true if EventQuit has appeared
func SlurpEvents() (quit bool) {
	quit = false
	for {
		ev := PollEvent()
		if ev == nil {
			return
		}
		switch ev.(type) {
		case *EventQuit:
			quit = true
		}
	}
}

// Returns true if key is pressed, false otherwise
func IsKeyPressed(kcode int) bool {
	return C.isKeyPressed(C.int(kcode)) == 1
}
