package goroguego

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

const (
	CLEAR_COLOUR = 'w'
)

var keypress_chan = make(chan string)
var key_query_chan = make(chan chan string)

// var logfile, _ = os.Create("logfile.txt")

// ----------------------------------------------------------

type id_object struct {
	current			int
}

func (i *id_object) next() int {
	i.current += 1
	return i.current
}

var id_maker id_object

// ----------------------------------------------------------

type ByteSlice []byte		// Define this so that such a thing can have its own MarshalJSON() method

func (b ByteSlice) MarshalJSON() ([]byte, error) {
	str := string(b)
	return json.Marshal(str)
}

// ----------------------------------------------------------

type Point struct {
	X				int					`json:"x"`
	Y				int					`json:"y"`
}

// ----------------------------------------------------------

type Window struct {
	Uid				int					`json:"uid"`
	Width			int					`json:"width"`
	Height			int					`json:"height"`
	Chars			ByteSlice			`json:"chars"`
	Colours			ByteSlice			`json:"colours"`
	Highlight		Point				`json:"highlight"`
}

// ----------------------------------------------------------

type NewMsgContent struct {
	Name			string				`json:"name"`
	Page			string				`json:"page"`
	Uid				int					`json:"uid"`
	Width			int					`json:"width"`
	Height			int					`json:"height"`
	BoxWidth		int					`json:"boxwidth"`
	BoxHeight		int					`json:"boxheight"`
	FontPercent		int					`json:"fontpercent"`
	Resizable		bool				`json:"resizable"`
}

type NewMsg struct {
	Command			string				`json:"command"`
	Content			NewMsgContent		`json:"content"`
}

// ----------------------------------------------------------

type FlipMsg struct {
	Command			string				`json:"command"`
	Content			*Window				`json:"content"`
}

// ----------------------------------------------------------

type IncomingMsgType struct {
	Type			string				`json:"type"`
}

// ----------------------------------------------------------

type IncomingKeyContent struct {
	Down			bool				`json:"down"`
	Uid				int					`json:"uid"`
	Key				string				`json:"key"`
}

type IncomingKey struct {
	Type			string				`json:"type"`
	Content			IncomingKeyContent	`json:"content"`
}

// ----------------------------------------------------------

func init() {
	go listener()
	go keymaster()
}

func listener() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()

		// logfile.WriteString(scanner.Text() + "\n")

		var type_obj IncomingMsgType

		err := json.Unmarshal(scanner.Bytes(), &type_obj)
		if err != nil {
			continue
		}

		if type_obj.Type == "key" {

			var key_msg IncomingKey

			err := json.Unmarshal(scanner.Bytes(), &key_msg)

			if err != nil {
				continue
			}

			if key_msg.Content.Down {
				keypress_chan <- key_msg.Content.Key
			}
		}
	}
}

func keymaster() {

	var keyqueue []string

	for {
		select {
		case response_chan := <- key_query_chan:
			if len(keyqueue) == 0 {
				response_chan <- ""
			} else {
				response_chan <- keyqueue[0]
				keyqueue = keyqueue[1:]
			}
		case keypress := <- keypress_chan:
			keyqueue = append(keyqueue, keypress)
		}
	}
}

func GetKeypress() string {
	response_chan := make(chan string)
	key_query_chan <- response_chan
	key := <- response_chan
	return key
}

// ----------------------------------------------------------

func (w *Window) Set(x, y int, char, colour byte) {
	index := y * w.Width + x
	if index < 0 || index >= len(w.Chars) || x < 0 || x >= w.Width || y < 0 || y >= w.Height {
		return
	}
	w.Chars[index] = char
	w.Colours[index] = colour
}

func (w *Window) SetHighlight(x, y int) {
	w.Highlight = Point{x, y}
}

func (w *Window) Clear() {
	for n := 0; n < len(w.Chars); n++ {
		w.Chars[n] = ' '
		w.Colours[n] = CLEAR_COLOUR
	}
	w.Highlight = Point{-1, -1}
}

func (w *Window) Flip() {

	m := FlipMsg{
		Command: "flip",
		Content: w,
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

	w.Highlight = Point{-1, -1}

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
