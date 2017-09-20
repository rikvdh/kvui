package main

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
)

const (
	treeView   = "tree"
	valueView  = "value"
	statusView = "status"

	keyPrefix = "  - "
	dbPrefix  = " db:"
)

var (
	currentView = treeView
	currentDb   = 0
	currentKey  = ""
)

func renderTree(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	databases, err := kvstore.Databases()
	if err != nil {
		fmt.Fprintln(v, err)
	}
	for i := 0; i < databases; i++ {
		if i == currentDb {
			fmt.Fprintf(v, "-%s%d\n", dbPrefix, i)
			keys, err := kvstore.Keys("*")
			if err != nil {
				return err
			}
			for _, k := range keys {
				fmt.Fprintf(v, "%s%s\n", keyPrefix, k)
			}
		} else {
			fmt.Fprintf(v, "+%s%d\n", dbPrefix, i)
		}
	}
	_, pos := v.Cursor()
	l, err := v.Line(pos)
	if err != nil {
		fmt.Fprintf(v, "error: %v", err)
		return nil
	}
	if strings.HasPrefix(l, keyPrefix) {
		currentKey = l[len(keyPrefix):]
		vv, _ := g.View(valueView)
		renderValue(vv)
	}
	return nil
}

func renderValue(v *gocui.View) error {
	v.Clear()
	if currentKey != "" {
		t, err := kvstore.Type(currentKey)
		fmt.Fprintln(v, t, err)
		s, err := kvstore.Get(currentKey)
		fmt.Fprintln(v, s, err)
	} else {
		fmt.Fprintln(v, time.Now().Format(time.Stamp), currentView)
	}
	return nil
}

func renderStatus(v *gocui.View) error {
	v.Clear()
	con, conerr := kvstore.Connected()
	if con {
		v.FgColor = gocui.ColorGreen
		fmt.Fprintf(v, " connected")
	} else {
		v.FgColor = gocui.ColorRed
		fmt.Fprintf(v, " disconnected (%v)", conerr)
	}

	return nil
}

func renderLayout(g *gocui.Gui) error {
	sizeX, sizeY := g.Size()
	treeSize := int(math.Floor(float64(sizeX) * 0.2))

	_, err := g.SetView(treeView, 0, 0, treeSize, sizeY-4)
	if err != nil {
		return err
	}
	_, err = g.SetView(valueView, treeSize+1, 0, sizeX-1, sizeY-4)
	if err != nil {
		return err
	}
	_, err = g.SetView(statusView, 0, sizeY-3, sizeX-1, sizeY-1)
	return err
}

func switchViewRight(g *gocui.Gui, v *gocui.View) error {
	if currentView == treeView {
		_, err := g.SetCurrentView(valueView)
		if err == nil {
			currentView = valueView
		}
	}
	return nil
}

func switchViewLeft(g *gocui.Gui, v *gocui.View) error {
	if currentView == valueView {
		_, err := g.SetCurrentView(treeView)
		if err == nil {
			currentView = treeView
		}
	}
	return nil
}

func redraw(g *gocui.Gui, v *gocui.View) error {
	nv, err := g.View(currentView)
	if err != nil {
		return err
	}
	switch currentView {
	case treeView:
		renderTree(g, nv)
	case valueView:
		renderValue(nv)
	}
	return nil
}
