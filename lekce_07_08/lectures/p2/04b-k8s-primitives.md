
<style>{less: ../../../../main.less}</style>

<section class="navigation">
<a href="../../README.md" class="button is-small has-tooltip-bottom is-pulled-right" data-tooltip="Return to the Course Overview"><i class="fas fa-home"></i></a>
<div class="columns is-centered">
    <a href="../p1/04a-k8s-primitives.md" class="button is-small has-tooltip-bottom" data-tooltip="Previous"><i class="fas fa-arrow-left"></i></a>
    <a href="./05-k8s-multi-tenancy.md" class="button is-small has-tooltip-bottom" data-tooltip="Next"><i class="fas fa-arrow-right"></i></a>
</div>
</section>

# <a name="ch:kubernetes-primitives"></a> Kubernetes Primitives

<section class="intro-section">

<!-- Intro -->

As we already know, Kubernetes is declarative. Basically anything in Kubernetes can be described as a **manifest**. The way users usually interact with Kubernetes is by creating such manifests and submitting them via `kubectl` CLI tool. This is a convenient, declarative way to interact with the Kubernetes API. The manifests can be in JSON or YAML format. `kubectl` interacts with the control plane, or more specifically, with the **Kubernetes API server**. Kubernetes controllers then watch for differences between the *desired* state and the current *observed* state of a resource. This is called **declarative state management**. In terms of Kubernetes, the desired state is the **spec** and the observed state is the **status**.

In this chapter, we'll take a look at Kubernetes primitives, that is, the basic kinds of resoruces that we can create and manage when deploying an application to a Kubernetes cluster. There's in fact huge amount of resources and the amount is continuously increasing ever since Custom Resource Definition resources (CRDs) have been introduced &mdash; we'll talk about these in the chapter dedicated to extending Kubernetes. Hence, be prepared to do a lot of research on your own<sup class="has-tooltip-multiline" data-tooltip="Kubernetes is backed by a huge community and usually anything related to it is well documented and contains a quickstart guide."><i class="fas fa-sm fa-question"></i></sup>. The purpose of this chapter, however, is to get you started with a basic deployment workflow of an application, so that you can then just expand your knowledge while having solid knowledge about the basic concepts.

<table border="0" class="questions-table">
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    How do you interact with Kubernetes?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    How do you interact with Kubernetes?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    How do you create or delete a resource?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is <code>Deployment</code>?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is <code>Service</code>?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is <code>Ingress</code>?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is <code>PersistentVolumeClaim</code>?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is dynamic provisioning?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is the difference between <code>Deployment</code>, <code>StatefulSet</code> and <code>DaemonSet</code>?
    </span>
</td></tr>
</table>

<table class="shortcuts-table">
<tr>
    <td>CRD</td>
    <td>CustomResourceDefinition</td>
</tr>
<tr>
    <td>PVC</td>
    <td>PersistentVolumeClaim</td>
</tr>
<tr>
    <td>SC</td>
    <td>StorageClass</td>
</tr>
</table>


</section>

## <a name="sec:kubernetes-primitives/application"></a> Application

Let's first describe the application that we're going to deploy. We could probably find a hello-world application in the Kubernetes docs, but I am not a huge fan of these. Rather, let's write something simple ourselves so that we could modify it for our purposes later on the road.

There is a simple HTTP server written in go in the [app/](./app/) directory with a single root path `/`. Kubernetes is a container orchestration engine, hence it requires a containerized application. Let's create a [Dockerfile](./app/Dockerfile) for it, so that we could build an image.

Build that image and push it to a docker registry. I recommend a public one to begin with. Pulling images from a private registry a bit more complicated.

Now, running the echo server and hitting the endpoint simply with

```sh
curl -v localhost:8080
```

should give us a simple response:

```
Engeto: Kubernetes Example Application

Request received: &http.Request{Method:"GET", URL:(*url.URL)(0xc0000f2100), Proto:"HTTP/1.1", ProtoMajor:1, ProtoMinor:1, Header:http.Header{"Accept":[]string{"*/*"}, "User-Agent":[]string{"curl/7.64.1"}}, Body:http.noBody{}, GetBody:(func() (io.ReadCloser, error))(nil), ContentLength:0, TransferEncoding:[]string(nil), Close:false, Host:"localhost:8080", Form:url.Values(nil), PostForm:url.Values(nil), MultipartForm:(*multipart.Form)(nil), Trailer:http.Header(nil), RemoteAddr:"172.17.0.1:48238", RequestURI:"/", TLS:(*tls.ConnectionState)(nil), Cancel:(<-chan struct {})(nil), Response:(*http.Response)(nil), ctx:(*context.cancelCtx)(0xc0000ca3c0)}
```

