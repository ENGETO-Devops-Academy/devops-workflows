
<style>{less: ../../../../main.less}</style>

<section class="navigation">
<a href="../../README.md" class="button is-small has-tooltip-bottom is-pulled-right" data-tooltip="Return to the Course Overview"><i class="fas fa-home"></i></a>
<div class="columns is-centered">
    <a href="./03-k8s-api.md" class="button is-small has-tooltip-bottom" data-tooltip="Previous"><i class="fas fa-arrow-left"></i></a>
    <a href="../p2/04b-k8s-primitives.md" class="button is-small has-tooltip-bottom" data-tooltip="Next"><i class="fas fa-arrow-right"></i></a>
</div>
</section>

# <a name="ch:kubernetes-primitives"></a> Kubernetes Primitives: Introduction

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

## <a name="sec:kubernetes-primitives/kubectl"></a> Setting up `kubectl`

In order for us to communicate with the cluster, we need to setup the `kubectl` tool properly. There is already a thorough [documentation](https://kubernetes.io/docs/tasks/tools/install-kubectl/) about how to install the `kubectl` tool, so I am not going to go through the installation in this text.
However, we will setup `kubectl` toghether, as it might get sort of tricky when running in a local setup with `multipass` and `microk8s` like we do.

So, to begin with, the source of truth for the `kubectl` tool is a **config** file, we call it **kubeconfig** and by default, it resides in `~/.kube/config`.

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
There is also <code>KUBECONFIG</code> environment variable and <code>--kubeconfig</code> parameter to most of the Kubernetes components) which you can use to point to a specific config file.
</div>
</article>

As we've provisioned our cluster with `microk8s` and `multipass`, we need to retrieve the `kubectl` from the VM which is the Kubernetes **master**, that is `k8s-1`. As `microk8s` stores it in a different location, rather than searching for it (I can tell you where it is, but I am not going to spoil the fun in case you wanted to look for it), we can retrieve it by using `microk8s kubectl` command which is a proxy to the `kubectl`.

```
microk8s kubectl config view
```

```yaml
apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: DATA+OMITTED
    server: https://127.0.0.1:16443
  name: microk8s-cluster
contexts:
- context:
    cluster: microk8s-cluster
    user: admin
  name: microk8s
current-context: microk8s
kind: Config
preferences: {}
users:
- name: admin
  user:
    token: "************************************************************"
```

You can retrieve it with the following command:

```
multipass exec k8s-1 -- microk8s kubectl config view --raw > ~/.kube/config.microk8s
```

The `--raw` parameter is important, otherwise the `certificate-authority-data` field would be redacted. Now, the last thing we have to do is to change the `server` url to the host ip. In my case, that would be `https://192.168.64.2:16443`.

Now just change the `KUBECONFIG` environment variable to point to the kubeconfig file. Verify that you've successfully setup `kubectl` by listing the nodes that you have in the cluster. You should see similar output:

```
$ kubectl get nodes
NAME           STATUS   ROLES    AGE   VERSION
192.168.64.3   Ready    <none>   29h   v1.18.5
192.168.64.4   Ready    <none>   29h   v1.18.5
k8s-1          Ready    <none>   29h   v1.18.5
```

That's it! You can now communicate with the Kubernetes API server via `kubectl`.

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">

<b>TIP: </b> Take a look at [oh-my-zsh kubectl plugin](https://github.com/ohmyzsh/ohmyzsh/blob/master/plugins/kubectl/kubectl.plugin.zsh) if you're using `zsh` (if not, I can only recommend that you do!). It makes working with `kubectl` so much more fun. E. g. the previous command to list nodes would be aliased simply to `kgno`.
</div>
</article>

## <a name="sec:kubernetes-primitives/namespace"></a> Kubernetes Primitives: Namespace

First of all, we need a `Namespace` to deploy the application to. A `Namespace` is a resource as any other. There are only a couple of namespaces by default and as a rule of thumb, you shouldn't mess with them very much unless you know what you're doing. The system namespace is called `kube-system` and it usually runs the vital Kubernetes components like **controllers**, **api-server** or **metrics-server**.

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
<b>TIP:</b> Try to list the namespaces that are currently present in the cluster. You can do that by running `kubectl get namespaces` or `kgns` with the zsh plugin.
</div>
</article>

That being said, `KUBECONFIG` stores the current **context**, i. e. the cluster to communicate with and also the namespace. By default you're placed in the `default` namespace.

Let's create a new namespace called `engeto`.

```
kubectl create namespace engeto
```

and change the context

```
kubectl config set-context --current --namespace=engeto
```

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
This command is really heavy to remember. With the zsh plugin, you would just type `kcn engeto`.
</div>
</article>

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">

<b>TIP:</b> I recommend using [k9s](https://k9scli.io) to visually explore your cluster, but you can also setup Kubernetes dashboard (`microk8s` provides an addon for it as well).
</div>
</article>


## References

<!-- -->

<a name="ref:k8s.io/contepts/secrets"></a>Kubernetes (n. d.). <i>Secrets</i>. Retrieved July 13, 2020, from <a href="https://kubernetes.io/docs/concepts/configuration/secret/">https://kubernetes.io/docs/concepts/configuration/secret/</a>

<a name="ref:k8s.io/contepts/services"></a>Kubernetes (n. d.). <i>Service</i>. Retrieved July 13, 2020, from <a href="https://kubernetes.io/docs/concepts/services-networking/service/">https://kubernetes.io/docs/concepts/services-networking/service/</a>

<a name="ref:k8s.io/contepts/ingress"></a>Kubernetes (n. d.). <i>Ingress</i>. Retrieved July 14, 2020, from <a href="https://kubernetes.io/docs/concepts/services-networking/ingress/">https://kubernetes.io/docs/concepts/services-networking/ingress/</a>

<a name="ref:k8s.io/contepts/persistent-volumes"></a>Kubernetes (n. d.). <i>Persistent Volumes</i>. Retrieved July 24, 2020, from <a href="https://kubernetes.io/docs/concepts/storage/persistent-volumes/">https://kubernetes.io/docs/concepts/storage/persistent-volumes/</a>
