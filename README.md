# Custom Openshift Payload CICD Pipeline

## Intro

This is a simple POC project that executes on a local machine (desktop) that has a dependency on podman, a local docker registry and golang (ver 1.19.3).
The intention is to build all the necessary core operators for openshift/okd without the need of kind, tekton etc, it makes use of simple yaml files and golang structures 
to build and push these operators for later use.

The build time are reduced largely due to the local golang build cache and podman build cache. This was the main driver for opting to build this POC, as Tekton needs some extra thought implementing 
caching in the PV allocation

This was succesfully test with golang 1.19.3, podman 4.2.1 and docker registry v2.7.0-2013-gcec2cad8.m

**NB** This is a WIP 

## Description

This custom cicd uses concepts (similar to Tekton) i.e pipeline, tasks and taskruns to compile the golang openshift core operators and then push them to a local registry

### Clone the repository and build

```bash
git clone git@github.com:luigizuccarelli/custom-openshiftpayload-cicd

cd custome-openshiftpayload-cicd
make clean
make build

```

## Usage

Generate the relevant taskrun objects (the generated files will be stored in the folder manifests/taskruns)

Execute the following command


```bash
# this will use the repositories text file to autogenerate the taskruns fro the project
./build/cicd -g <folder-to-buildconfigs> -d <destination-folder-for taskruns> 

```

Execute the following to start a pipeline 

```bash
# NB a directory will be created from the field spec.workspaces[0].name in the Pipeline (file reference in manifests/pipeline)
# The command below will execute all the taskruns in the folder that points to the main kustomization folder
./build/cicd -k environments/overlays/cicd
```