Let's see how we can deploy the same application into Kubernetes.

## <a name="sec:kubernetes-primitives/deployment"></a> Kubernetes Primitives: Deployment

Let's get to the business and deploy our application. For that, we need to **create a `Deployment` resource** in the cluster.

The manifest of the `Deployment` could look like this:

```yaml
# file: deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:  # optional, but it's a good practice
    app: echoserver
  name: engeto-echoserver
  namespace: engeto  # optional, by default to the current namespace
spec:
  selector:
    matchLabels:
      app: echoserver
  template:
    metadata:
      labels:
        app: echoserver
    spec:
      containers:
      - name: main
        image: <namespace>/echoserver:latest  # make sure to provide your own docker namespace
```

This is about as simple as it gets. There is a bunch of additional stanzas we can use, but this is sufficient for the time being. So what did we do here?<br>
The resource is identified by its `apiVersion`, `kind` and `name`. It is a good practice to provide also the `app` label which is then common among multiple resources that are somewhat logically linked together. The `spec` then defines the `template` of a `Pod` which will be created. A `Pod` is a smallest deployable unit of computing that can be created and managed in Kubernetes. Each `Pod` has **at least** one container. `Pod`s are usually not created directly, but via `Deployments`, `DaemonSets` or `StatefulSets` (more on that in the <a href="#sec:kubernetes-primitives/others">last chapter</a>). A `Pod` has it's own `spec` which requires `containers` to be specified. Here, we provide the docker image with our application.

We can then submit that deployment to the cluster using the following command:

```
kubectl create -f deployment.yaml
```

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
<b>TIP:</b> If you want to perform <code>kubectl apply</code> on this object in the future, it is useful to add `--safe-config` parameter to store the current configuration so that Kubernetes could apply the patch.
</div>
</article>

So, in a nutshell, when we submit a `Deployment`, it flows like this:

<ol>
<li><code>Deployment</code> is submitted via <code>kubectl create -f deployment.yaml</code></li>
<li>API server</li>
<li>deployment-controller</li>
<li>scheduler</li>
<li><code>Pod</code> is scheduled</code> on a <code>Node</code></ol>

Let's go ahead and create the `Deployment`. A `Pod` should have been created:

```
$ kubectl get pods -n engeto
NAME                                READY   STATUS    RESTARTS   AGE
engeto-echoserver-7d9b55884-bflbh   1/1     Running   0          15s
```

Now if we get the logs from the pod, or rather the container running in that pod, we should see the message that the echo serves is being served on port `8080`:

```
$ kubectl logs engeto-echoserver-7d9b55884-bflbh 
Serving Echo Server on: http://0.0.0.0:8080
```

### <a name="subsec:kubernetes-primitives/deployment/customization"></a> Customizing a deployment

Very often we'd need to allow users of our applications to customize it to their needs. As we most likely don't want them to build their own images with custom parameters, like the `port` for example, there are two approaches how to handle that &ndash; either by using config files or environment variables.

For now, let's go with the latter. The application can read the `SERVER_HOST` and `SERVER_PORT` environment variables. The way we feed these to the `Deployment` is via `spec.template.spec.containers[].env[]`:

Try to modify that section with the following example:

```yaml
  ...
    containers:
    - name: main
      image: cermakm/echoserver:latest
      env:
      - name: SERVER_HOST
        value: "5000"
```

Check out the log again to verify that the server is running on the port `5000` now. We'll get to the second case &mdash; providing parameters via a configuration file &mdash; in the next section.

## <a name="sec:kubernetes-primitives/configmap-and-secrets"></a> Kubernetes Primitives: ConfigMaps and Secrets

`ConfigMap`s and `Secret`s are resources which store a configuration of an application.

### ConfigMaps

`ConfigMap`s are quite simple, they just store `data` as key=value pairs. An example of a `ConfigMap` which stores the `SERVER_HOST` and `SERVER_PORT` values:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: echoserver-config
  labels:
    app: echoserver
data:
  host: "0.0.0.0"
  port: "5000"
