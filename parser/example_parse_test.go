package parser

import (
	"fmt"
)

func ExampleParse() {
	path := "../testdata/deployment"
	res, err := Parse(path)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)

	// Output:
	// [{header.yaml false 0 apiVersion: apps/v1
	// kind: Deployment} {metadata.yaml false 0 metadata:
	//   name: nginx-deployment
	//   labels:
	//     app: nginx} {spec true 0 spec} {spec/replicas.yaml false 1 replicas: 3} {spec/selector.yaml false 1 selector:
	//   matchLabels:
	//     app: nginx} {spec/template true 1 template} {spec/template/metadata.yaml false 2 metadata:
	//   labels:
	//     app: nginx} {spec/template/spec true 2 spec} {spec/template/spec/containers true 3 containers} {spec/template/spec/containers/nginx.yaml false 4 - name: nginx
	//   image: nginx:1.21-alpine
	//   ports:
	//   - containerPort: 80}]
}
