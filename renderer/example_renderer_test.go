package renderer

import (
	"fmt"

	"github.com/lucasepe/yamldo/parser"
)

func ExampleRenderer_plain() {
	path := "../testdata/deployment"
	blocks, err := parser.Parse(path)
	if err != nil {
		panic(err)
	}

	res, err := New().Render(blocks)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", res)

	// Output:
	// apiVersion: apps/v1
	// kind: Deployment
	// metadata:
	//   name: nginx-deployment
	//   labels:
	//     app: nginx
	// spec:
	//   replicas: 3
	//   selector:
	//     matchLabels:
	//       app: nginx
	//   template:
	//     metadata:
	//       labels:
	//         app: nginx
	//     spec:
	//       containers:
	//         - name: nginx
	//           image: nginx:1.21-alpine
	//           ports:
	//           - containerPort: 80
}
