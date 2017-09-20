// Copyright 2017 The KVUI Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/rikvdh/redisui/kv/types"
)

const (
	treeView     = "tree"
	valueView    = "value"
	statusView   = "status"
	subValueView = "subvalue"

	keyPrefix = "  - "
	dbPrefix  = " db:"
)

var (
	currentView    = treeView
	currentDb      = 0
	currentKey     = ""
	currentKeyType = types.KVTypeInvalid
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
	} else if strings.HasPrefix(l, "-"+dbPrefix) || strings.HasPrefix(l, "+"+dbPrefix) {
		currentKey = ""
	}

	vv, _ := g.View(valueView)
	if vv != nil {
		if err := renderValue(g, vv); err != nil {
			sv, _ := g.View(statusView)
			renderStatus(sv, err)
		}
	}
	return nil
}

func dbSelect(g *gocui.Gui, v *gocui.View) error {
	if v.Name() == treeView {
		_, pos := v.Cursor()
		l, err := v.Line(pos)
		if err != nil {
			fmt.Fprintf(v, "error: %v", err)
			return nil
		}
		if strings.HasPrefix(l, "+"+dbPrefix) {
			currentDb, err = strconv.Atoi(l[len("+"+dbPrefix):])
			kvstore.Database(currentDb)
			renderTree(g, v)
		} else if strings.HasPrefix(l, "-"+dbPrefix) {
		}
	}
	return nil
}

func renderValue(g *gocui.Gui, v *gocui.View) error {
	v.Clear()
	if currentKey != "" {
		t, err := kvstore.Type(currentKey)
		if err != nil {
			return err
		}
		if currentKeyType != t {
			currentKeyType = t
			renderLayout(g)
		}
		switch t {
		case types.KVTypeString:
			s, err := kvstore.Get(currentKey)
			if err != nil {
				return err
			}
			fmt.Fprintf(v, s)
		case types.KVTypeMap:
			s, err := kvstore.HKeys(currentKey)
			if err != nil {
				return err
			}
			for _, i := range s {
				fmt.Fprintf(v, "- %v\n", i)
			}
			_, p := v.Cursor()
			str, _ := v.Line(p)
			if len(str) >= 3 {
				renderSubValue(g, str[2:])
			}
		case types.KVTypeList:
			s, err := kvstore.LGet(currentKey)
			if err != nil {
				return err
			}
			for _, i := range s {
				fmt.Fprintf(v, "- %v\n", i)
			}
		}
	} else {
		fmt.Fprintln(v, time.Now().Format(time.Stamp), currentView)
	}
	return nil
}

func renderSubValue(g *gocui.Gui, field string) error {
	v, err := g.View(subValueView)
	if err != nil {
		return err
	}
	v.Clear()
	val, err := kvstore.HGet(currentKey, field)
	if err != nil {
		return err
	}
	fmt.Fprintf(v, val)
	return nil
}

var lastErr error

func renderStatus(v *gocui.View, err ...error) error {
	v.Clear()
	con, conerr := kvstore.Connected()
	if con {
		v.FgColor = gocui.ColorGreen
		fmt.Fprintf(v, " connected")
	} else {
		v.FgColor = gocui.ColorRed
		fmt.Fprintf(v, " disconnected (%v)", conerr)
	}

	if len(err) >= 1 {
		lastErr = err[0]
	}
	if lastErr != nil {
		fmt.Fprintf(v, "\t\tERROR: %v", lastErr)
	}
	return nil
}

func renderLayout(g *gocui.Gui) error {
	sizeX, sizeY := g.Size()
	treeSize = int(math.Floor(float64(sizeX) * 0.2))

	_, err := g.SetView(treeView, 0, 0, treeSize, sizeY-4)
	if err != nil {
		return err
	}
	if currentKeyType == types.KVTypeMap {
		vView, err := g.SetView(valueView, treeSize+1, 0, treeSize*2, sizeY-4)
		if err != nil {
			return err
		}
		vView.Highlight = true
		svView, err := g.SetView(subValueView, treeSize*2+1, 0, sizeX-1, sizeY-4)
		svView.Wrap = true
		if err != nil {
			return err
		}
	} else {
		vView, err := g.SetView(valueView, treeSize+1, 0, sizeX-1, sizeY-4)
		if err != nil {
			return err
		}
		vView.Highlight = false
		g.DeleteView(subValueView)
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
		err := renderValue(g, nv)
		if err != nil {
			sv, _ := g.View(statusView)
			renderStatus(sv, err)
		}
	}
	return nil
}
