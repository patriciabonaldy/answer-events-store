1.- create a cluster
kind create cluster --config=k8s/kind.yaml --name=patycluster

2.- change to cntext created
kubectl cluster-info --context kind-patycluster
kubectl config use-context kind-patycluster

3.- create deployment yaml file
4.- apply yam file
   kubectl apply -f k8s/deployment.yaml
