# GOGAME

This is a simple 2D game library for Go. It's incomplete and I'm writing it for my personal projects. Feel free to suggest features!!

## Install instructions (linux only)

I use debian. If you're like me, you will need to install the following packages:

    libsdl2-dev
    libsdl2-image-dev
    libsdl2-ttf-dev

Others distros shoud have equivalent libraries avaliable.

And now, the easy part (make sure you have a valid $GOPATH):

    go get github.com/pnegre/gogame

## Documentation

In progress...

## Usage

Simple example:

    package main

    import (
        "github.com/pnegre/gogame"
        "github.com/pnegre/gogame/cache"
        "github.com/pnegre/gogame/sprite"
        "log"
    )

    const (
        WINTITLE = "test"
        WIN_W    = 800
        WIN_H    = 600
        IMAGE    = "someimage.png"
    )

    type Target struct {
        *sprite.Simple
        x, y int
    }

    func NewTarget() *Target {
        tex := cache.GetTexture(IMAGE)
        return &Target{sprite.NewSimple(tex), WIN_W / 2, WIN_H / 2}
    }

    func (self *Target) Update() {
        self.x += 10
        self.Rect.SetCenter(self.x, self.y)
        if self.x >= WIN_W {
            self.x = 0
        }
    }

    func main() {
        if err := gogame.Init(WINTITLE, WIN_W, WIN_H); err != nil {
            log.Fatal(err)
        }
        defer gogame.Quit()
        defer cache.DestroyAll()

        target := NewTarget()
        quit := false
        for !quit {
            for {
                var ev gogame.Event
                if ev = gogame.PollEvent(); ev == nil {
                    break
                }

                switch ev.(type) {

                case *gogame.QuitEvent:
                    quit = true
                }
            }

            target.Update()

            gogame.RenderClear()
            target.Draw()
            gogame.RenderPresent()

            gogame.Delay(50)
        }

    }
