package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type includeFlags []string

func (i *includeFlags) String() string {
	return "includeFlags"
}

func (i *includeFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var includes includeFlags

type excludeFlags []string

func (i *excludeFlags) String() string {
	return "excludeFlags"
}

func (i *excludeFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var excludes excludeFlags

var (
	only = flag.String("f", "", "file path(optional)")
	// include = flag.SliceString("ex", "", "regexp to exclude file name(optional)")
	// exclude = flag.String("ex", "", "regexp to exclude file name(optional)")
)

type TotalStat struct {
	Bytes int
	Lines int
	Files int
}

type FileStat struct {
	Bytes int
	Lines int
}

func main() {
	flag.Var(&includes, "in", "regexp to include file name(optional)")
	flag.Var(&excludes, "ex", "regexp to exclude file name(optional)")
	flag.Parse()

	fs := make([]string, 0)
	if *only != "" {
		fs = []string{*only}
	} else {
		inre := make([]*regexp.Regexp, 0)
		exre := make([]*regexp.Regexp, 0)
		if len(excludes) > 0 {
			for _, v := range excludes {
				exre = append(exre, regexp.MustCompile(v))
			}
		}
		if len(includes) > 0 {
			for _, v := range includes {
				inre = append(inre, regexp.MustCompile(v))
			}
		}
		fs = flattenDir("./", exre, inre)
	}

	var lock sync.Mutex

	ffsMap := make(map[string]FileStat, len(fs))
	var wg sync.WaitGroup
	for _, v := range fs {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			fs, err := statFile(path)
			if err != nil {
				fmt.Println(err)
				return
			}
			lock.Lock()
			ffsMap[path] = fs
			lock.Unlock()
		}(v)
	}
	wg.Wait()

	ts := TotalStat{}
	for _, v := range ffsMap {
		ts.Bytes += v.Bytes
		ts.Lines += v.Lines
		ts.Files++
	}

	fmt.Fprintf(os.Stdout, "TotalStats\n")
	fmt.Fprintf(os.Stdout, totalOut(ts))
	fmt.Fprintf(os.Stdout, "Each files\n")
	for k, v := range ffsMap {
		fmt.Fprintf(os.Stdout, statFileOut(k, v))
	}
}

func totalOut(ts TotalStat) string {
	var header, content bytes.Buffer
	header.WriteString("| ")
	content.WriteString("| ")
	fh := "Files "
	fc := fmt.Sprintf("%d ", ts.Files)
	if len(fh) > len(fc) {
		fc += strings.Repeat(" ", len(fh)-len(fc))
	} else if len(fh) < len(fc) {
		fh += strings.Repeat(" ", len(fc)-len(fh))
	}
	header.WriteString(fh + "| ")
	content.WriteString(fc + "| ")

	lh := "Lines "
	lc := fmt.Sprintf("%d ", ts.Lines)
	if len(lh) > len(lc) {
		lc += strings.Repeat(" ", len(lh)-len(lc))
	} else if len(lh) < len(lc) {
		lh += strings.Repeat(" ", len(lc)-len(lh))
	}
	header.WriteString(lh + "| ")
	content.WriteString(lc + "| ")

	bh := "Bytes "
	bc := fmt.Sprintf("%d ", ts.Bytes)
	if len(bh) > len(bc) {
		bc += strings.Repeat(" ", len(bh)-len(bc))
	} else if len(bh) < len(bc) {
		bh += strings.Repeat(" ", len(bc)-len(bh))
	}
	header.WriteString(bh + "|")
	content.WriteString(bc + "|")
	headerStr := header.String()
	contentStr := content.String()
	l := strings.Repeat("-", len(headerStr))
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", l, headerStr, l, contentStr, l)
}

func statFileOut(fileName string, fs FileStat) string {
	content := bytes.NewBuffer([]byte(fmt.Sprintf("| %d | %d ", fs.Lines, fs.Bytes)))
	fname := bytes.NewBuffer([]byte(fmt.Sprintf("| %s ", fileName)))
	if fname.Len() > content.Len() {
		content.WriteString(strings.Repeat(" ", fname.Len()-content.Len()))
	} else if fname.Len() < content.Len() {
		fname.WriteString(strings.Repeat(" ", content.Len()-fname.Len()))
	}
	content.WriteString("|")
	fname.WriteString("|")
	l := strings.Repeat("-", content.Len())
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n", l, fname.String(), l, content.String(), l)
}

func statFile(path string) (FileStat, error) {
	f, err := os.Open(path)
	if err != nil {
		return FileStat{}, err
	}
	defer f.Close()

	fs := FileStat{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		b := scanner.Bytes()
		fs.Bytes += len(b)
		fs.Lines++
	}
	return fs, nil
}

func flattenDir(dir string, exre []*regexp.Regexp, inre []*regexp.Regexp) []string {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var paths []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}
		isEx := false
		for _, re := range exre {
			if re.MatchString(file.Name()) {
				isEx = true
				break
			}
		}
		if isEx {
			continue
		}
		if file.IsDir() {
			paths = append(paths, flattenDir(filepath.Join(dir, file.Name()), exre, inre)...)
			continue
		}
		if len(inre) > 0 {
			for _, re := range inre {
				if re.MatchString(file.Name()) {
					paths = append(paths, filepath.Join(dir, file.Name()))
				}
			}
			continue
		}
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths
}
