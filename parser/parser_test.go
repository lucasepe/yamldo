package parser

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDepth(t *testing.T) {
	table := []struct {
		path string
		want int
	}{
		{"header.yaml", 0},
		{"spec", 0},
		{"spec/template", 1},
		{"spec/template/spec", 2},
		{"spec/template/spec/containers", 3},
	}

	for _, el := range table {
		t.Run(el.path, func(t *testing.T) {
			got := depth(el.path)
			if got != el.want {
				t.Errorf("test fail! want: '%d', got: '%d'", el.want, got)
			}
		})
	}
}

func TestParseDir(t *testing.T) {
	want := expectedResult()

	got, err := parseDir("../testdata/deployment")
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(got, want, cmp.AllowUnexported(Fragment{}))
	if len(diff) > 0 {
		t.Fatal(diff)
	}
}

func TestParseZip(t *testing.T) {
	want := expectedResult()

	got, err := parseZip("../testdata/deployment.zip")
	if err != nil {
		t.Fatal(err)
	}

	diff := cmp.Diff(got, want, cmp.AllowUnexported(Fragment{}))
	if len(diff) > 0 {
		t.Fatal(diff)
	}
}

func expectedResult() []Fragment {
	return []Fragment{
		{path: "header.yaml", isKey: false, depth: 0, content: `apiVersion: apps/v1
kind: Deployment`},
		{path: "metadata.yaml", isKey: false, depth: 0, content: `metadata:
  name: nginx-deployment
  labels:
    app: nginx`},
		{path: "spec", isKey: true, depth: 0, content: "spec"},
		{path: "spec/replicas.yaml", isKey: false, depth: 1, content: "replicas: 3"},
		{path: "spec/selector.yaml", isKey: false, depth: 1, content: `selector:
  matchLabels:
    app: nginx`},
		{path: "spec/template", isKey: true, depth: 1, content: "template"},
		{path: "spec/template/metadata.yaml", isKey: false, depth: 2, content: `metadata:
  labels:
    app: nginx`},
		{path: "spec/template/spec", isKey: true, depth: 2, content: "spec"},
		{path: "spec/template/spec/containers", isKey: true, depth: 3, content: "containers"},
		{path: "spec/template/spec/containers/nginx.yaml", isKey: false, depth: 4, content: `- name: nginx
  image: nginx:1.21-alpine
  ports:
  - containerPort: 80`},
	}
}
