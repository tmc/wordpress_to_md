// Package wp2md converts a wordpress export into a jkl-compatible list of markdown files
package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/tmc/utils/pathext"
	wp2md "github.com/tmc/wordpress_to_md"
	"github.com/tmc/wxr"
)

var (
	inputFile = flag.String("input", "", "input wxr XML file")
	outputDir = flag.String("outputDir", ".", "output directory for markdown files")
)

func main() {
	flag.Parse()

	if *inputFile == "" {
		flag.PrintDefaults()
		flag.Usage()
		os.Exit(1)
	}
	log.Println("Opening", *inputFile)

	var (
		r   io.Reader
		err error
		rss *wxr.RSS
	)
	if *inputFile, err = pathext.ExpandUser(*inputFile); err != nil {
		log.Fatalf("error expanding user: %s\n", err)
	}
	r, err = os.Open(*inputFile)
	if err != nil {
		log.Fatalf("error opening %s: %s\n", *inputFile, err)
	}

	if rss, err = wxr.NewRSS(r); err != nil {
		log.Fatalln(err)
	}

	if files, err := wp2md.RSSToFiles(rss); err != nil {
		log.Fatalln(err)
	} else {
		for _, file := range files {
			outputFile := filepath.Join(*outputDir, file.Filename)
			if outputFile, err = pathext.ExpandUser(outputFile); err != nil {
				log.Fatalf("error expanding user: %s\n", err)
			}
			err := ioutil.WriteFile(outputFile, file.Contents, 0644)
			if err != nil {
				log.Fatalf("error writing to %s: %s\n", outputFile, err)
			}
			log.Println("wrote", outputFile)
		}
	}
}
