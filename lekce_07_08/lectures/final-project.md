
<style>{less: ../../../main.less}</style>

<section class="navigation">
<a href="../README.md" class="button is-small has-tooltip-bottom is-pulled-right" data-tooltip="Return to the Course Overview"><i class="fas fa-home"></i></a>
</section>

# <a name="ch:final-project"></a> Final project

<section class="intro-section">

Kubernetes is quite a beast. I know that and you should know that too. In this course we've just briefly scratched the surface of what's possible. However, we did go through the vital parts that should give you a pretty good idea how to deal with a Kubernetes cluster when you're sit in front of one. In your final assignment, practice these vital parts... and a little bit more.

Your assignment will consist of multiple steps and it will test your knowledge in other areas as well.

</section>

## Assignment

1) Use Kubespray do provision a Kubernetes cluster called `engeto.cluster.local` on **3 nodes** with **1 master**.
    - Let Kubespray deploy `cert-manager`
    - Let Kubespray deploy `ingress` controller
    - Let Kubespray deploy `metallb`
1) Create a docker builder container which will be used to build golang images, that container is used as a toolbox, you can exec into that container, fetch your git repo and build it. The container **uses its hosts docker daemon**, it doesn't create its own socket (do NOT use DIND<sup class="has-tooltip-multiline" data-tooltip="Short for Docker-in-docker."><i class="fas fa-sm fa-question"></i></sup>).
2) Build  the example application in that container.
3) Deploy the example application such that it is deployed to *each node in the cluster*. That deployment to nodes should be dynamic, so that each time a new node is added to the cluster, the application gets deployed to that node automatically.
4) Make the application accessible from outside of the cluster.
5) [optional] Give it a nice domain name.
