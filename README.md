## Simple Kubernetes Kubeconfig Cluster Selector

![Screenshot](docs/ks.png)

### Requirements:

- go 1.22

### Usage:

```
NAME:
   cluster - Select kubeconfig

USAGE:
   ks [global options] command [command options] [arguments...]

VERSION:
   1.2 (bf8ed0e)

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --debug        Enable debug
   --help, -h     show help
   --version, -v  print the version
```

### Build:

`make`

### Install:

`make install`

or download binary in [releases](https://github.com/bvankampen/kubeconfig-selector/releases).

### Config:

`~/.config/ks.conf`

```
kubeconfigDir: ~/.kube
kubeconfigFile: config
extraKubeconfigDirs:
    - ~/Downloads
showKubeconfig: true
createLink: true
```
