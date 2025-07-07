# GOGAME

This is a simple 2D game library for Go. It's incomplete and I'm writing it for my personal projects. Feel free to suggest features!!

## Documentation

Once you have installed the package, just run:

    go doc github.com/pnegre/gogame | less
    go doc github.com/pnegre/gogame/sprite | less
    go doc github.com/pnegre/gogame/cache | less

## Usage

Simple example:

    package main

    import (
        "github.com/pnegre/gogame"
        "github.com/pnegre/gogame/cache"
        "github.com/pnegre/gogame/sprite"
        "log"
        "runtime"
    )

    const (
        WINTITLE = "test"
        WIN_W    = 800
        WIN_H    = 600
        IMAGE    = "sprite.png"
    )

    type Target struct {
        *sprite.Simple
    }

    func NewTarget() *Target {
        tex := cache.GetTexture(IMAGE)
        s := sprite.NewSimple(tex)
        s.Rect.SetCenter(WIN_W/2, WIN_H/2)
        return &Target{s}
    }

    func (self *Target) Update() {
        if gogame.IsKeyPressed(gogame.K_LEFT) {
            self.Rect.X -= 10
        }
        if gogame.IsKeyPressed(gogame.K_RIGHT) {
            self.Rect.X += 10
        }
        if gogame.IsKeyPressed(gogame.K_UP) {
            self.Rect.Y -= 10
        }
        if gogame.IsKeyPressed(gogame.K_DOWN) {
            self.Rect.Y += 10
        }
    }

    func main() {
        runtime.LockOSThread()
        if err := gogame.InitSDL(); err != nil {
            log.Fatal(err)
        }
        if err := gogame.Init(WINTITLE, WIN_W, WIN_H); err != nil {
            log.Fatal(err)
        }
        defer gogame.Quit()
        defer cache.DestroyAll()

        target := NewTarget()

        for {
            if quit := gogame.SlurpEvents(); quit == true {
                break
            }

            target.Update()

            gogame.RenderClear()
            target.Draw()
            gogame.RenderPresent()

            gogame.Delay(16)
        }

    }