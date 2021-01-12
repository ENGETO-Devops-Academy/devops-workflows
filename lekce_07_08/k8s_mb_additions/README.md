# K8s part 2.

Let's add autocomplete, see https://kubernetes.io/docs/tasks/tools/install-kubectl/ 

```
echo 'source <(kubectl completion bash)' >>~/.bashrc
```

### Volumes

#### Simple volumes

On master:
```
apt update
apt install -y nfs-kernel-server
mkdir /var/nfs/general -p

echo "/var/nfs/general    *(rw,sync,no_subtree_check)" > /etc/exports
systemctl restart nfs-kernel-server
```

On nodes:
```
apt-get install -y nfs-common
```

Let's create pod:

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: pod-using-nfs
spec:
  volumes:
    - name: nfs-volume
      nfs: 
        server: 10.0.1.160
        path: /var/nfs/general
  containers:
    - name: app
      image: alpine
      volumeMounts:
        - name: nfs-volume
          mountPath: /var/nfs
      command: ["/bin/sh"]
      args: ["-c", "while true; do date >> /var/nfs/dates.txt; sleep 5; done"]
```

Details at: https://matthewpalmer.net/kubernetes-app-developer/articles/kubernetes-volumes-example-nfs-persistent-volume.html

Let's exec `kubectl exec -ti pod-using-nfs -- sh`.

I case of failure, change node to correct IP:

```
cat /var/lib/kubelet/kubeadm-flags.env
KUBELET_KUBEADM_ARGS="--network-plugin=cni --pod-infra-container-image=k8s.gcr.io/pause:3.2 --node-ip=10.0.1.161"
```

Let's check the data from pod:

```
mount
cat /var/nfs/dates.txt
```

and from master:

```
cat /var/nfs/general/dates.txt
```

#### Persistent volumes and volume claims

See https://docs.openshift.com/enterprise/3.1/install_config/storage_examples/shared_storage.html

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs-pv 
spec:
  capacity:
    storage: 1Gi 
  accessModes:
    - ReadWriteMany 
  persistentVolumeReclaimPolicy: Retain 
  nfs: 
    path: /var/nfs/general
    server: 10.0.1.160
    readOnly: false
```

See:

```
kubectl get pv
NAME     CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE
nfs-pv   1Gi        RWX            Retain           Available                                   5s

```

Let's get claim:

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs-pvc  
spec:
  accessModes:
  - ReadWriteMany      
  resources:
     requests:
       storage: 1Gi
```

Why? We do not need to know about underlaying protocol.

Let's create pod:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nfs-pod 
  labels:
    name: nfs-pod
spec:
  containers:
    - name: app
      image: alpine
      volumeMounts:
        - name: nfs-volume
          mountPath: /var/nfs
      command: ["/bin/sh"]
      args: ["-c", "while true; do date >> /var/nfs/dates.txt; sleep 5; done"]
  volumes:
    - name: nfs-volume
      persistentVolumeClaim:
        claimName: nfs-pvc
```

Let's again see:

```
kubectl exec -ti nfs-pod -- sh
cat /var/nfs/dates.txt
```

See pod description:

```
kubectl describe pod nfs-pod
...
Volumes:
  nfs-volume:
    Type:       PersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)
    ClaimName:  nfs-pvc
    ReadOnly:   false
...
```

### ConfigMaps

Have configuration files and environment variables.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-configuration
data:
  variable: "and content of the variable" 
  config.txt: |
    abc=def
    something
    testtest
```

See

```
kubectl get ConfigMap app-configuration -o jsonpath='{.data}'
```

Let's create a pod which contains a variable and a file from configmap:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: configmap-pod
spec:
  containers:
    - name: app
      image: alpine
      command: [ "/bin/sh", "-c", "sleep", 6000 ]
      volumeMounts:
        - name: config-volume
          mountPath: /etc/config
      env:
        - name: SPECIAL_LEVEL_KEY
          valueFrom:
            configMapKeyRef:
              name: app-configuration
              key: variable
  volumes:
    - name: config-volume
      configMap:
        name: app-configuration
  restartPolicy: Never
```

Now see what's in the container:

```
kubectl exec -ti configmap-pod -- sh
/ # set
...
SPECIAL_LEVEL_KEY='and content of the variable'
...
/ # cat /etc/config/config.txt 
abc=def
something
testtest
```

See https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/

### Secrets

Similar to configmaps but keeps secrets ENCODED into base64.

Let's create a secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secret
type: Opaque
data:
  secret: dmVyeV9zZWNyZXRfY29udGVudAo=
  secret_file.txt: YmxhYmxhYmxhCg==
