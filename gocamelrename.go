package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var recursive bool
var pretend bool
var verbose bool

func init() {
	flag.BoolVar(&verbose, "verbose", false, "Output more information")
	flag.BoolVar(&recursive, "recursive", true, "Rename files recursively")
	flag.BoolVar(&pretend, "pretend", true, "Spit out what would be done but don't do it")
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not open current working directory - %v", wd)
	}
	files, err := ioutil.ReadDir(wd)
	if err != nil {
		log.Fatalf("Could not list files in %s - %v", wd, err)
	}
	flag.Parse()
	renameFiles(files) // ignoring error, we don't care if the parent directory does not exist, we're already done
}

// Only returns an error if the current working directory could not be changed to the "parent" upon exit
func renameFiles(fileList []os.FileInfo) error {
	for _, fi := range fileList {
		if fi.IsDir() {
			if string(fi.Name()[0]) == "." {
				// skip this, it's a hidden directory
				continue
			}
			if recursive {
				files, err := ioutil.ReadDir(fi.Name())
				if err != nil {
					log.Printf("Could not list files in %s - %v", fi.Name(), err)
				} else {
					err := os.Chdir(fi.Name())
					if err != nil {
						log.Printf("Could not change working directory to %s, skipping it - %v", fi.Name(), err)
						return nil
					} else {
						err := renameFiles(files)
						if err != nil {
							log.Fatalf("Could not change to a parent directory we've already been to, failing - %v", err)
						}
					}
				}
			}
			continue
		}
		orig := fi.Name()
		camelName := strings.Replace(orig, string(orig[0]), strings.ToUpper(string(orig[0])), 1)
		for x := 1; x < len(camelName); x++ { // start at one since first char is uppercased already
			if len(camelName) > len(orig)*2 {
				log.Printf("Working on char[%d] %s in %s", x, string(camelName[x]), camelName)
				log.Fatalf("Something went wrong, do not continue on file %s", orig)
			}
			if strings.ContainsAny(string(camelName[x]), "ABCDEFGHIJKLMNOPQRSTUVWXYZ") && string(camelName[x-1]) != " " {
				camelName = strings.Replace(camelName, string(camelName[x]), " "+strings.ToUpper(string(camelName[x])), 1)
			}
		}
		if camelName != orig {
			if !pretend {
				err := os.Rename(orig, camelName)
				if err != nil {
					log.Printf("Error renaming %s to %s - %v", orig, camelName, err)
				}
				log.Printf("Renamed %s to %s", orig, camelName)
			} else {
				log.Printf("Would have renamed %s to %s", orig, camelName)
			}
		} else if verbose {
			log.Printf("Did not have to rename %s", orig)
		}
	}
	return os.Chdir("..")
}
