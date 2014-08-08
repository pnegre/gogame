package gogame

/*
#cgo pkg-config: sdl2
#include "SDL.h"

int getEventType(SDL_Event *e) {
    return e->type;
}

int getKeyCode(SDL_Event *e) {
    return e->key.keysym.scancode;
}

int isKeyRepeat(SDL_Event *e) {
    return e->key.repeat;
}

int isKeyPressed(int kcode) {
    const Uint8 * state = SDL_GetKeyboardState(NULL);
    return state[kcode];
    // int sc = SDL_GetScancodeFromKey(kcode);
    // return state[sc];
}

int getWindowEvent(SDL_Event *e) {
    return e->window.event;
}

void getWinData(SDL_Event *e, int *w, int *h) {
    *w = e->window.data1;
    *h = e->window.data2;
}

void getMouseCoords(SDL_Event *e, int *x, int *y, int *down) {
    *x = e->button.x;
    *y = e->button.y;
    if (e->button.state == SDL_PRESSED) {
        *down = 1;
    } else {
        *down = 0;
    }
}

void getMouseWheel(SDL_Event *e, int *x, int *y) {
    *x = e->wheel.x;
    *y = e->wheel.y;
}

*/
import "C"

const (
	K_LEFT   = C.SDL_SCANCODE_LEFT
	K_RIGHT  = C.SDL_SCANCODE_RIGHT
	K_UP     = C.SDL_SCANCODE_UP
	K_DOWN   = C.SDL_SCANCODE_DOWN
	K_SPACE  = C.SDL_SCANCODE_SPACE
	K_ESC    = C.SDL_SCANCODE_ESCAPE
	K_RETURN = C.SDL_SCANCODE_RETURN
	K_A      = C.SDL_SCANCODE_A
	K_B      = C.SDL_SCANCODE_B
	K_C      = C.SDL_SCANCODE_C
	K_D      = C.SDL_SCANCODE_D
	K_E      = C.SDL_SCANCODE_E
	K_F      = C.SDL_SCANCODE_F
	K_G      = C.SDL_SCANCODE_G
	K_H      = C.SDL_SCANCODE_H
	K_I      = C.SDL_SCANCODE_I
	K_J      = C.SDL_SCANCODE_J
	K_K      = C.SDL_SCANCODE_K
	K_L      = C.SDL_SCANCODE_L
	K_M      = C.SDL_SCANCODE_M
	K_N      = C.SDL_SCANCODE_N
	K_O      = C.SDL_SCANCODE_O
	K_P      = C.SDL_SCANCODE_P
	K_Q      = C.SDL_SCANCODE_Q
	K_R      = C.SDL_SCANCODE_R
	K_S      = C.SDL_SCANCODE_S
	K_T      = C.SDL_SCANCODE_T
	K_U      = C.SDL_SCANCODE_U
	K_V      = C.SDL_SCANCODE_V
	K_W      = C.SDL_SCANCODE_W
	K_X      = C.SDL_SCANCODE_X
	K_Y      = C.SDL_SCANCODE_Y
	K_Z      = C.SDL_SCANCODE_Z
)

type Event interface{}

type EventQuit interface{}

type EventUnknown interface{}

type EventExposed interface{}

type EventMouseClick struct {
	X, Y int
	Down bool
}

type EventMouseWheel struct {
	X, Y int
}

type EventKey struct {
	Code int
	Down bool
}

type EventResize struct {
	W, H int
}

// Wait for next event
func WaitEvent() Event {
	var cev C.SDL_Event
	if 0 == C.SDL_WaitEvent(&cev) {
		return nil
	}
	return classifyEvent(&cev)
}

func WaitEventTimeout(timeout int) Event {
	var cev C.SDL_Event
	if 0 == C.SDL_WaitEventTimeout(&cev, C.int(timeout)) {
		return nil
	}
	return classifyEvent(&cev)
}

// Poll for pending envents. Return nil if there is no event available
func PollEvent() Event {
	var cev C.SDL_Event
	if 0 == C.SDL_PollEvent(&cev) {
		return nil
	}
	return classifyEvent(&cev)
}

func classifyEvent(cev *C.SDL_Event) Event {
	switch C.getEventType(cev) {

	case C.SDL_QUIT:
		return new(EventQuit)

	case C.SDL_KEYDOWN:
		// Ignore repeat key events
		if C.isKeyRepeat(cev) != 0 {
			break
		}
		kde := new(EventKey)
		kde.Code = int(C.getKeyCode(cev))
		kde.Down = true
		return kde

	case C.SDL_KEYUP:
		kde := new(EventKey)
		kde.Code = int(C.getKeyCode(cev))
		kde.Down = false
		return kde

	case C.SDL_WINDOWEVENT:
		wet := C.getWindowEvent(cev)
		if wet == C.SDL_WINDOWEVENT_RESIZED {
			var w, h C.int
			C.getWinData(cev, &w, &h)
			wr := new(EventResize)
			wr.W = int(w)
			wr.H = int(h)
			return wr

		} else if wet == C.SDL_WINDOWEVENT_EXPOSED {
			return new(EventExposed)
		}

	case C.SDL_MOUSEBUTTONDOWN:
		var x, y, dw C.int
		var down bool
		C.getMouseCoords(cev, &x, &y, &dw)
		if dw == 1 {
			down = true
		}
		return &EventMouseClick{X: int(x), Y: int(y), Down: down}

	case C.SDL_MOUSEWHEEL:
		var x, y C.int
		C.getMouseWheel(cev, &x, &y)
		return &EventMouseWheel{int(x), int(y)}
	}
	return new(EventUnknown)
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
