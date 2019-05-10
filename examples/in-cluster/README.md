# Running clustereye in cluster

This directory contains Kubernetes manifests for running the dashboard in a cluster.

Steps:

* Create a secret that contains your kubeconfig

```sh
 kubectl create secret generic clustereye-kubeconfig --from-file=/path/to/kubeconfig
```

* Update deployment manifest to point to secret

```yaml
# ====== 8< ======
args: ["-v", "--kubeconfig", "/kube/<file name>"]
# ====== 8< ======
```

* Apply deployment

```sh
kubectl kustomize | kubectl apply -f -
```

* Create a port forward to the pod in clustereye or using kubectl
