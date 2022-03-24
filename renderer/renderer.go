package renderer

import (
	"bytes"

	"github.com/lucasepe/yamldo/parser"
)

type Renderer interface {
	Render([]parser.Fragment) ([]byte, error)
}

func New() Renderer {
	return &plainRenderer{}
}

type plainRenderer struct{}

func (r *plainRenderer) Render(items []parser.Fragment) ([]byte, error) {
	res := bytes.NewBufferString("")

	for _, el := range items {
		_, err := res.WriteString(el.String())
		if err != nil {
			return nil, err
		}

		_, err = res.WriteString("\n")
		if err != nil {
			return nil, err
		}
	}

	return res.Bytes(), nil
}
