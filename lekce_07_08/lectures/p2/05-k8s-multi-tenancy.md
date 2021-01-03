
<style>{less: ../../../../main.less}</style>

<section class="navigation">
<a href="../../README.md" class="button is-small has-tooltip-bottom is-pulled-right" data-tooltip="Return to the Course Overview"><i class="fas fa-home"></i></a>
<div class="columns is-centered">
    <a href="./04b-k8s-primitives.md" class="button is-small has-tooltip-bottom" data-tooltip="Previous"><i class="fas fa-arrow-left"></i></a>
    <a href="./06-k8s-kubernetes-in-production.md" class="button is-small has-tooltip-bottom" data-tooltip="Next"><i class="fas fa-arrow-right"></i></a>
</div>
</section>

# <a name="ch:kubernetes-multi-tenancy"></a> Multi-Tenancy and cluster policies

<section class="intro-section">

It is often the case that a cluster is shared by multiple users. These are refered to as <q>tenants</q>. In this chapter, we'll take a look at how we can isolate the tenants from each other to minimize the damange that a compromised tenant can do to the cluster or other tenants and how to restrict or grant the rights for particular users.

<table border="0" class="questions-table">
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is multi-tenancy?
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What is RBAC?
</td></tr>
<tr><td class="question">
    <i class="fas fa-question-circle"></i>
    <span role="complementary" aria-label="question">
    What are the best practices when it comes to cluster security and namespace isolation?
</td></tr>
</table>

<table class="shortcuts-table">
<tr>
    <td>RBAC</td>
    <td>Role base access control</td>
</tr>
</table>


</section>

## <a name="sec:kubernetes-multi-tenancy/multi-tenancy"></a> Multi-tenancy

The Kubernetes multi-tenancy SIG defines a tenant as representing a group of Kubernetes users that has access to a subset of cluster resources (compute, storage, networking, control plane and API resources) as well as resource limits and quotas for the use of those resources. Resource limits and quotas lay out tenant boundaries. These boundaries extend to the control plane allowing for grouping of the resources owned by the tenant, limited access or visibility to resources outside of the control plane domain and tenant authentication <sup>[<a href="#ref:replex/multi-tenancy">3</a>]</sup>.

<div class="subsection">

**Soft multi-tenancy**

Soft multi-tenancy trusts tenants to be good actors and assumes them to be non-malicious. Soft multi-tenancy is focused on minimising accidents and managing the fallout if they do.

**Hard multi-tenancy**

Hard multi-tenancy assumes tenants to be malicious and therefore advocates zero trust between them. Tenant resources are isolated and access to other tenant’s resources is not allowed. Clusters are configured in a way that isolate tenant resources and prevent access to other tenant’s resources.

</div>

## <a name="sec:kubernetes-multi-tenancy/isolating-tenants"></a> Isolating tenants in namespaces

**Tenants**

Multi-tenancy is an alternative to managing many single-tenant clusters, which would have much greater operational overhead and would also be significantly more resource-consuming.

Kubernetes uses `Namespace`s to isolate tenants. These tenants are often distinct teams within the organization. A `Namespace` resource has been presented in the previous chapter as a way to organize applications in a logical way. But that is by no means the only reason the namespaces exist. Within a namespace, **policies** can be configured with respect to **`Pod`s**, **resource quotas** and **API access**.

## <a name="sec:kubernetes-multi-tenancy/enterprise-multi-tenancy"></a> Enterprise multi-tenancy

In an enterprise environment, the tenants of a cluster are distinct teams within the organization. Typically, each tenant has a corresponding namespace. Alternative models of multi-tenancy with a tenant per cluster, or a tenant per Google Cloud project, are harder to manage. Network traffic within a namespace is unrestricted. Network traffic between namespaces must be explicitly whitelisted. These policies can be enforced using Kubernetes network policy <sup>[<a href="#ref:cloud.google/concepts/multitenancy">0</a>]</sup>.

The users of the cluster are divided into three different roles, depending on their privilege <sup>[<a href="#ref:cloud.google/concepts/multitenancy">0</a>]</sup>:

<div class="subsection">

**Cluster administrator**

This role is for administrators of the entire cluster, who manage all tenants. Cluster administrators can create, read, update, and delete any policy object. They can create namespaces and assign them to namespace administrators.