```

There are three ways you would consume a `ConfigMap` in an application:

1. Using the value of a key as an environment variable
2. Load the whole `ConfigMap` into the environment
3. Mount the `ConfigMap` as a config file into the container

ad 1)

```yaml
  ...
    env:
    - name: SERVER_PORT
      valueFrom:
      configMapKeyRef:
        name: echoserver-config
        key: port
        optional: true  # start the container even if the key is missing
```

ad 2)

```yaml
  ...
    envFrom:
    - configMapRef:
        name: echoserver-config
```

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
Note that in this case, the environment keys will be called <code>host</code> and <code>port</code>.
</div>
</article>


ad 3)

Consuming the `ConfigMap` as a volume is a little more complicated. We're going to talk about `Volume`s in a separate section, however, an example of such use case might look like this:

First, the `ConfigMap` itself would look a little bit different. Each key in the `ConfigMap` will be **mounted as a separate file** at the `mountPath`. Therefore, we would usually write something like this:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: echoserver-config
  labels:
    app: echoserver
data:
  config.json: |-
    {
      "host": "0.0.0.0",
      "port": "5000"
    }
```

And then we could mount the `ConfigMap` as follows:

```yaml
  containers:
    - name: main
      ...
      volumeMounts:
      - name: config  # references the `volumes[].name`
        mountPath: /etc/echoserver/  # will result in a file /etc/echoserver/config.json
  volumes:
    - name: config
      configMap:
        name: echoserver-config
```

The application logic would of course need to be adjusted for such change. Feel free to do so!

### Secrets

In Kubernetes, `Secret`s let you store sensitive information, such as passwords, OAuth tokens, and ssh keys. Storing confidential information in a Secret is safer and more flexible than putting it verbatim in a Pod definition or in a container image <sup>[<a href="#ref:k8s.io/concepts/secrets">0</a>]</sup>. The primary reason being that the values stored in secrets are `base64` encoded, that is not really a safety measure (base64 is considered plain text), but it does at least prevents someone from beaing able to directly read a secret password.

A `Secret` might for example look as such:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: example-secret
type: Opaque
data:
  username: YWRtaW4=
  password: cGFzc3dvcmQ=
```

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
To encode a the `username` used above into <code>base64</code> encoded string, issue <code>echo -n "admin" | base64</code>.
</div>
</article>

`Secret`s are quite similar to `ConfigMap`s and they are also consumed in a similar fashion, however, having these two configuration concepts devided has some advantages:

- Secrets are less verbatim
- Secrets have additional mechanisms that make their creation easier, i. e. `kubectl create secret` for reference.
- Logical differentiation between application configuration and secrets / credentials. This allows to control access to secrets in general as well by using native Kubernetes RBAC.

It's also necessary to list the risk that are associated with `Secrets` <sup>[<a href="#ref:k8s.io/concepts/secrets">0</a>]</sup>:

<div class="subsubsection">

- In the API server, secret data is stored in etcd ; therefore:

  - Administrators should enable encryption at rest for cluster data (requires v1.13 or later).
  - Administrators should limit access to etcd to admin users.
  - Administrators may want to wipe/shred disks used by etcd when no longer in use.

- If running etcd in a cluster, administrators should make sure to use SSL/TLS for etcd peer-to-peer communication.
- If you configure the secret through a manifest (JSON or YAML) file which has the secret data encoded as base64, sharing this file or checking it in to a source repository means the secret is compromised. Base64 encoding is not an encryption method and is considered the same as plain text.

</div>

Applications still need to protect the value of secret after reading it from the volume, such as not accidentally logging it or transmitting it to an untrusted party. **A user who can create a Pod that uses a secret can also see the value of that secret**. Even if the API server policy does not allow that user to read the `Secret`, the user could run a Pod which exposes the secret. Currently, anyone with root permission on any node can read any secret from the API server, by impersonating the kubelet. It is a planned feature to only send secrets to nodes that actually require them, to restrict the impact of a root exploit on a single node.

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
Have a look at solutions like <a href="https://www.vaultproject.io">Vault</a> from Hashicorp which allow to manage cluster secrets in a much safer way.
</div>
</article>

## <a name="sec:kubernetes-primitives/service"></a> Kubernetes Primitives: Services

Having a container running inside a Kubernetes cluster is a decent achievement, but in most cases, it's utterly useles unles we can somehow interact with it and integrate it with other applications. By default, Kubernetes pods are **only accessible from within the cluster**. Each pod is assigned a unique cluster `IP` &mdash; which is a virtual IP &mdash; which identifies it within the cluster. However, in most cases we do want to expose these applications (like servers, microservices, controllers, UIs) to the outer world.

For example, let's take a look at the `Pod` resource of our sample application:

```
kubectl get pods -n engeto
```

There are multiple ways to expose the containers (or pods, to follow Kubernetes terminology more precisely) depending on the use case. To expose a process running inside a container basically means to expose the port of that container.

### Port forwarding

Let's start with the easiest one (but also, probably a completely useless one, at least for production purposes) &ndash; **port forwarding**. Port forwarding capability is provided by the `kubectl` tool via `kubectl port-forward` command. It is very similar to the `docker -p $HOST_PORT:$TARGET_PORT` command, which you should already be familiar with. The `port-forward` command creates a **proxy** to the container port and binds it to a localhost port.

This is sometimes useful for debugging purposes but you certainly don't want to use this in a production setting and you won't come across this very often (except for tutorials). Also, pods are mortal and if the pod that you're connected to dies, so does your connection.

### Services

Kubernetes `Service` is an abstract way to expose an application running on a set of Pods as a network service. In Kubernetes, a `Service` is an abstraction which defines a logical set of `Pods` and a policy by which to access them. For example, consider a stateless image-processing backend which is running with 3 replicas. Those replicas are fungible &mdash; frontends do not care which backend they use. While the actual Pods that compose the backend set may change, the frontend clients should not need to be aware of that, nor should they need to keep track of the set of backends themselves. The Service abstraction enables this decoupling. <sup>[<a href="#ref:k8s.io/concepts/services">1</a>]</sup>.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: example-service
spec:
  selector:
    app: example-app
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 9376
  type: ClusterIP
```

