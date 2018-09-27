package log

import (
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func WriteLog() {
	d := []byte("hello\ngo\n")
	f, err := os.Create("1.log")
	check(err)

	defer f.Close()

	d = []byte("haha\ngo\n")

	f.Write(d)

}
