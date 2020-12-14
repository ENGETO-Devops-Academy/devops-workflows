#### Teached footnote, if you are interested in K8s, wait for next lecture, deployment is in folder k8s, this are just short notes for me to make k8s work

First we have to setup cluster credentials in `Operations` -> `Kubernetes` tab.

For further info, see https://docs.gitlab.com/ee/user/project/clusters/add_remove_clusters.html#existing-kubernetes-cluster

Getting credentials for project (you'll be told those, no worries):
```
Get API IP: kubectl cluster-info | grep 'Kubernetes master' | awk '/http/ {print $NF}'
Get CA: kubectl get secret default-token-zqp2j -n kube-system -o jsonpath="{['data']['ca\.crt']}" | base64 --decode
Get token: kubectl -n kube-system describe secret $(kubectl -n kube-system get secret | grep gitlab | awk '{print $1}')
```

if token could not be retrieved, steal one with shorter timespan:
```
/usr/lib/google-cloud-sdk/bin/gcloud config config-helper --format=json --min-expiry=1h
```

First deployment of our applications:

```
kubectl apply -f ingress.yaml
helm upgrade --install --set="appName=deploy-0" deploy-0 ./app-chart
```
