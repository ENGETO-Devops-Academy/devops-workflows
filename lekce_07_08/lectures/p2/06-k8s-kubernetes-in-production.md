
<style>{less: ../../../../main.less}</style>

<section class="navigation">
<a href="../../README.md" class="button is-small has-tooltip-bottom is-pulled-right" data-tooltip="Return to the Course Overview"><i class="fas fa-home"></i></a>
<div class="columns is-centered">
    <a href="./05-k8s-multi-tenancy.md" class="button is-small has-tooltip-bottom" data-tooltip="Previous"><i class="fas fa-arrow-left"></i></a>
    <a href="../final-project.md" class="button is-small has-tooltip-bottom" data-tooltip="Next"><i class="fas fa-arrow-right"></i></a>
</div>
</section>

# <a name="ch:kubernetes-in-production"></a> Kubernetes in production

<section class="intro-section">

Kubernetes has been increasingly adopted by companies all over the world. It has come to the point where it's already being used in production environments, and that's a big thing! However, as the Kubernetes ecosystem has taken off, the fast pace comes with a price &mdash; and not a small one &mdash; and that is **maintainance**. In this final lesson, we will review a couple things that are vital in order to maintain a productioin quality Kubernetes cluster and stain sane while doing so.

<table border="0" class="questions-table">
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What should you be wary of the most if you run a Kubernetes in production?
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    How many nodes should there be at least in a production cluster? How many is the recommended minimum? Why?
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is the tool of choice to provision a cluster?
</td></tr>
</table>

</section>

## <a name="sec:kubernetes-in-production/maintenance"></a> Maintaining a healthy cluster

Maintaining a cluster is surprisingly enough a pretty hard job. Let me give you a few tips to stay sane while doing so.

What strike me as especially painful was the quick releases. To the date of this writing, there is a 14 day release cycle! Now imagine that one week you bump up a version of the cluster and the next one that version is depricated. True story! My recommendation is to stay in touch with the Kubernetes community, read the mailing list or join the slack channel to be prepared for the changes that are to come (like deprecated APIs, for example).

Besides cluster updates, you should always strive for the perfection when it comes to automation. Have your provisioning scripts, ansible playbooks, Jenkins jobs, Terraform jobs or whatever you deem necessary at your fingertips so that you could rotate a cluster in a blink of an eye if necessary. The cluster *will* go down eventually, be prepared for that.

And the third advice, have your monitoring and logging apart of the cluster ((at least the production one). It will save you a lot of trouble while debugging why the cluster is down. If you have your monitoring and logging deployed in the cluster (I know, it is tempting given the ease of deployment with a Helm chart), you'll eventually hit a chicken-egg problem &mdash; *How do I debug the cluster, when the cluster is down and I can't access my monitoring dashboard deployed in that cluster?!*.

## <a name="sec:kubernetes-in-production/high-availability"></a> High availability

When it comes to production environments, reliability and high availability is the key to fulfiling SLAs. In order to achieve high availability in Kubernetes, it's seemingly simple &ndash; just add more nodes. Why *seemingly*? Well, let's discuss what happens when there are multiple nodes in a cluster.
First, the nodes can either be **master** or **worker** nodes. Multi-master architecture is a little bit more complicated, but it's worth it when it comes to HA. How ever, we have to concider the **number of masters**. That is because a phenomenon called <q>quorum</q>. In HA scenario, Kubernetes will maintain a copy of the etcd database on each master, but elects a leader of the control plane. The election will happen in a distributed consensus algorithm called **raft** and a **quorum has to be maintained** which is decided by (n/2)+1 voters. For that reason, it is often a good case to have **odd number of nodes** and NOT less than 3 (ideally 5) out of which at least are master nodes.

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
If you happen to have only two nodes, don't ever make both of them masters. That would make things even worse as they would be in a serious fight in order to reach a quorum and the cluster might go stale. In this case, <b>one master is better than two masters</b>!
</div>
</article>

## <a name="sec:kubernetes-in-production/multi-cluster"></a> Multiple clusters

Multiple clusters are a great way to isolate production ready services from the environment used by developers. In practice, there are often at least two clusters &ndash; **staging (or test) cluster** and the **production** cluster. Sometimes there can even be a smaller **dev** cluster dedicated for experimental workloads. It is not uncommon for the operations to create a separate cluster for logging / monitoring and other operations, something like a <q>supervisor cluster</q>.

In some rare cases, it can also be desirable for the clusters to create **a mesh** (or to **federate clusters**). A cluster mesh is a rare occurance and is quite demanding to setup and maintain, especially when it comes to networking and storage. This capability is not native to Kubernetes, but there exist third-party open source solutions created by the community, I would recommend you to take a look at [Cilium](https://cilium.io).

## <a name="sec:kubernetes-in-production/cluster-provisioning"></a> Cluster Provisioning

Provisioning a lightweight cluster is quite easy with the tools like [minikube](https://minikube.sigs.k8s.io/docs/) or [microk8s](https://microk8s.io), as we've already discuss. These tools, however, are not really production-ready and they also don't provide many configuration options. When it comes to production-ready deployment, there is swiss knife called [Kubespray](https://github.com/kubernetes-sigs/kubespray), which is being actively developed and maintained and also provides a certain degree of quality backed by integration tests. If Kubespray feels too heavy for you, I would then choose microk8s as my second choice. It has a great support and the community pushes it forward.

## <a name="sec:kubernetes-in-production/enterprise-kubernetes"></a> Enterprise Kubernetes

If you don't really want to provision your own cluster on a bare-metal machine, there is of course a way to opt for a cloud provider. The most commonly used are [Google Cloud (GKE)](https://cloud.google.com), [Amazon Web Services (AWS)](https://aws.amazon.com) and [Microsoft Azure](https://azure.microsoft.com/cs-cz/).
