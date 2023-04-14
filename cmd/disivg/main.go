package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/reactivego/ivg/decode"
)

func main() {
	var out = flag.String("o", "stdout", "the filename to write the disassembled IVG data to")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "%[1]s is a tool for disassembling IVG icons.\n\n"+
			"Usage:\n\n"+
			"  %[1]s [flags] filepath\n\n"+
			"The flags are:\n\n", flag.CommandLine.Name())
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output())
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	filename := flag.Arg(0)
	ivgData, err := os.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		log.Fatalf("%s: ReadFile: %v", filename, err)
	}
	dis, err := decode.Disassemble(ivgData)
	if err != nil {
		log.Fatalf("%s: disassemble: %v", filename, err)
	}
	if out == nil || *out == "stdout" {
		_, err := os.Stdout.Write(dis)
		if err != nil {
			log.Fatalf("%s: Write: %v", *out, err)
		}
	} else if err := os.WriteFile(filepath.FromSlash(filepath.FromSlash(*out)), dis, 0666); err != nil {
		log.Fatalf("%s: WriteFile: %v", filename, err)
	}
}
