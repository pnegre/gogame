package gogame

/*
#cgo pkg-config: SDL2_ttf
#include "SDL_ttf.h"

SDL_Texture *getTextureFromFont(SDL_Renderer *ren, TTF_Font *ttf, char *t, unsigned char r, unsigned char g, unsigned char b) {
    SDL_Color color = {r,g,b};
    SDL_Surface *s = TTF_RenderText_Blended(ttf, t, color);
    SDL_Texture *tex = SDL_CreateTextureFromSurface(ren, s);
    SDL_FreeSurface(s);
    return tex;
}

*/
import "C"

func init() {
	C.TTF_Init()
	// TODO: must call TTF_Quit() when quitting...
}

type Font struct {
	cttf *C.TTF_Font
}

// Load font from a file, use at "size" size
func NewFont(fname string, size int) *Font {
	f := C.TTF_OpenFont(C.CString(fname), C.int(size))
	// TODO: set finalizer (TTF_CloseFont)
	return &Font{f}
}

// Destroy font
func (self *Font) Destroy() {
	C.TTF_CloseFont(self.cttf)
}

// Get resulting texture of text rendered using the font object and color supplied
func (self *Font) GetTexture(text string, color *Color) *Texture {
	tex := C.getTextureFromFont(renderer, self.cttf, C.CString(text), C.uchar(color.R), C.uchar(color.G), C.uchar(color.B))
	return getNewTexture(tex)
}

// Render directly to screen. Creates texture and discards it after rendering to screen
func (self *Font) RenderToScreen(text string, x, y int, color *Color) {
	tex := self.GetTexture(text, color)
	tex.Blit(x, y)
	tex.Destroy()
}

// Same as RenderToScreen, but center the resulting texture in x,y coordinates
func (self *Font) RenderToScreenCenter(text string, x, y int, color *Color) {
	tex := self.GetTexture(text, color)
	w, h := tex.GetDimensions()
	xx := x - w/2
	yy := y - h/2
	tex.Blit(xx, yy)
	tex.Destroy()
}