**Namespace administrator**

This role is for administrators of specific, single tenants. A namespace administrator can manage the users in their namespace.

**Developer**

Members of this role can create, read, update, and delete namespaced non-policy objects like Pods, Jobs, and Ingresses. Developers only have these privileges in the namespaces they have access to.

</div>

## <a name="sec:kubernetes-multi-tenancy/policy-enforcement"></a> Policy enforcement

There are multiple ways to enforce policies in a Kubernetes cluster. We'll focus only on those which are Kubernetes native &mdash; there are other policy enforcement tools specific to the various cluster providers. These Kubernetes native policy enforcement tools are primarly **Role base access control (RBAC)**, **network policies**, **resource quotas** and **pod security policies**.

### RBAC

RBAC is built into Kubernetes and grants granular permissions for specific resources and operations within your clusters <sup>[<a href="#ref:cloud.google/concepts/mutitenancy">0</a>]</sup>. RBAC permissiones are linked to **users**, **groups** and **service accounts**.

You define your RBAC permissions by creating the following kinds of Kubernetes objects <sup>[<a href="#ref:cloud.google/concepts/rbac">1</a>]</sup>:

<div class="subsection">

`ClusterRole` or `Role`: defines a set of resource types and operations that can be assigned to a user or group of users in a cluster (`ClusterRole`), or a `Namespace` (`Role`), but does not specify the user or group of users.

`ClusterRoleBinding` or `RoleBinding`: assigns a `ClusterRole` or `Role` to a user or group of users. A `ClusterRoleBinding` works with a `ClusterRole`, and a `RoleBinding` works with either a `ClusterRole` or a `Role`.

</div>

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
RBAC roles are purely additive &ndash; there are no "deny" rules. When structuring your RBAC roles, you should think in terms of "granting" users access to cluster resources.
</div>
</article>

<!-- TODO: Example -->

### Network Policies

Cluster network policies give you control over the communication between your cluster's Pods. Policies specify which namespaces, labels, and IP address ranges a Pod can communicate with <sup>[<a href="#ref:cloud.google/concepts/multitenancy">0</a>]</sup>.

<!-- TODO: Example -->

### Resource Quotas

Resource quotas are a tool for administrators to address the problem of one team using more than its fair share of resources.

Resource quotas work like this <sup>[<a href="ref:cloud.google/concepts/resource-quotas">2</a>]</sup>:

- Different teams work in different namespaces. Currently this is voluntary, but support for making this mandatory via ACLs is planned.
- The administrator creates one ResourceQuota for each namespace.
- Users create resources (pods, services, etc.) in the namespace, and the quota system tracks usage to ensure it does not exceed hard resource limits defined in a ResourceQuota.
- If creating or updating a resource violates a quota constraint, the request will fail with HTTP status code 403 FORBIDDEN with a message explaining the constraint that would have been violated.
- If quota is enabled in a namespace for compute resources like cpu and memory, users must specify requests or limits for those values; otherwise, the quota system may reject pod creation. Hint: Use the LimitRanger admission controller to force defaults for pods that make no compute resource requirements. See the walkthrough for an example of how to avoid this problem.

### Pod security policies

This is a topic that is out of the scope of this lesson and you will typically not stumble upon this. Pod security policies require an admission controller, which, again, is a bit of an advanced topic. Feel free to explore the official [documentation](https://kubernetes.io/docs/concepts/policy/pod-security-policy/) yourself to learn more.

## <a name="sec:kubernetes-multi-tenancy/best-practices"></a> Multi-Tenancy & Security Best practices

There're a couple of points to follow when talking about Kubernetes best practices as far as multi-tenancy and security policies are concerned.

### Categorize namespaces

1) System namespaces

There is at least one namespace which should be used by NONE but the cluster admins &ndash; **kube-system**. In this namespace resides the very core functionality of the cluster, like controllers, DNS, API, proxy, etc... NO users (and probably not even un-experienced cluster admins) should be allowed to temper with these and NO applications should be deployed to this namespace unless you *really* have a reason to do so and you *really* know what you're doing.<br>
The other namespace which is usually pretty untouched is the **default** namespace. There's really no reason to put anything into it. It should serve as an entrypoint to the new tenants.

2) Service & Operation namespaces