The `spec` of a `Service` is quite straightforward. There is a `selector` to map the service to an app, `ports` to choose ports that should be mapped and a `type`.


**Publishing Services (ServiceTypes)**

For some parts of your application (for example, frontends) you may want to expose a Service onto an external IP address, that's outside of your cluster.

Kubernetes ServiceTypes allow you to specify what kind of Service you want. The default is `ClusterIP`.

Type values and their behaviors are <sup>[<a href="#ref:k8s.io/concepts/services">1</a>]</sup>:

<div class="subsection">

`ClusterIP` exposes the Service on a cluster-internal IP. Choosing this value makes the Service only reachable from within the cluster. This is the default ServiceType.

<center>
<img src="../../assets/kubernetes-service-clusterip.png" alt="Kubernetes Service: Cluster IP">
<div class="caption">
<a name="fig:kubernetes-service-clusterip" role="figure" aria-label="reference">
<figcaption>Kubernetes Service: ClusterIP</figcaption>
</a>
<span class="has-text-weight-light is-size-7">Courtesy of <a href="https://medium.com/u/2cac56571879?source=post_page-----922f010849e0----------------------">Ahmet Alp Balkan</a></span>
</div>
</center>


`NodePort` exposes the Service on each Node's IP at a static port (the NodePort). A ClusterIP Service, to which the NodePort Service routes, is automatically created. You'll be able to contact the NodePort Service, from outside the cluster, by requesting <NodeIP>:<NodePort>.

<center>
<img src="../../assets/kubernetes-service-nodeport.png" alt="Kubernetes Service: NodePort">
<div class="caption">
<a name="fig:kubernetes-service-nodeport" role="figure" aria-label="reference">
<figcaption>Kubernetes Service: NodePort</figcaption>
</a>
<span class="has-text-weight-light is-size-7">Courtesy of <a href="https://medium.com/u/2cac56571879?source=post_page-----922f010849e0----------------------">Ahmet Alp Balkan</a></span>
</div>
</center>

`LoadBalancer` exposes the Service externally **using a cloud provider's load balancer**<sup class="has-tooltip-multiline" data-tooltip="In a bare-metal deployment, this can be handled for example with metallb"><i class="fas fa-sm fa-question"></i></sup>. NodePort and ClusterIP Services, to which the external load balancer routes, are automatically created.

<center>
<img src="../../assets/kubernetes-service-loadbalancer.png" alt="Kubernetes Service: LoadBalancer">
<div class="caption">
<a name="fig:kubernetes-service-loadbalancer" role="figure" aria-label="reference">
<figcaption>Kubernetes Service: LoadBalancer</figcaption>
</a>
<span class="has-text-weight-light is-size-7">Courtesy of <a href="https://medium.com/u/2cac56571879?source=post_page-----922f010849e0----------------------">Ahmet Alp Balkan</a></span>
</div>
</center>

