/*
Copyright © 2024 marcelo-fm

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gocolly/colly"
	"github.com/marcelo-fm/arcpy2go/arcpy-scraper/web"
)

const defaultPackageName = "arcpy"

func main() {
	packageMode := flag.Bool("package", false, "saves the file in a package folder named after the tool")
	packageName := flag.String("package-name", defaultPackageName, "name of the package of the generated Go file")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags] <url> [path]\n\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "Flags:")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 || len(args) > 2 {
		flag.Usage()
		os.Exit(2)
	}

	url := args[0]
	file := os.Stdout

	homeConfigDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	appConfigDir := filepath.Join(homeConfigDir, "arcpy2go")
	cacheDir := filepath.Join(appConfigDir, "cache")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	c := colly.NewCollector(
		colly.CacheDir(cacheDir),
	)

	data, err := web.Parse(c, url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	data.PackageName = *packageName

	if len(args) == 2 {
		path := args[1]
		filePath := path
		if *packageMode {
			if err := os.MkdirAll(path, 0755); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			filePath = filepath.Join(path, fmt.Sprintf("%s.go", data.FunctionName))
		} else {
			dirPath := filepath.Dir(path)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}

		file, err = os.Create(filePath)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer file.Close()
	}

	if err := data.Render(file); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
