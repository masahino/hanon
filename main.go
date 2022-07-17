package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	path := os.Getenv("HANON_PATH")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n  %s [-n] X.Y\n", os.Args[0], os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	var speed bool
	flag.BoolVar(&speed, "n", false, "use natural speed")
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
	}
	num := flag.Arg(0)
	opt := "slow"
	if speed == true {
		opt = "natural"
	}
	filepattern := path + "* " + num + "_" + opt + ".mp3"

	files, err := filepath.Glob(filepattern)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		fmt.Println(f)
		playMp3(f)
	}
}

func playMp3(path string) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	st, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer st.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	done := make(chan bool)
	speaker.Play(beep.Seq(st, beep.Callback(func() {
		done <- true
	})))
	<-done
}
