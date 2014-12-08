package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

var checkTime bool

func startTimer() {
	checkTime = false
	time.Sleep(time.Second * 30)
	checkTime = true
}

func oldDelete(filename string, rotnum int) {
	const lookBackDays int = 7
	for i := 1; i < lookBackDays+1; i++ {
		t := time.Now().AddDate(0, 0, i*-1-rotnum)
		fn := fmt.Sprintf("%s-%d%02d%02d", filename, t.Year(), t.Month(), t.Day())
		if _, err := os.Stat(filename); !os.IsNotExist(err) {
			os.Remove(fn)
		}
	}
}

func main() {
	var rotnum int
	const fileflag int = os.O_CREATE | os.O_APPEND | os.O_RDWR
	const filemode os.FileMode = 0644

	flag.IntVar(&rotnum, "n", 7, "The number of rotated files to keep")
	flag.Parse()

	filename := flag.Arg(0)

	if len(filename) < 1 {
		fmt.Printf("Usage: %s [options] <filepath>\nType -h for options\n", os.Args[0])
		os.Exit(1)
	}

	day := time.Now().Day()

	f, err := os.OpenFile(filename, fileflag, filemode)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	reader := bufio.NewReader(os.Stdin)

	go oldDelete(filename, rotnum)
	go startTimer()

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if checkTime {
			if time.Now().Day() != day {
				if _, err := os.Stat(filename); !os.IsNotExist(err) {
					t := time.Now()
					filenameRot := fmt.Sprintf("%s-%d%02d%02d", filename, t.Year(), t.Month(), t.Day())
					f.Close()
					os.Rename(filename, filenameRot)
					go oldDelete(filename, rotnum)
					f, err = os.OpenFile(filename, fileflag, filemode)
				}
				day = time.Now().Day()
			}
			go startTimer()
		}

		f.WriteString(line)
	}
}
