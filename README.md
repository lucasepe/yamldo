# yamldo

Create YAML documents from a directory tree or a ZIP archive.

---

## Description 

Build a YAML document from a directory tree or a ZIP archive containing YAML fragments.

It is useful everywhere complex YAML configuration is employed: 

- CI pipelines 
- Crossplane Compositions 
- Kubernetes manifests, etc.

For instance, Kubernetes resources can be complex, and in turn make the YAML very verbose.

The aim is to make the YAML document easier to reason about and maintain.

## How to use

`yamldo` parses a directory tree:

- a directory is a yaml key
- a file (a yaml fragment) has content that is rendered at the current indentation level

Consider the following [example](./testdata/deployment):

```sh
$ tree testdata/deployment 
testdata/deployment
├── header.yaml
├── metadata.yaml
└── spec
    ├── replicas.yaml
    ├── selector.yaml
    └── template
        ├── metadata.yaml
        └── spec
            └── containers
                └── nginx.yaml

4 directories, 6 files
```

Running `yamldo` with the following command:

```sh
$ yamldo testdata/deployment
```

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.21-alpine
          ports:
          - containerPort: 80
```

You can also work on a tree in a [ZIP archive](./testdata/deployment.zip):

```sh
$ yamldo testdata/deployment.zip
```

And eventually you can pipe the result straight to `kubectl`.

### Working with linter

`yamldo` can work side by side with the linter [yamllint](https://yamllint.readthedocs.io/en/stable/quickstart.html):

```sh
$ yamldo testdata/deployment | yamllint -
stdin
  1:1       warning  missing document start "---"  (document-start)
  21:11     error    wrong indentation: expected 12 but found 10  (indentation)
```

and to check your fragments, for debugging, you can use the `-debug` flag:

```sh
$ yamldo -debug testdata/deployment
+-----+----+---------------------+
| ln  | ws | 📝 header.yaml      |
+-----+----+---------------------+
|   1 |  0 | apiVersion: apps/v1 |
|   2 |  0 | kind: Deployment    |
+-----+----+---------------------+

+-----+----+--------------------------+
| ln  | ws | 📝 metadata.yaml         |
+-----+----+--------------------------+
|   3 |  0 | metadata:                |
|   4 |  2 | ··name: nginx-deployment |
|   5 |  2 | ··labels:                |
|   6 |  4 | ····app: nginx           |
+-----+----+--------------------------+

+-----+----+---------+
| ln  | ws | 📂 spec |
+-----+----+---------+
|   7 |  0 | spec:   |
+-----+----+---------+

+-----+----+-----------------------+
| ln  | ws | 📝 spec/replicas.yaml |
+-----+----+-----------------------+
|   8 |  2 | ··replicas: 3         |
+-----+----+-----------------------+

+-----+----+-----------------------+
| ln  | ws | 📝 spec/selector.yaml |
+-----+----+-----------------------+
|   9 |  2 | ··selector:           |
|  10 |  4 | ····matchLabels:      |
|  11 |  6 | ······app: nginx      |
+-----+----+-----------------------+

+-----+----+------------------+
| ln  | ws | 📂 spec/template |
+-----+----+------------------+
|  12 |  2 | ··template:      |
+-----+----+------------------+

+-----+----+--------------------------------+
| ln  | ws | 📝 spec/template/metadata.yaml |
+-----+----+--------------------------------+
|  13 |  4 | ····metadata:                  |
|  14 |  6 | ······labels:                  |
|  15 |  8 | ········app: nginx             |
+-----+----+--------------------------------+

+-----+----+-----------------------+
| ln  | ws | 📂 spec/template/spec |
+-----+----+-----------------------+
|  16 |  4 | ····spec:             |
+-----+----+-----------------------+

+-----+----+----------------------------------+
| ln  | ws | 📂 spec/template/spec/containers |
+-----+----+----------------------------------+
|  17 |  6 | ······containers:                |
+-----+----+----------------------------------+

+-----+----+---------------------------------------------+
| ln  | ws | 📝 spec/template/spec/containers/nginx.yaml |
+-----+----+---------------------------------------------+
|  18 |  8 | ········- name: nginx                       |
|  19 | 10 | ··········image: nginx:1.21-alpine          |
|  20 | 10 | ··········ports:                            |
|  21 | 10 | ··········- containerPort: 80               |
+-----+----+---------------------------------------------+
```

The linter (obviously) was right! the fragment _spec/template/spec/containers/nginx.yaml_ with the line _ln:21_ has only 10 spaces.

---

### Reuse and maintenance of your fragments

This tool comes in handy with really complex YAML documents where you can split reusable fragments.

Check the [./testdata/platform-ref-aws/cluster/eks](./testdata/platform-ref-aws/cluster/eks/) for a more complex YAML (a Crossplane composition).

```sh
$ tree testdata/platform-ref-aws/cluster/eks
testdata/platform-ref-aws/cluster/eks
├── header.yaml
├── metadata.yaml
└── spec
    ├── header.yaml
    └── resources
        ├── cluster.yaml
        ├── nodegroup.yaml
        ├── oidc_provider.yaml
        ├── provider_config.yaml
        ├── role_controlplane.yaml
        ├── role_nodegroup.yaml
        ├── role_policy_attachment_controlplane.yaml
        └── role_policy_attachment_nodegroup.yaml
```

## How to install

Binary downloads of `yamldo` can be found on the [Releases page](https://github.com/lucasepe/yamldo/releases).

Unpack the binary and add it to your PATH and you are good to go!

### MacOS

```sh
$ brew update
$ brew tap lucasepe/yamldo
$ brew install yamldo
```

or if you have already installed `yamldo` using brew, you can upgrade it by running:

```sh
$ brew upgrade yamldo
```

---

#### Credits

Credits to [sampointer/dy](https://github.com/sampointer/dy) for the basic idea. 

`yamldo` works a little differently but the inspiration came from that project.