```

We can see the secrets are plain visible:

```yaml
kubectl get secrets app-secret -o yaml
apiVersion: v1
data:
  secret: dmVyeV9zZWNyZXRfY29udGVudAo=
  secret_file.txt: YmxhYmxhYmxhCg==
kind: Secret
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"secret":"dmVyeV9zZWNyZXRfY29udGVudAo=","secret_file.txt":"YmxhYmxhYmxhCg=="},"kind":"Secret","metadata":{"annotations":{},"name":"app-secret","namespace":"default"},"type":"Opaque"}
  creationTimestamp: "2021-01-09T18:21:01Z"
  managedFields:
  - apiVersion: v1
    fieldsType: FieldsV1
    fieldsV1:
      f:data:
        .: {}
        f:secret: {}
        f:secret_file.txt: {}
      f:metadata:
        f:annotations:
          .: {}
          f:kubectl.kubernetes.io/last-applied-configuration: {}
      f:type: {}
    manager: kubectl-client-side-apply
    operation: Update
    time: "2021-01-09T18:21:01Z"
  name: app-secret
  namespace: default
  resourceVersion: "41033"
  selfLink: /api/v1/namespaces/default/secrets/app-secret
  uid: 557bde8c-6881-46fc-be99-dcd9d9f76966
type: Opaque
```

Let's create a pod accessing the secrets:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: secret-pod
spec:
  containers:
    - name: app
      image: alpine
      command: [ "/bin/sh", "-c", "sleep", 6000 ]
      volumeMounts:
        - name: app-secret
          mountPath: /etc/secrets
      env:
        - name: SECRET
          valueFrom:
            secretKeyRef:
              name: app-secret
              key: secret
  volumes:
    - name: app-secret
      secret:
        secretName: app-secret
  restartPolicy: Never
```

See https://kubernetes.io/docs/concepts/configuration/secret/

#### SealedSecrets

See: https://github.com/bitnami-labs/sealed-secrets

## Jobs & CronJobs

One time thing. Mainly for computation.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
```

CronJob:

```yaml
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "*/1 * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox
            imagePullPolicy: IfNotPresent
            args:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure
```

See https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/

## Advanced topics

### Readyness vs liveliness

Let's deploy deployment with service:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: test-pi-server
  labels:
    app: test-pi-server
spec:
  replicas: 1
  selector:
    matchLabels:
      app: test-pi-server
  template:
    metadata:
      labels:
        app: test-pi-server
    spec:
      containers:
      - name: test-pi-server
        image: beranm14/test_pi_server
        ports:
        - containerPort: 3000
---
apiVersion: v1
kind: Service
metadata:
  name: test-pi-server
spec:
  selector:
    app: test-pi-server
  ports:
    - protocol: TCP
      port: 80
      targetPort: 3000
```

Now let's see if it's running:

```
kubectl run -ti --image=curlimages/curl curl --rm -- http://test-pi-server/pi
kubectl run -ti --image=curlimages/curl curl --rm -- http://test-pi-server/healthz
```

Inspect the top (heapster needed):

```
watch kubectl top pods --containers
```

### QOS classes vs reservations and allocations



See: https://kubernetes.io/docs/tasks/configure-pod-container/quality-service-pod/

```yaml
spec:
  containers:
    ...
    resources:
      limits:
        cpu: 700m
        memory: 200Mi
      requests:
        cpu: 700m
        memory: 200Mi
    ...
status:
  qosClass: Guaranteed
```

Read https://blog.pipetail.io/posts/2020-05-04-most-common-mistakes-k8s/

### HPA

See https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/

```
kubectl autoscale deployment test-pi-server --cpu-percent=50 --min=1 --max=10
```

Let's benchmark it again.

### VPA

Scaling verticaly https://cloud.google.com/kubernetes-engine/docs/how-to/vertical-pod-autoscaling

It's more & less about statistics.

## Ingress

The example is in GCP cluster: https://cloud.google.com/kubernetes-engine/docs/concepts/ingress

Let's create:

```
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: pi-ingress
spec:
  rules:
  - http:
      paths:
      - path: /pi/*
        backend:
          serviceName: my-discounted-products
          servicePort: 80
```

And inspect how it goes in GCP console.

Other ingress controllers could be nginx, haproxy, traefik, ...

See https://kubernetes.io/docs/concepts/services-networking/ingress-controllers/

## Operators & CRDs

CRD: https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/

Operators: https://kubernetes.io/docs/concepts/extend-kubernetes/operator/

Let's go to K8s patterns https://www.redhat.com/cms/managed-files/cm-oreilly-kubernetes-patterns-ebook-f19824-201910-en.pdf

## Manifests management

Helm & Kustomize

For helm, let's inspect CICD lecture at `k8s_from_cicd` folder, install command is:

```
helm upgrade --install --set="appName=deploy-0" deploy-0 ./app-chart
```
