# Kafka on Kubernetes for local development

This document explains running single node local kafka infra on kubernetes it covers kafka, zookeeper and
schema-registry in kubernetes using [Kind](https://kind.sigs.k8s.io/).


1. storage-class-setup folder has the setup running with the default Kind Storage Class and Stateful sets.
2. pv-pvc-setup folder has the setup running with Persistent Volumes and Persistent Volume Claims explicitly.
3. helm foler has the setup using Helm, it uses the same Kubernetes approach as 1. storage-class-setup using Helm 
   instead of plain Kubernetes which brings flexibility. The Helm setup is kept simple.


### Pre-reqs, install:

- [Docker](https://docs.docker.com/get-docker/)
- [Kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
- [Kind](https://kind.sigs.k8s.io/)
### Step-by-step instructions to run storage-class-setup

1. Open a terminal and cd to the storage-class-setup folder.
2. Run kind specifying configuration: `kind create cluster --config=kind-config.yaml`. This will start a kubernetes
   control plane + worker
4. Run kubernetes configuration for kafka `kubectl apply -f kafka-k8s`
5. When done stop kubernetes objects: `kubectl delete -f kafka-k8s` and then if you want also stop the kind cluster
   which:
   will also delete the storage on the host machine: `kind delete cluster`

> After running `kind create...` command on step 3 above, if you have the images already downloaded in your local docker
> registry you should load them on Kind so it won't try to download the images every time, to do that use: 
> `kind load docker-image $image_name` i.e - `kind load docker-image confluentinc/cp-kafka:7.0.1`
> You can check the loaded messages entering one of the Kind docker containers(worker or control plane) and use `crictl images`

