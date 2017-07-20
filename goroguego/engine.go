package goroguego

import (
	"encoding/json"
	"fmt"
)

const (
	CLEAR_COLOUR = 'w'
)

// ----------------------------------------------------------

type id_object struct {
	current			int
}

func (i *id_object) next() int {
	i.current += 1
	return i.current
}

var id_maker id_object;

// ----------------------------------------------------------

type Window struct {
	Uid				int
	Width			int
	Height			int
	Chars			[]byte
	Colours			[]byte
}

// ----------------------------------------------------------

type NewMsgContent struct {
	Name			string			`json:"name"`
	Page			string			`json:"page"`
	Uid				int				`json:"uid"`
	Width			int				`json:"width"`
	Height			int				`json:"height"`
	BoxWidth		int				`json:"boxwidth"`
	BoxHeight		int				`json:"boxheight"`
	FontPercent		int				`json:"fontpercent"`
	Resizable		bool			`json:"resizable"`
}

type NewMsg struct {
	Command			string			`json:"command"`
	Content			NewMsgContent	`json:"content"`
}

// ----------------------------------------------------------

type FlipMsgContent struct {
	Uid				int				`json:"uid"`
	Chars			string			`json:"chars"`
	Colours			string			`json:"colours"`
}

type FlipMsg struct {
	Command			string			`json:"command"`
	Content			FlipMsgContent	`json:"content"`
}

// ----------------------------------------------------------

func (w *Window) Set(x, y int, char, colour byte) {
	index := y * w.Width + x
	if index < 0 || index >= len(w.Chars) || x < 0 || x >= w.Width || y < 0 || y >= w.Height {
		panic("index in set()")
	}
	w.Chars[index] = char
	w.Colours[index] = colour
}

func (w *Window) Clear() {
	for n := 0; n < len(w.Chars); n++ {
		w.Chars[n] = ' '
		w.Colours[n] = CLEAR_COLOUR
	}
}

func (w *Window) Flip() {

	m := FlipMsg{Command: "flip", Content: FlipMsgContent{
			Uid: w.Uid,
			Chars: string(w.Chars),
			Colours: string(w.Colours),
		},
	}

	s, err := json.Marshal(m)
	if err != nil {
		panic("Failed to Marshal")
	}
	fmt.Printf("%s\n", string(s))
}

func NewWindow(name, page string, width, height, boxwidth, boxheight, fontpercent int, resizable bool) *Window {

	uid := id_maker.next()

	w := Window{Uid: uid, Width: width, Height: height}

	w.Chars = make([]byte, width * height)
	w.Colours = make([]byte, width * height)

	// Create the message to send to the server...

	m := NewMsg{Command: "new", Content: NewMsgContent{
			Name: name,
			Page: page,
			Uid: uid,
			Width: width,
			Height: height,
			BoxWidth: boxwidth,
			BoxHeight: boxheight,
			FontPercent: fontpercent,
			Resizable: resizable,
		},
	}

	s, err := json.Marshal(m)
	if err != nil {
		panic("Failed to Marshal")
	}
	fmt.Printf("%s\n", string(s))

	return &w
}
