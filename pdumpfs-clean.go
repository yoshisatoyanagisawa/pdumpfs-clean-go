package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	keep             = flag.String("keep", "2Y6M6W6D", "use KEEPARGS to decide delete directories")
	verbose          = flag.Bool("verbose", false, "verbose output")
	dryrun           = flag.Bool("dryrun", false, "dryrun")
	remove_empty_dir = flag.Bool("remove-empty", false, "remove empty directory")
)

type keepDuration struct {
	year  int
	month int
	week  int
	day   int
}

func parseKeep(input string) (keepDuration, error) {
	prev := 0
	ret := make(map[string]int)
	for _, c := range []string{"Y", "M", "W", "D"} {
		if prev >= len(input) {
			break
		}
		idx := strings.Index(input[prev:], c) + prev
		if idx != -1 {
			parsed, err := strconv.Atoi(input[prev:idx])
			if err != nil {
				return keepDuration{}, err
			}
			ret[c] = parsed
		}
		prev = idx + 1
	}
	return keepDuration{
		year:  ret["Y"],
		month: ret["M"],
		week:  ret["W"],
		day:   ret["D"],
	}, nil
}

func createKeepDirs(k keepDuration) map[string]bool {
	now := time.Now()
	ret := make(map[string]bool)
	for i := 0; i < k.year; i++ {
		ret[now.AddDate(-i, 0, 0).Format("2006")+"/01/01"] = true
	}
	for i := 0; i < k.month; i++ {
		ret[now.AddDate(0, -i, 0).Format("2006/01")+"/01"] = true
	}
	for i := 0; i < k.week; i++ {
		weekday := int(now.Weekday())
		ret[now.AddDate(0, 0, -7*i-weekday).Format("2006/01/02")] = true
	}
	for i := 0; i < k.day; i++ {
		ret[now.AddDate(0, 0, -i).Format("2006/01/02")] = true
	}
	return ret
}

func isEmptyDir(d string) (bool, error) {
	info, err := os.Stat(d)
	if err != nil {
		return false, err
	}
	if !info.IsDir() {
		return false, fmt.Errorf("%v is not directory", d)
	}
	f, err := os.Open(d)
	defer f.Close()
	if err != nil {
		return false, err
	}
	names, err := f.Readdirnames(-1)
	if err != nil {
		return false, err
	}
	return len(names) == 0, nil
}

func deleteEmptyDir(d string) error {
	empty, err := isEmptyDir(d)
	if empty {
		err = os.Remove(d)
	}
	return err
}

func main() {
	flag.Parse()
	k, err := parseKeep(*keep)
	if err != nil {
		panic(err)
	}
	keepDirs := createKeepDirs(k)
	var kept []string
	for _, basedir := range flag.Args() {
		m, err := filepath.Glob(basedir + "/[0-9][0-9][0-9][0-9]/[0-1][0-9]/[0-3][0-9]")
		if err != nil {
			panic(err)
		}
		for _, d := range m {
			r, err := filepath.Rel(basedir, d)
			if err != nil {
				panic(err)
			}
			if keepDirs[r] {
				kept = append(kept, d)
				continue
			}
			if *verbose {
				fmt.Printf("Deleting %q ...", d)
			}
			if !*dryrun {
				err = os.RemoveAll(d)
				if err != nil {
					fmt.Printf("failed to remove %q: %v\n", d, err)
				}
			}
			if *verbose {
				fmt.Println(" done.")
			}
		}
		if *remove_empty_dir {
			m, err = filepath.Glob(basedir + "/[0-9][0-9][0-9][0-9]/[0-1][0-9]")
			if err != nil {
				panic(err)
			}
			for _, d := range m {
				err = deleteEmptyDir(d)
				if err != nil {
					fmt.Printf("failed to remove %q: %v\n", d, err)
				}
			}
			m, err = filepath.Glob(basedir + "/[0-9][0-9][0-9][0-9]")
			if err != nil {
				panic(err)
			}
			for _, d := range m {
				err = deleteEmptyDir(d)
				if err != nil {
					fmt.Printf("failed to remove %q: %v\n", d, err)
				}
			}
		}
	}
}
