<style>{less: ../../../../main.less}</style>

<section class="navigation">
<a href="../../README.md" class="button is-small has-tooltip-bottom is-pulled-right" data-tooltip="Return to the Course Overview"><i class="fas fa-home"></i></a>
<div class="columns is-centered">
    <a href="./02-k8s-cluster.md" class="button is-small has-tooltip-bottom" data-tooltip="Next"><i class="fas fa-arrow-right"></i></a>
    <!-- <a href=""><i class="fas fa-arrow-right"></i></a> -->
</div>
</section>

# <a name="ch:introduction"></a> Kubernetes: Introduction

<section class="intro-section">

<!-- Intro -->

Kubernetes is a container-orchestration platform for automating deployment, scaling and operations of applications <sup>[<a href="#ref:k8s.io">0</a>]</sup>. Open-sourced by Google in 2014, Kubernetes was built based on the search giant's own experience with running containers in production. It's now under the aegis of the Cloud Native Computing Foundation (CNCF), which reports that Kubernetes is the most popular container management tool among large enterprises, used by 83% of respondents in a recent CNCF survey. And in case you're wondering, the name "Kubernetes" originates from Greek and means helmsman or pilot <sup>[<a href="#ref:blog/how-to-monitor-kubernetes-best-practices">1</a>]</sup>.

<table border="0" class="questions-table">
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is Kubernetes?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What's the purpose of Kubernetes?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What are the main concepts of Kubernetes?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What's the difference between Kubernetes and Docker? Is it even comparible? How about Docker swarm?
    </span>
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What are the benefits of Kubernetes?
    </span>
</td></tr>
</table>

<table class="shortcuts-table">
<tr>
    <td>K8s</td>
    <td>Kubernetes</td>
</tr>
</table>

</section>

<br>

## Concepts

To understand the value of Kubernetes, we must first take a look back at how enterprises have deployed applications over the years. In traditional deployments, applications ran on physical servers, an approach that led to resource allocation issues. If multiple applications were running on a single server, for instance, one application could consume most of the resources, hampering performance of other applications. One solution was to run each app on a separate physical server, but this approach was too expensive and led to underutilization of resources <sup>[<a href="#ref:blog/appdynamics/kubernetes">1</a>]</sup>.

Virtualization was the next step, which addressed the limitations of physical servers by running multiple virtual machines (VMs)—each running its own components, including the operating system (OS) and applications—on top of a physical server's CPU. VMs offered numerous benefits, including improved utilization of server resources, lower hardware costs, easier application upgrades, and other scalability enhancements <sup>[<a href="#ref:blog/appdynamics/kubernetes">1</a>]</sup>.
However, VMs are also very costly in terms of physical resources and (in comparison to containers) creating a VM is relatively time-consuming. There is also operational complexity when it comes to provisioning the machines and maintaining them. This adds complexity to the software development lifecycle and limits the portability of apps between public and private clouds and data centers.

This is where Kubernetes and containers come into play. As we've discussed in the Part I of this course, containers are <q>lightweight single-purpose VMs</q> providing runtime for an application, so to speak (I recommend not using this expression too often as it is far from reality and there are significant differencies, as you already know, but it at least gives an idea about their behaviour). As such, containers can be spawned, deleted and re-created very quickly. Which is exactly what Kubernetes does. It manages and orchestrates containers and everything around them, including internal and external networking, storage, monitoring, etc...

Kubernetes provides a framework to run distributed systems with resilience. Upon deployment, you get a Kubernetes cluster &ndash; a set of machines, or nodes, that run containerized applications that Kubernees manages. Cluster and its architectural components will be described in greater detail the following chapters.

#### Kubernetes as container orchestration engine

Since containers are far more efficient, fast, and lightweight than traditional virtualization, it's unsurprising that enterprises with large application deployments may deploy multiple containers as one or more container clusters. But this environment comes with its own set of challenges, as large, distributed containerized applications often become difficult to coordinate.

Kubernetes is used most often with Docker, the leading containerization platform. But Kubernetes also supports other container systems that meet container image format and runtime standards set by the Open Container Initiative (OCI), an open source technical community supervised by The Linux Foundation. Alternatives to Kubernetes include Docker Swarm, a container orchestrator bundled with Docker, and Apache Mesos.

#### How Kubernetes works

Kubernetes groups containers that make up an application into logical units for easy management and discovery. Kubernetes builds upon 15 years of experience of running production workloads at Google, combined with best-of-breed ideas and practices from the community <sup>[<a href="#ref:k8s.io">0</a>]</sup>. It runs everywehere, giving you the freedom to take advantage of on-premises, hybrid, or public cloud infrastructure, letting you effortlessly move workloads to where it matters to you <sup>[<a href="#ref:blog/appdynamics/kubernetes">1</a>]</sup>.

