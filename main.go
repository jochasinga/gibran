package main

import (
	"flag"
	"fmt"
	"os"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage of Gibran CLI:")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, " gibran new <project_name>")
		fmt.Fprintln(os.Stderr, "")
		flag.PrintDefaults()
	}

	flag.CommandLine.Init("", flag.ExitOnError)
}

func main() {

}
