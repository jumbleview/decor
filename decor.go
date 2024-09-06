// go build -ldflags -H=windowsgui
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jumbleview/decor/screen"
)

func main() {
	var isSilent bool
	flag.BoolVar(&isSilent, "s", false, "silent mode")
	flag.Parse()
	cmd := flag.Args()

	if !isSilent {
		screen.MakeConsole()
	}
	greeting := "decor.exe [-s] path_to_jpeg_file"
	var Usage = func() {
		fmt.Println(greeting)
	}
	if len(cmd) != 1 {
		fmt.Printf("Wrong number of arguments %d\n", len(cmd))
		Usage()
		fmt.Printf("Hit \"Enter\" to exit... ")
		s := ""
		fmt.Scanln(&s)
		os.Exit(1)
	}

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	execPath := filepath.Dir(ex)
	err = screen.SetWallpaper(cmd[0], execPath)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("Hit \"Enter\" to exit... ")
	s := ""
	fmt.Scanln(&s)
}