`ExternalName` maps the Service to the contents of the externalName field (e.g. foo.bar.example.com), by returning a CNAME record with its value. No proxying of any kind is set up. This is quite uncommon and in most cases, you won't probably meet with this type of service.

</div>

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
Note: You need either kube-dns version 1.7 or CoreDNS version 0.0.8 or higher to use the ExternalName type.
</div>
</article>

You can also use `Ingress` to expose your `Service`. `Ingress` is not a `Service` type, but it acts as the entry point for your cluster. It lets you consolidate your routing rules into a single resource as it can expose multiple services under the same IP address. We'll discuss `Ingress` in the next section.

Let's create a service for our example application. As we haven't got a load balancer service deployed in the cluster, we're basically limited to `NodePort` and `ClusterIP`. Let's first try out the default `ClusterIP` type to see how it behaves.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: echoserver  # it is good practice to use the same name for both service and application
spec:
  selector:
    app: echoserver
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 5000  # make sure to use the correct port, remember, we've configured it via `ConfigMap`
  type: ClusterIP
```

Now, if we curl the Cluster IP address from our localhost, it won't work. We need to ssh to a node in the cluster and send the request from there.

```sh
curl -v <ClusterIP>:8080
```

That might not be the desired solution. What we can do then is to change the `type` to `NodePort`. Try it out for yourself and try to access the service using the host IP (as assigned by `multipass`).

## <a name="sec:kubernetes-primitives/ingress"></a> Kubernetes Primitives: Ingress

`Ingress` exposes HTTP and HTTPS routes from outside the cluster to services within the cluster. Traffic routing is controlled by rules defined on the Ingress resource <sup>[<a href="#ref:https://kubernetes.io/docs/concepts/services-networking/ingress/">2</a>]</sup>. Ingresses are a bit more complicated, they can also provide load balancing and SSL / TLS and they require an **ingress controller**.

<center>
<img src="../../assets/kubernetes-ingress.png" alt="Kubernetes Ingress">
<div class="caption">
<a name="fig:kubernetes-ingress" role="figure" aria-label="reference">
<figcaption>Kubernetes Ingress</figcaption>
</a>
<span class="has-text-weight-light is-size-7">Courtesy of <a href="https://medium.com/u/2cac56571879?source=post_page-----922f010849e0----------------------">Ahmet Alp Balkan</a></span>
</div>
</center>

A minimal Ingress resource example:

```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: test-ingress
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
spec:
  rules:
  - http:
      paths:
      - path: /testpath
        pathType: Prefix
        backend:
          serviceName: test
          servicePort: 80
```

Ingress can also specify a `host` to provide named based virtual hosting, i.e.:

```
foo.bar.com --|                 |-> foo.bar.com service1:80
              | 178.91.123.132  |
bar.foo.com --|                 |-> bar.foo.com service2:80
```

which matches the following manifest:

```yaml
apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: name-virtual-host-ingress
spec:
  rules:
  - host: foo.bar.com
    http:
      paths:
      - backend:
          serviceName: service1
          servicePort: 80
  - host: bar.foo.com
    http:
      paths:
      - backend:
          serviceName: service2
          servicePort: 80
```

Ingresses are quite complicated and we certainly did not cover all of it in this section. I recommend checking out the [official documentation](https://kubernetes.io/docs/concepts/services-networking/ingress/) to learn more about them.

## <a name="sec:kubernetes-primitives/volumes"></a> Kubernetes Primitives: Volumes and PVCs

Now, let's head back to the scenario with the `echoserver` application. Say we want to store some kind of data to a **persistent volume**, for example to log each request that comes to the server or on the contrary, load an existing data to a running application. We want to create a `PersistentVolume` for that.

```yaml
# example: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pv0003
spec:
  capacity:
    storage: 5Gi
  volumeMode: Filesystem
  accessModes:
    - ReadWriteOnce
  persistentVolumeReclaimPolicy: Recycle
  storageClassName: slow
  mountOptions:
    - hard
    - nfsvers=4.1
  nfs:
    path: /tmp
    server: 172.17.0.2
