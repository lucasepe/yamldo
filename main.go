/*
Copyright (c) 2012-2022 Luca Sepe

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
"Software"), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/lucasepe/yamldo/parser"
	"github.com/lucasepe/yamldo/renderer"
	"github.com/lucasepe/yamldo/renderer/debug"
)

const (
	banner = `                ┏┓    ┏┓ 
┏┓ ┏┓ ┏━━┓ ┏┓┏┓ ┃┃  ┏━┛┃ ┏━━┓
┃┗━┛┃ ┃┏┓┃ ┃┗┛┃ ┃┃  ┃┏┓┃ ┃┏┓┃
┗━┓┏┛ ┃┏┓┃ ┃┃┃┃ ┃┗┓ ┃┗┛┃ ┃┗┛┃  ver: VERSION
┗━━┛  ┗┛┗┛ ┗┻┻┛ ┗━┛ ┗━━┛ ┗━━┛  cid: BUILD`
)

var (
	Version string
	Build   string

	flagDebug      bool
	flagShowIndent bool
)

// no more pain to  keep the right indentation.
func main() {
	configureFlags()

	if flag.CommandLine.Arg(0) == "" {
		flag.CommandLine.Usage()
		os.Exit(1)
	}

	dir := flag.Args()[0]

	blocks, err := parser.Parse(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s while parsing '%s' directory\n", err.Error(), dir)
		os.Exit(1)
	}

	var rndr renderer.Renderer
	if flagDebug {
		rndr = debug.New(flagShowIndent)
	} else {
		rndr = renderer.New()
	}

	res, err := rndr.Render(blocks)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s while rendering '%s' directory\n", err.Error(), dir)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "%s", res)
}

func configureFlags() {
	name := appName()

	flag.CommandLine.Usage = func() {
		printBanner()
		fmt.Println("Create YAML documents from a directory tree or a ZIP archive.")
		fmt.Println(" ➽ no more pain to keep the right indentation")
		fmt.Println()

		fmt.Print("USAGE:\n\n")
		fmt.Printf("  %s [flags] <directory or zip archive with yaml fragments>\n\n", name)

		fmt.Print("EXAMPLE(s):\n\n")
		fmt.Printf("  %s /path/to/dir/with/yaml/fragments\n", name)
		fmt.Printf("  %s /path/to/archive.zip\n", name)
		fmt.Printf("  %s -debug /path/to/dir/with/yaml/fragments\n", name)
		fmt.Println()

		fmt.Print("FLAGS:\n\n")
		flag.CommandLine.SetOutput(os.Stdout)
		flag.CommandLine.PrintDefaults()
		flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
		fmt.Print("  -help\n\tprints this message\n")
		fmt.Println()
		fmt.Printf("Crafted with passion by Luca Sepe <https://github.com/lucasepe>\n\n")
	}

	flag.CommandLine.SetOutput(ioutil.Discard) // hide flag errors
	flag.CommandLine.Init(os.Args[0], flag.ExitOnError)

	flag.CommandLine.BoolVar(&flagDebug, "debug", false, "show all the yaml parts without generating the final document")

	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s while parsing args\n", err.Error())
		os.Exit(1)
	}
}

func printBanner() {
	res := strings.Replace(banner, "VERSION", Version, 1)
	res = strings.Replace(res, "BUILD", Build, 1)
	fmt.Print(res, "\n\n")
}

func appName() string {
	return filepath.Base(os.Args[0])
}
