## Simple Kubernetes Kubeconfig Cluster Selector

![Screenshot](docs/ks.png)

### Requirements:

- go 1.25

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
rancherFix: false
```

### Fix:
For mitigating these issues:
- https://github.com/rancher/rancher/issues/55031
- https://github.com/rancher/rancher/issues/55034

set configuration varialble `rancherFix` to `true`


### Note:

If the config option `createLink` is set to `true` and the kubeconfig file `~/.kube/config` exists or it's set to `false` and the config file is not managed by ks, then you need to remove the kubeconfig file first. This is a safety measure to not override existing configurations and this file needs to be removed or rename by hand.