You will most likely deploy plenty of applications into the cluster. Some of them, however, will be only for the operational purposes, like logging, monitoring, administration and other services more or less vital for the cluster users &mdash; but not directly operated by them and often even not in direct interaction with tenants. These will typically have dedicated namespaces, often logically grouping multiple services that share similar purpose, like `monitoring` namespace or `operators`.

Examples of namespaces for dedicated components are `cert-manager`, `db`, `gitlab`, `nginx-ingress`, `rook`, `vault`, etc...

3) Tenant namespaces

Basically, the rest of the namespaces will belong to the tenants. It can be the case that a tenant has multiple namespaces, be prepared for such requirements, but it should not be the standard scenario.

### Create cluster personas

Another best practice is to create a hierarchy of cluster personas scoped to varying levels of permissions and the operations they can perform. The Kubernetes multi-tenancy SIG outlines four such personas <sup>[<a href="#ref:replex/multi-tenancy">3</a>]</sup>:

<div class="subsection">

**Cluster admin**

Full read/write privileges for all resources in the cluster including those owned by tenants

**Cluster view**

Read privileges for all resources in the cluster including those owned by tenants

**Namespace admin**

Owner of a namespace with full privileges scoped to that particular namespace.

**Namespace user**

A user of a namespace. This tenant usually has the lowest privileges of all which only allow him do deploy an application to a namespace and to view/list namespaced resources.

</div>

### Isolate tenant namespaces

The whole purpose of multi-tenancy is to isolate tenant namespaces. That means that the namespaces should not be able to communicate with each other or even be visible to each other. A tenant from one namespace should not be able to access resources in another namespace, schedule a workload to that namespace or even list namespaces (that is a cluster-level operation).

### Limit access to shared resources

There are resources which are shared among all tenants in the cluster. These are called **cluster-scoped** or just **cluster-level** resources. If not handled carefully, a tenant of a single namespace can modify that resource which would affect the whole cluster. That is certainly NOT acceptable and you should prevent such situations. An example of such resources are `ClusterRole`,`ClusterRoleBinding` or `CustomResourceDefinition`. Tenants might therefore have only `get` and `list` access to these resources to make use of them.<br>
The contrary to these resources are **namespaced** resources. These are tied to a particular namespace, like `Role`s, `RoleBinding`s and `ServiceAccount`s and are in the ownership of namespace admins.

<article class="message is-info">
<div class="message-header">
<i class="fas fa-2x fa-info-circle"></i>
<span>Note</span>
</div>
<div class="message-body" role="complementary" aria-label="note">
There are other ways to manage security of a cluster, like admission controllers and webhooks, network policies, etc... These are beyond the scope of this lecture.
</div>
</article>

### Prevent use of HostPath volumes to increase security

We've already talked about pod security policies. A practical use case is to prevent tenants to use a `HostPath` volume, which gives a `Pod` direct access to the underlying host's filesystem.

<br>

## References

<!-- -->

<a name="ref:cloud.google/concepts/multitenancy"></a>Kubernetes (n. d.). <i>Cluster multi-tenancy</i>. Retrieved July 24, 2020, from <a href="https://cloud.google.com/kubernetes-engine/docs/concepts/multitenancy-overview">https://cloud.google.com/kubernetes-engine/docs/concepts/multitenancy-overview</a>

<a name="ref:cloud.google/concepts/rbac"></a>Kubernetes (n. d.). <i>RBAC</i>. Retrieved July 24, 2020, from <a href="https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control">https://cloud.google.com/kubernetes-engine/docs/how-to/role-based-access-control</a>

<a name="ref:cloud.google/concepts/resource-quotasj"></a>Kubernetes (n. d.). <i>Resource Quotas</i>. Retrieved July 24, 2020, from <a href="https://kubernetes.io/docs/concepts/policy/resource-quotas/">https://kubernetes.io/docs/concepts/policy/resource-quotas/</a>

<a name="ref:replex/multi-tenancy"></a>Replex (2020). <i>Kubernetes in Production: Best Practices and Checklist for Multi-Tenancy</i>. Retrieved July 30, 2020, from <a href="https://www.replex.io/blog/kubernetes-in-production-best-practices-and-checklist-for-multi-tenancy">https://www.replex.io/blog/kubernetes-in-production-best-practices-and-checklist-for-multi-tenancy</a>
