package electronbridge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	CLEAR_COLOUR = 'w'
)

var keypress_chan = make(chan string)
var key_query_chan = make(chan chan string)
var keyclear_chan = make(chan bool)

var mousedown_chan = make(chan Point)
var mouse_query_chan = make(chan chan Point)

var effect_done_channels = make(map[int]chan bool)
var effect_done_channels_MUTEX sync.Mutex

// ----------------------------------------------------------

type id_object struct {
	current			int
}

func (i *id_object) next() int {
	i.current += 1
	return i.current
}

var id_maker			id_object
var effect_id_maker		id_object

// ----------------------------------------------------------

type ByteSlice []byte		// Define this so that such a thing can have its own MarshalJSON() method

func (b ByteSlice) MarshalJSON() ([]byte, error) {
	str := string(b)
	return json.Marshal(str)
}

// ----------------------------------------------------------

type Point struct {
	X				int						`json:"x"`
	Y				int						`json:"y"`
}

type Spot struct {
	Char			byte
	Colour			byte
}

// ----------------------------------------------------------

type Window struct {
	Uid				int						`json:"uid"`
	Width			int						`json:"width"`
	Height			int						`json:"height"`
	Chars			ByteSlice				`json:"chars"`
	Colours			ByteSlice				`json:"colours"`
	Highlight		Point					`json:"highlight"`
}

// ----------------------------------------------------------

type NewMsgContent struct {
	Name			string					`json:"name"`
	Page			string					`json:"page"`
	Uid				int						`json:"uid"`
	Width			int						`json:"width"`
	Height			int						`json:"height"`
	BoxWidth		int						`json:"boxwidth"`
	BoxHeight		int						`json:"boxheight"`
	FontPercent		int						`json:"fontpercent"`
	Resizable		bool					`json:"resizable"`
}

type NewMsg struct {
	Command			string					`json:"command"`
	Content			NewMsgContent			`json:"content"`
}

type FlipMsg struct {
	Command			string					`json:"command"`
	Content			*Window					`json:"content"`
}

type AlertMsg struct {
	Command			string					`json:"command"`
	Content			string					`json:"content"`
}

type SpecialMsgContent struct {
	Effect			string					`json:"effect"`
	EffectID		int						`json:"effectid"`
	Uid				int						`json:"uid"`
	Args			[]interface{}			`json:"args"`
}

type SpecialMsg struct {
	Command			string					`json:"command"`
	Content			SpecialMsgContent		`json:"content"`
}

// ----------------------------------------------------------

type IncomingMsgType struct {
	Type			string					`json:"type"`
}

// ----------------------------------------------------------

type IncomingKeyContent struct {
	Down			bool					`json:"down"`
	Uid				int						`json:"uid"`
	Key				string					`json:"key"`
}

type IncomingKey struct {
	Type			string					`json:"type"`
	Content			IncomingKeyContent		`json:"content"`
}

// ----------------------------------------------------------

type IncomingMouseContent struct {
	Down			bool					`json:"down"`
	Uid				int						`json:"uid"`
	X				int						`json:"x"`
	Y				int						`json:"y"`
}

type IncomingMouse struct {
	Type			string					`json:"type"`
	Content			IncomingMouseContent	`json:"content"`
}

// ----------------------------------------------------------

type IncomingEffectDoneContent struct {
	Uid				int							`json:"uid"`
	EffectID		int							`json:"effectid"`
}

type IncomingEffectDone struct {
	Type			string						`json:"type"`
	Content			IncomingEffectDoneContent	`json:"content"`
}

// ----------------------------------------------------------

func init() {
	go listener()
	go keymaster()
	go mousemaster()
}

func listener() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		scanner.Scan()

		Logf("%v", scanner.Text())

		if strings.TrimSpace(scanner.Text()) == "" {
			continue
		}

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

		if type_obj.Type == "mouse" {

			var mouse_msg IncomingMouse

			err := json.Unmarshal(scanner.Bytes(), &mouse_msg)

			if err != nil {
				continue
			}

			if mouse_msg.Content.Down {
				mousedown_chan <- Point{mouse_msg.Content.X, mouse_msg.Content.Y}
			}
		}

		if type_obj.Type == "effect_done" {

			var effect_done_msg IncomingEffectDone

			err := json.Unmarshal(scanner.Bytes(), &effect_done_msg)

			if err != nil {
				continue
			}

			effect_done_channels_MUTEX.Lock()
			ch := effect_done_channels[effect_done_msg.Content.EffectID]
			effect_done_channels_MUTEX.Unlock()

			if ch != nil {
				go effect_notifier(ch)
			} else {
				Logf("Received done for effect %d but no notifier was known", effect_done_msg.Content.EffectID)
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
		case <- keyclear_chan:
			keyqueue = nil
		}
	}
}

