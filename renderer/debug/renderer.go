package debug

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	"github.com/lucasepe/yamldo/parser"
	"github.com/lucasepe/yamldo/renderer"
	"github.com/lucasepe/yamldo/table"
)

func New(showIndent bool) renderer.Renderer {
	return &debugRenderer{
		showIndent: showIndent,
	}
}

type line struct {
	nr      int
	ws      int
	content string
	path    string
}

type debugRenderer struct {
	showIndent bool
}

func (r *debugRenderer) Render(frags []parser.Fragment) ([]byte, error) {
	res := bytes.NewBufferString("")

	counter := 0
	for _, el := range frags {

		lines := computeLines(&el, counter)
		counter = counter + len(lines)

		//fmt.Fprintf(res, "%s\n", el.Path())

		tbl := &table.TextTable{}
		if el.IsKey() {
			tbl.SetHeader("ln", "ws", fmt.Sprintf("\U0001f4c2 %s", el.Path()))
		} else {
			tbl.SetHeader("ln", "ws", fmt.Sprintf("\U0001f4dd %s", el.Path()))
		}

		for _, ln := range lines {
			row := []string{
				fmt.Sprintf("%3d", ln.nr),
				fmt.Sprintf("%2d", ln.ws),
				highlightLeftSpaces(ln.content, '·'),
			}

			tbl.AddRow(row...)
		}

		tbl.AddRowLine()

		err := tbl.DrawInBuffer(res)
		if err != nil {
			return nil, err
		}

		fmt.Fprintln(res)
	}

	return res.Bytes(), nil
}

func leadingSpaces(line string, space rune) int {
	count := 0
	for _, v := range line {
		if v == space {
			count++
		} else {
			break
		}
	}

	return count
}

func highlightLeftSpaces(ln string, ch rune) string {
	var sb strings.Builder
	var i int
	var r rune
	for i, r = range ln {
		if r != ' ' {
			break
		}

		sb.WriteRune(ch)
	}
	sb.WriteString(ln[i:])

	return sb.String()
}

func computeLines(frag *parser.Fragment, counter int) []line {
	res := []line{}
	input := frag.String()

	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		counter = counter + 1
		txt := scanner.Text()

		res = append(res, line{
			nr:      counter,
			ws:      leadingSpaces(txt, ' '),
			content: txt,
			path:    frag.Path(),
		})
	}

	return res
}
