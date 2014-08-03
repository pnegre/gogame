package gogame

/*
#cgo pkg-config: sdl2 SDL2_image
#include "SDL.h"
#include "SDL_image.h"

SDL_Texture * makeTexture( char *f, SDL_Renderer *ren ) {
    SDL_Texture *tex = IMG_LoadTexture(ren, f);
    return tex;
}

void render( SDL_Renderer *ren, SDL_Texture *tex, int x, int y, int w, int h) {
    static SDL_Rect dst;
    dst.x = x;
    dst.y = y;
    dst.w = w;
    dst.h = h;
    SDL_RenderCopy(ren, tex, NULL, &dst);
}

void render2( SDL_Renderer *ren, SDL_Texture *tex, int x, int y, int w, int h, int ox, int oy) {
    static SDL_Rect org;
    static SDL_Rect dst;
    dst.x = x;
    dst.y = y;
    dst.w = w;
    dst.h = h;
    org.x = ox;
    org.y = oy;
    org.w = w;
    org.h = h;
    SDL_RenderCopy(ren, tex, &org, &dst);
}

void render3( SDL_Renderer *ren, SDL_Texture *tex, int w, int h, int x, int y, int dx, int dy) {
    static SDL_Rect org;
    static SDL_Rect dst;
    org.x = 0;
    org.y = 0;
    org.w = w;
    org.h = h;
    dst.x = x;
    dst.y = y;
    dst.w = dx;
    dst.h = dy;
    SDL_RenderCopy(ren, tex, &org, &dst);
}

void queryTexture(SDL_Texture *t, int *h, int *v) {
    SDL_QueryTexture(t, NULL, NULL, h, v);
}

int intersects(int x1, int y1, int w1, int h1, int x2, int y2, int w2, int h2) {
    static SDL_Rect a;
    static SDL_Rect b;
    a.x = x1; a.y = y1; a.w = w1; a.h = h1;
    b.x = x2; b.y = y2; b.w = w2; b.h = h2;
    return SDL_HasIntersection(&a, &b);
}


*/
import "C"

type Drawable interface {
	BlitRect(*Rect)
	GetDimensions() (int, int)
}

type Texture struct {
	tex *C.SDL_Texture
	w   int
	h   int
}

type Rect struct {
	X, Y, W, H int
}

// Set center of Rect to x,y
func (self *Rect) SetCenter(x, y int) {
	self.X = x - self.W/2
	self.Y = y - self.H/2
}

// Determine if two rectangles intersect
func (self *Rect) Intersects(r2 *Rect) bool {
	if 0 == C.intersects(C.int(self.X), C.int(self.Y), C.int(self.W), C.int(self.H),
		C.int(r2.X), C.int(r2.Y), C.int(r2.W), C.int(r2.H)) {
		return false
	}

	return true
}

// Construct a new texture from a image file
func NewTexture(filename string) *Texture {
	tex := C.makeTexture(C.CString(filename), renderer)
	return getNewTexture(tex)
}

func getNewTexture(tex *C.SDL_Texture) *Texture {
	t := new(Texture)
	t.tex = tex
	var w C.int
	var h C.int
	C.queryTexture(t.tex, &w, &h)
	t.w, t.h = int(w), int(h)
	// El finalizer no funciona. Suposo, fent proves, que és perquè el Garbage Collector
	// s'executa dins la seva pròpia goroutine, i sembla que no es pot cridar a
	// SDL_DestroyTexture si no és des del thread original
	// 	runtime.SetFinalizer(t, func(t *Texture) {
	// 		log.Printf("Finalizing texture...")
	// 		C.SDL_DestroyTexture(t.tex)
	// 	})
	return t
}

// Destroy texture. Must be called explicitly, no automatic free for texture data
func (self *Texture) Destroy() {
	C.SDL_DestroyTexture(self.tex)
}

// Get texture dimensions, (horizontal, vertical)
func (self *Texture) GetDimensions() (int, int) {
	return self.w, self.h
}

// Blit texture to screen, using provided rect
func (self *Texture) BlitRect(r *Rect) {
	C.render(renderer, self.tex, C.int(r.X), C.int(r.Y),
		C.int(r.W), C.int(r.H))
}

// Blit texture to screen, using x,y coordinates (top-left corner)
func (self *Texture) BlitToScreen(x, y int) {
	C.render(renderer, self.tex, C.int(x), C.int(y),
		C.int(self.w), C.int(self.h))
}

// Blit texture to screen, stretching the image
func (self *Texture) BlitToScreenStretch(dst *Rect) {
	C.render3(renderer, self.tex, C.int(self.w), C.int(self.h),
		C.int(dst.X), C.int(dst.Y), C.int(dst.W), C.int(dst.H))
}

// Get subtexture
func (self *Texture) SubTex(x, y, w, h int) *SubTexture {
	return &SubTexture{self, &Rect{x, y, w, h}}
}

type SubTexture struct {
	tex  *Texture
	rect *Rect
}

// Blit subtexture to screen, using provided rect
func (self *SubTexture) BlitRect(r *Rect) {
	C.render2(renderer, self.tex.tex, C.int(r.X), C.int(r.Y), C.int(r.W),
		C.int(r.H), C.int(self.rect.X), C.int(self.rect.Y))
}

// Get subtexture dimensions
func (self *SubTexture) GetDimensions() (int, int) {
	return self.rect.W, self.rect.H
}
