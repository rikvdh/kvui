package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/rikvdh/redisui/kv"
)

var (
	no256   *bool
	host    *string
	port    *uint
	kvtype  *string
	kvstore kv.KV
)

func init() {
	no256 = flag.Bool("no256", false, "Disable 256-color")
	host = flag.String("h", "localhost", "Host to connect to")
	port = flag.Uint("p", 6379, "Port to connect to")
	kvtype = flag.String("type", "redis", "KV-storage type")
}

func exit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func main() {
	flag.Parse()
	c := gocui.Output256
	if *no256 {
		c = gocui.OutputNormal
	}
	g, err := gocui.NewGui(c)
	if err != nil {
		panic(err)
	}
	defer g.Close()

	kvstore, err = kv.New(*kvtype, fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		panic(err)
	}

	log.SetOutput(os.Stderr)
	g.SetManagerFunc(renderLayout)
	g.SelFgColor = gocui.ColorYellow
	g.Highlight = true
	g.Cursor = true

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, exit); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, switchViewRight); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, switchViewLeft); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, -1, true)
		return redraw(g, v)
	}); err != nil {
		panic(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		v.MoveCursor(0, 1, true)
		return redraw(g, v)
	}); err != nil {
		panic(err)
	}

	sizeX, sizeY := g.Size()
	treeSize := int(math.Floor(float64(sizeX) * 0.2))
	treeView, err := g.SetView(treeView, 0, 0, treeSize, sizeY-4)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	treeView.Wrap = true
	treeView.Highlight = true
	treeView.SelBgColor = gocui.ColorWhite
	treeView.SelFgColor = gocui.ColorBlack
	renderTree(g, treeView)
	valueView, err := g.SetView(valueView, treeSize+1, 0, sizeX-1, sizeY-4)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	valueView.Wrap = true
	renderValue(valueView)

	statusView, err := g.SetView(statusView, 0, sizeY-3, sizeX-1, sizeY-1)
	if err != nil && err != gocui.ErrUnknownView {
		panic(err)
	}
	go func() {
		tm := time.NewTicker(time.Second)
		for range tm.C {
			g.Update(func(g *gocui.Gui) error {
				return renderStatus(statusView)
			})
		}
	}()

	g.SetCurrentView(currentView)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
	}
}