Kubernetes provides a framework to run distributed systems with resilience. Upon deployment, you get a Kubernetes cluster—a set of machines, or nodes, that run containerized applications that Kubernees manages. A cluster has at least one <sup>[<a href="#ref:blog/appdynamics/kubernetes">1</a>]</sup>:

- Worker node to host the pods. (Each pod is a group of one or more containers.)
- Master node to manage the worker nodes and pods in the cluster. 

#### Features & Advanteges

The following set of features describes Kubernetes and its advanteges for production workloads <sup>[<a href="#ref:k8s.io">0</a>]</sup>:

<div class="tile is-vertical is-centered is-ancestor">

  <div class="tile is-centered is-parent">
    <div class="tile is-three-fifths is-child box">
      <p class="title has-text-left is-size-6">Service discovery and load balancing</p>
      <p class="content">
        No need to modify your application to use an unfamiliar service discovery mechanism. Kubernetes gives Pods their own IP addresses and a single DNS name for a set of Pods, and can load-balance across them.
      </p>
    </div>
    <div class="tile is-three-fifths is-right is-child box">
      <p class="title has-text-left is-size-6">Service Topology</p>
      <p class="content">
        Routing of service traffic based upon cluster topology.
      </p>
    </div>
  </div>

  <div class="tile is-parent">
    <div class="tile is-child box">
      <p class="title has-text-left is-size-6">Storage orchestration</p>
      <p class="content">
        Automatically mount the storage system of your choice, whether from local storage, a public cloud provider such as GCP or AWS, or a network storage system such as NFS, iSCSI, Gluster, Ceph, Cinder, or Flocker.
      </p>
    </div>
    <div class="tile is-child box">
      <p class="title has-text-left is-size-6">Self-healing</p>
      <p class="content">
      Restarts containers that fail, replaces and reschedules containers when nodes die, kills containers that don't respond to your user-defined health check, and doesn't advertise them to clients until they are ready to serve.
      </p>
    </div>
  </div>

  <div class="tile is-parent">
    <div class="tile is-child box">
      <p class="title has-text-left is-size-6">Automated rollouts and rollbacks</p>
      <p class="content">
        Kubernetes progressively rolls out changes to your application or its configuration, while monitoring application health to ensure it doesn't kill all your instances at the same time. If something goes wrong, Kubernetes will rollback the change for you. Take advantage of a growing ecosystem of deployment solutions.
      </p>
    </div>
    <div class="tile is-child box">
      <p class="title has-text-left is-size-6">Secret and configuration management</p>
      <p class="content">
        Deploy and update secrets and application configuration without rebuilding your image and without exposing secrets in your stack configuration.
      </p>
    </div>
  </div>

  <div class="tile is-parent">
    <div class="tile is-child box">
      <p class="title has-text-left is-size-6">Automatic bin packing</p>
      <p class="content">
        Automatically places containers based on their resource requirements and other constraints, while not sacrificing availability. Mix critical and best-effort workloads in order to drive up utilization and save even more resources.
      </p>
    </div>
    <div class="tile is-child box">
      <p class="title has-text-left is-size-6">Batch execution</p>
      <p class="content">
        In addition to services, Kubernetes can manage your batch and CI workloads, replacing containers that fail, if desired.
      </p>
    </div>
  </div>

  <div class="tile is-parent">
    <div class="tile is-child box">
      <p class="title has-text-left is-size-6">IPv4/IPv6 dual-stack</p>
      <p class="content">
        Allocation of IPv4 and IPv6 addresses to Pods and Services
      </p>
    </div>
    <div class="tile is-child box">
      <p class="title has-text-left is-size-6"> Horizontal scaling</p>
        Scale your application up and down with a simple command, with a UI, or automatically based on CPU usage.
      </p>
    </div>
  </div>

</div>

<br>

## References

<a name="ref:k8s.io"></a>Kubernetes.io (n.d.). Retrieved June 21, 2020, from <a href="https://kubernetes.io/">https://kubernetes.io/</a>

<a name="ref:blog/appdynamics/kubernetes"></a>appdynamics (n.d.). <i>What is Kubernetes?</i>. Retrieved June 16, 2020, from <a href="https://www.appdynamics.com/solutions/cloud/cloud-monitoring/kubernetes-monitoring/how-to-monitor-kubernetes-best-practices">https://www.appdynamics.com/solutions/cloud/cloud-monitoring/kubernetes-monitoring/how-to-monitor-kubernetes-best-practices</a>

<!-- -->

[kubernetes]: https://kubernetes.io