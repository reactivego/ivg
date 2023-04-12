package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/reactivego/ivg/mdicons"
)

func main() {
	var pkg = flag.String("package", "icons", "the name of the package to generate")
	var test = flag.Bool("test", false, "pass this flag to generate data_test.go")
	var size = flag.Float64("size", 48, "width and height (in ideal vector space) of the "+
		"generated IVG graphic, regardless of the size of the input SVG")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%[1]s is a tool for converting icons from svg to ivg.\n\n"+
			"Usage:\n\n"+
			"  %[1]s [flags] directory\n\n"+
			"The flags are:\n\n", flag.CommandLine.Name())
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output())
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	if err := mdicons.Parse(flag.Arg(0), *pkg, true, *test, float32(*size)); err != nil {
		log.Fatal(err.Error())
	}
}