func GetKeypress() (string, error) {

	response_chan := make(chan string)
	key_query_chan <- response_chan

	key := <- response_chan
	var err error = nil

	if key == "" {
		err = fmt.Errorf("GetKeypress(): nothing on queue")
	}

	return key, err
}

func ClearKeyQueue() {
	keyclear_chan <- true
}

func mousemaster() {

	var mousequeue []Point

	for {
		select {
		case response_chan := <- mouse_query_chan:
			if len(mousequeue) == 0 {
				response_chan <- Point{-1, -1}					// Note this: -1, -1 is used as a flag for empty queue
			} else {
				response_chan <- mousequeue[0]
				mousequeue = mousequeue[1:]
			}
		case mousedown := <- mousedown_chan:
			mousequeue = append(mousequeue, mousedown)
		}
	}
}

func GetMousedown() (Point, error) {

	response_chan := make(chan Point)
	mouse_query_chan <- response_chan

	point := <- response_chan
	var err error = nil

	if point.X < 0 {											// Note this: -1, -1 is used as a flag for empty queue
		err = fmt.Errorf("GetMousedown(): nothing on queue")
	}

	return point, err
}

func effect_notifier(ch chan bool) {
	ch <- true
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

func (w *Window) SetPointSpot(point Point, spot Spot) {
	w.Set(point.X, point.Y, spot.Char, spot.Colour)
}

func (w *Window) Get(x, y int) Spot {
	index := y * w.Width + x
	if index < 0 || index >= len(w.Chars) || x < 0 || x >= w.Width || y < 0 || y >= w.Height {
		return Spot{Char: ' ', Colour: CLEAR_COLOUR}
	}
	char := w.Chars[index]
	colour := w.Colours[index]
	return Spot{Char: char, Colour: colour}
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
		Command: "update",
		Content: w,
	}

	s, err := json.Marshal(m)
	if err != nil {
		panic("Failed to Marshal")
	}

	fmt.Printf("%s\n", string(s))
}

func (w *Window) Special(effect string, args []interface{}) {

	c := SpecialMsgContent{
		Effect: effect,
		Uid: w.Uid,
		EffectID: effect_id_maker.next(),
		Args: args,
	}

	m := SpecialMsg{
		Command: "special",
		Content: c,
	}

	s, err := json.Marshal(m)
	if err != nil {
		panic("Failed to Marshal")
	}

	// We make a channel for the purpose of receiving a message when the effect completes,
	// and add it to the global map of such channels.

	ch := make(chan bool)

	timeout := time.NewTimer(5 * time.Second)

	effect_done_channels_MUTEX.Lock()
	effect_done_channels[c.EffectID] = ch
	effect_done_channels_MUTEX.Unlock()

	fmt.Printf("%s\n", string(s))

	// Now we wait for the message that the effect completed...
	// Or the timeout ticker to fire.

	ChanLoop:
	for {
		select {
		case <- ch:
			break ChanLoop
		case <- timeout.C:
			break ChanLoop
		}
	}

	effect_done_channels_MUTEX.Lock()
	delete(effect_done_channels, c.EffectID)
	effect_done_channels_MUTEX.Unlock()
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

func Alertf(format_string string, args ...interface{}) {

	msg := fmt.Sprintf(format_string, args...)

	m := AlertMsg{
		Command: "alert",
		Content: msg,
	}

	s, err := json.Marshal(m)
	if err != nil {
		panic("Failed to Marshal")
	}
	fmt.Printf("%s\n", string(s))
}

func Logf(format_string string, args ...interface{}) {

	msg := fmt.Sprintf(format_string, args...)

	if len(msg) < 1 {
		return
	}

	if msg[len(msg) - 1] != '\n' {
		msg += "\n"
	}

	fmt.Fprintf(os.Stderr, "%s", msg)
}