```

I encourage you to check the [documentation]https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistent-volumes) to learn more about different types of persistent volumes, like for example `nfs`, `hostPath`, `CSI` or `cephFS`. Explanation how persistent volume provisioning works and what `StorageClass`es are is beyond the scope of this lecture.

Alright, that looks quite complicated. That is why, in most cases, we don't use `PersistentVolume`s directly, however, instead we make use of **dynamic provisioning**. How that works, is that we only specify a `PersistentVolumeClaim`, which is a **request** for a storage of certain qualities. The persistent volume is then provisioned dynamically according to the requirements.

<article class="message is-primary">
<div class="message-header">
<i class="fas fa-2x fa-book-open"></i>
<span>Definition</span>
</div>
<div class="message-body" role="complementary" aria-label="definition">
<a name="def:persistent-volume-claim" aria-label="reference">
A PersistentVolumeClaim (PVC) is a request for storage by a user. It is similar to a Pod. Pods consume node resources and PVCs consume PV resources. Pods can request specific levels of resources (CPU and Memory). Claims can request specific size and access modes (e.g., they can be mounted ReadWriteOnce, ReadOnlyMany or ReadWriteMany, see AccessModes).
</a><sup>[<a href="#ref:k8s.io/concepts/persistent-volumes">3</a>]</sup>
</div>
</article>

There are certain conditions that have to be met in order for dynamic provisioning to work, hence that would typically not be available in a bare-metal Kubernetes deployment. Most cloud providiers, however, already have these setup for you to consume. We would use create a `PersistentVolumeClaim` as such:

```
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: echoserver-storage
spec:
  accessModes:
    - ReadWriteOnce
  storageClassName: fast
  resources:
    requests:
      storage: 2Gi
```

The `StorageClass` is provided by the cloud provider and is directly linked to the volume provisioner (quite an advanced topic). If `storageClassName` is left out, the default `StorageClass` is used.

The `PVC` is then consumed in the following way:

```yaml
# file: deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:  # optional, but it's a good practice
    app: echoserver
  name: engeto-echoserver
  namespace: engeto  # optional, by default to the current namespace
spec:
  selector:
    matchLabels:
      app: echoserver
  template:
    metadata:
      labels:
        app: echoserver
    spec:
      containers:
      - name: main
        image: cermakm/echoserver:latest
        volumeMounts:
        - mountPath: "/var/echoserver/storage"
          name: storage
    volumes:
    - name: storage
      persistentVolumeClaim:
        claimName: echoserver-storage

```


## <a name="sec:kubernetes-primitives/others"></a> Kubernetes Primitives: Other kinds of deployments

There are other ways to deploy an application other than using the `Deployment`. The two other commonly used approaches include `StatefulSet` and `DaemonSet`. We're not going to cover these here in greater detail, but you should at least know when to use them.

`DaemonSet` is quite specific in that it is used to deploy an application to **every node** in the cluster (if not explicitly prohibited). `StatefulSet` is very similar to `Deployment`, but it is used for **stateful** applications. The primary difference is in the way it handles **volumes** when using multiple replicas. With `Deployment`, the volume is **shared by all replicas**, whereas each pod in a `StatefulSet` gets its own volume. Also, `StatefulSet` stores its state on a persistent storage along with additional data, like the order in which it starts the pods etc. It's far more complicated than a <q>simple</q> `Deployment`.

<br>

There is certainly much more to cover here, but that overview should give you a high level intuition about Kubernetes primitives and how they are used together. You should now be able to deploy your own application.

<br>

## References

<!-- -->

<a name="ref:k8s.io/contepts/secrets"></a>Kubernetes (n. d.). <i>Secrets</i>. Retrieved July 13, 2020, from <a href="https://kubernetes.io/docs/concepts/configuration/secret/">https://kubernetes.io/docs/concepts/configuration/secret/</a>

<a name="ref:k8s.io/contepts/services"></a>Kubernetes (n. d.). <i>Service</i>. Retrieved July 13, 2020, from <a href="https://kubernetes.io/docs/concepts/services-networking/service/">https://kubernetes.io/docs/concepts/services-networking/service/</a>

<a name="ref:k8s.io/contepts/ingress"></a>Kubernetes (n. d.). <i>Ingress</i>. Retrieved July 14, 2020, from <a href="https://kubernetes.io/docs/concepts/services-networking/ingress/">https://kubernetes.io/docs/concepts/services-networking/ingress/</a>

<a name="ref:k8s.io/contepts/persistent-volumes"></a>Kubernetes (n. d.). <i>Persistent Volumes</i>. Retrieved July 24, 2020, from <a href="https://kubernetes.io/docs/concepts/storage/persistent-volumes/">https://kubernetes.io/docs/concepts/storage/persistent-volumes/</a>
