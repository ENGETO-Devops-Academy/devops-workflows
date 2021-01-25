<p align="center">
  <img src="/assets/engeto-logo.png" alt="ENGETO Academy">
</p>

# Lesson #9: Monitoring, logging & alerting

## 1. Introduction
Logging, monitoring and alerting are at the core of a healthy infrastructure. As a DevOps engineer, you should always aim for perfection as far as these three pillars are concerned. 

Thorough, persistent and accessible application or system logs can save a great amount of time — and therefore money.

### What is your responsibility?
Devops engineers **ARE NOT** responsible for:
 * how an application logs **info**, **debug** and **error** messages
 * how thorough and verbose application logs are

Devops engineers** **ARE** responsible for:
 * log **persistency** and **accessibility**
 * monitoring systems and dashboards
 * proper alerting via suitable notification channels

## 2. Logging
> A log is a text output of a software which contains a record of the events related to the runtime of that particular software. To log means to direct that output to a particular destination, i.e. standard output (stdout), standard error (stderr) or a file.

At the very least, it is important to have a persistent and accessible storage of these logs and to allow developers to easily direct their logs to that storage. That is the job of a devops engineer.

### System log
> The system log (syslog) contains a record of the operating system (OS) events that indicates how the system processes and drivers were loaded. The syslog shows informational, error and warning events related to the computer OS. By reviewing the data contained in the log, an administrator or user troubleshooting the system can identify the cause of a problem or whether the system processes are loading successfully.

The OS maintains a log of events that helps in monitoring, administering and troubleshooting the system in addition to helping users get information about important processes.

Some events include system **errors**, **warnings**, **startup messages**, **system changes**, **abnormal shutdowns**, etc. This list is applicable to most versions of the three common OSs (Windows, Linux and Mac OS).

The logs handled by syslog are available in `/var/log/` directory on Linux system.

Example of log structure on Linux:
```bash
ansible@host:~$ ll /var/log/
total 17824
drwxrwxr-x  12 root      syslog             4096 Aug 15 00:00 ./
drwxr-xr-x  14 root      root               4096 Jul 12 09:20 ../
-rw-r--r--   1 root      root              24836 Aug 12 13:45 alternatives.log
drwxr-xr-x   2 root      root               4096 Aug 13 06:15 apt/
-rw-r-----   1 syslog    adm              477698 Aug 15 14:25 auth.log
-rw-r--r--   1 root      root             104003 Jul 12 09:18 bootstrap.log
-rw-rw----   1 root      utmp                384 Aug 15 08:45 btmp
-rw-r--r--   1 root      root               8258 Aug 12 13:30 cloud-init-output.log
-rw-r--r--   1 syslog    adm              133354 Aug 12 13:30 cloud-init.log
drwxr-xr-x   2 root      root               4096 Aug 13 12:43 containers/
drwx------   3 root      root               4096 Aug 12 13:46 crio/
drwxr-xr-x   2 root      root               4096 May  6 00:26 dist-upgrade/
-rw-r--r--   1 root      adm              149918 Aug 12 13:30 dmesg
-rw-r--r--   1 root      adm              150458 Aug 12 13:01 dmesg.0
-rw-r--r--   1 root      root             417448 Aug 13 06:15 dpkg.log
-rw-r--r--   1 root      root              32032 Aug 12 13:47 faillog
drwxr-xr-x   3 root      root               4096 Aug 12 12:59 installer/
drwxr-sr-x+  3 root      systemd-journal    4096 Aug 12 13:00 journal/
-rw-r-----   1 syslog    adm              390593 Aug 15 08:55 kern.log
drwxr-xr-x   2 landscape landscape          4096 Aug 12 13:02 landscape/
-rw-rw-r--   1 root      utmp             292292 Aug 15 14:25 lastlog
drwxr-xr-x  11 root      root               4096 Aug 13 12:01 pods/
drwx------   2 root      root               4096 Jul 12 09:17 private/
-rw-r-----   1 syslog    adm             5648042 Aug 15 14:30 syslog
-rw-r-----   1 syslog    adm             9334940 Aug 15 00:00 syslog.1
-rw-r-----   1 syslog    adm             1278250 Aug 14 00:00 syslog.2.gz
-rw-------   1 root      root                  0 Jul 12 09:18 ubuntu-advantage.log
drwxr-x---   2 root      adm                4096 Aug 12 13:01 unattended-upgrades/
-rw-rw-r--   1 root      utmp              15360 Aug 15 14:25 wtmp

```

As you can see, there are logs from various programs and services. Of particular importance are the following ones:
> **auth.log** contains information about system authorization including logins and other authentication mechanisms

> **wtmp** contains login records (who is currently logged into the system)

> **dmesg** contains kernel ring buffer information and low-level runtime info (useful for debugging container failures, like OOM)

> **journals** can be viewed using journalctl command. As you can see, it is a directory containing a bunch of *.journal files:
```shell
ansible@host:~$ tree /var/log/journal/
/var/log/journal/
└── aa393256cd6f4adfb003b5e9f42a0468
    ├── system.journal
    ├── system@f773452f1aa549a4b307cc18ee1057f6-0000000000000001-0005acac1a2f3b60.journal
    ├── system@f773452f1aa549a4b307cc18ee1057f6-0000000000001c0a-0005acacba9c5c7a.journal
    ├── user-1000.journal
    ├── user-1000@861deed8ca144c89b225064ab35a8a73-00000000000009ce-0005acac1f28f756.journal
    └── user-1000@861deed8ca144c89b225064ab35a8a73-0000000000001c5d-0005acacbbeba860.journal
1 directory, 6 files

```

An example of such logs could be:
```shell
ansible@host:~$ journalctl
-- Logs begin at Wed 2020-08-12 13:00:57 CEST, end at Sat 2020-08-15 14:56:48 CEST. --
Aug 12 13:00:57 host kernel: Linux version 5.4.0-42-generic (buildd@lgw01-amd64-038) (gcc version 9.3.0 (Ubuntu 9.3.0-10ubuntu2)) #46-Ubuntu SMP Fri Jul 10 00:24:02 UTC 2020 (Ubuntu 5.>
Aug 12 13:00:57 host kernel: Command line: BOOT_IMAGE=/boot/vmlinuz-5.4.0-42-generic root=UUID=fe845787-abaa-43af-88b1-f6e8c1fd6b7c ro mitigations=off transparent_hugepage=never
Aug 12 13:00:57 host kernel: KERNEL supported cpus:
Aug 12 13:00:57 host kernel:   Intel GenuineIntel
Aug 12 13:00:57 host kernel:   AMD AuthenticAMD
Aug 12 13:00:57 host kernel:   Hygon HygonGenuine
Aug 12 13:00:57 host kernel:   Centaur CentaurHauls
Aug 12 13:00:57 host kernel:   zhaoxin   Shanghai
Aug 12 13:00:57 host kernel: x86/fpu: Supporting XSAVE feature 0x001: 'x87 floating point registers'
Aug 12 13:00:57 host kernel: x86/fpu: Supporting XSAVE feature 0x002: 'SSE registers'
Aug 12 13:00:57 host kernel: x86/fpu: Supporting XSAVE feature 0x004: 'AVX registers'
Aug 12 13:00:57 host kernel: x86/fpu: xstate_offset[2]:  576, xstate_sizes[2]:  256
Aug 12 13:00:57 host kernel: x86/fpu: Enabled xstate features 0x7, context size is 832 bytes, using 'c

```

### Container Logs
By default, container emit logs to the `stdout` and `stderr` output streams. Containers are ephemeral, however, hence the logs are stored on the host. What takes care of that is a logging driver. You may have noticed the `containers/` folder in `/var/log` directory. It actually points in our case to `/var/log/pods`.

```shell
root@host:~# tree /var/log/pods/
/var/log/pods/
├── kube-system_cilium-2tl6f_1b5be83f-5de0-40a2-8e28-250cfcdb2c2d
│   ├── c06ce52cdcf060a2847263cc8320c8c2d4c4e41728ca98c90f0df5971486d7ee.log
│   ├── cilium-agent
│   │   ├── 0.log
│   │   └── 1.log
│   ├── clean-cilium-state
│   │   └── 1.log
│   └── da5a6993dec12950863d31301f2e4586e749580dd435594516ea763528efc6d4.log
├── kube-system_coredns-dff8fc7d-v4d97_96cab748-5233-490a-9c77-b07b4ac622c0
│   ├── 7a51946f855568f3817f9f43361475cd9dcd2b2820c95be93c7985f96b7ceb5d.log
│   └── coredns
│       └── 0.log
…
```

Which resembles the Kubernetes pods:
```shell
roothost:~# kubectl get pods -A -o json | jq '.items[] | select(.spec.nodeName=="host") | .metadata.name'
"cilium-2tl6f"
"coredns-dff8fc7d-v4d97"
…
```

### Gathering Logs
In a real-world scenario, we usually want to keep logs safe in a centralised storage so that we could control access to them, analyse them or even use them as a knowledge base.

There are some third-party solutions that can help you achieve that. Probably the most commonly used one is the [**ELK stack**](https://www.elastic.co/log-monitoring). The Elastic Stack (sometimes known as the ELK Stack) is the most popular open source logging platform. The ELK is a short for Elasticsearch, Logstash and Kibana, which are the three main components of the stack.

## 3. Monitoring
One of the best ways to gain this insight is with a robust monitoring system that gathers metrics, visualises data, and alerts operators when things appear to be broken.

> What we call **monitoring** here should for better clarity be called "metrics gathering and dashboarding", if we want to break the whole concept to three pillars (as name of this lesson suggest). However, in practice, we use the term "monitoring" for anything related to logging, metrics gathering, dashboarding and alerting. So, next time you're asked to build a proper monitoring infrastructure, make sure not to forget about logging and alerting!

### Parts of a monitoring system
**1. Distributed Monitoring Agents and Data Exporters**
 * a small application designed to collect and forward data to a collection endpoint
 * run as always-on daemons on each host (or a service or sidecar as part of a deployed containerized application) throughout the system
 * the agent must use minimal resources and be able to operate with little to no management
 * typically collect generic, host-level metrics, but agents to monitor software like web or database servers are available as well

**2. Metrics Ingress**
 * one of the busiest part of a monitoring system at any given time
 * for push-based systems, the metrics ingress endpoint is a central location on the network where each monitoring agent or stats aggregator sends its collected data
 * for pull-based systems, the corresponding component is the polling mechanism that reaches out and parses the metrics endpoints exposed on individual hosts
 * the metrics gathering process must be able to provide the correct credentials to log in and access the secure endpoint

>The corresponding component differs for **push-based** and **pull-based** systems. For push-based systems, the metrics ingress endpoint is a central location on the network where each monitoring agent or stats aggregator sends its collected data. In case of the push-based system, the ingress is under heavy load and needs to be load-balanced and replicated to guarantee the system stability.

> In **pull-based** monitoring systems, individual hosts are responsible for gathering, aggregating, and serving metrics in a known format at an accessible endpoint. The monitoring server polls the metrics endpoint on each host to gather the metrics data. The software that collects and presents the data through the endpoint has many of the same requirements as an agent, but often requires less configuration since it does not need to know how to access other machines.


**3. Data Management Layer (DML)**
 * responsible for organizing and recording incoming data from the metrics ingress component and responding to queries and data requests from the administrative layers
* for persistence over longer periods of time, the storage layer needs to provide a way to export data when the collection exceeds the local limitations for processing, memory, or storage

**4. Visualisation and Dashboard Layer**
 * on top of the data management layer are the interfaces that you interact with to understand the data being collected
 * commonly used graphs and data are often organized into saved dashboards
 * a dashboard with a detailed breakdown of physical storage capacity throughout a fleet can be important when capacity planning, but might not need to be referenced for daily administration

**5. Alerting and Threshold Functionality**
 * to reliably notify operators when data indicates an important change and to leave them alone otherwise
 * user-defined metric thresholds
 * important to find a balance between sufficient alerting and over-alerting
 * various notification channels

<p align="center">
  <img src="/assets/prometheus-architecture.png" alt="Prometheus Architecture">
</p>

### Metrics
> **Metrics** represent the raw measurements of resource usage or behavior that can be observed and collected throughout your systems. These might be low-level usage summaries provided by the operating system, or they can be higher-level types of data tied to the specific functionality or work of a component, like requests served per second or membership in a pool of web servers.

**1. Host-based and system metrics**
 * CPU/GPU
 * memory usage

**2. Storage metrics**
 * disk space
 * throughput on R/W

**3. Network metrics**
 * Connectivity
 * Latency
 * Throughput
 * Packet loss

**4. Application metrics**
 * /healthz
 * errors
 * number of requests per second
 * resource usage

#### Tracing
Tracing is useful for monitoring and troubleshooting microservices-based distributed systems, including:

 * Distributed context propagation
 * Distributed transaction monitoring
 * Root cause analysis
 * Service dependency analysis
 * Performance / latency optimisation

<p align="center">
  <img src="/assets/jaeger-2-opt.png" alt="Jaeger">
</p>

### Dashboards

One of the ways to monitor metrics and to share them with a broader audience is via Dashboards. These usually depend on the kind of metrics that you have. For Prometheus metrics there is Grafana (although Elastic can serve as a data source as well), whereas — as I've mentioned before — as part of the ELK stack, there is Kibana. Both of them provide a free version which gives you plenty of possibilities to visualize your metrics. On top of it, for popular open-source projects, there are already dashboards which you can download and use with your applications.

I recommend you to take a look at the [show-off](https://www.elastic.co/kibana) that **Kibana** offers and check out the examples of dashboards you can create with it.

<p align="center">
  <img src="/assets/kibana-dashboard-opt.png" alt="Kibana Dashboard">
</p>

**Grafana** development is also not stale, and it's up to your taste which one you prefer, Grafana or Kibana, both will get the job done. I recommend checking out the official documentation to learn more about [Grafana](https://grafana.com/docs/grafana/latest/getting-started/what-is-grafana/). It is a free yet powerful tool in your hands!

<p align="center">
  <img src="/assets/grafana-dashboard-opt.png" alt="Grafana Dashboard">
</p>

## 4. Alerting
While monitoring systems are incredibly useful for active interpretation and investigation, one of the primary benefits of a complete monitoring system is letting administrators disengage from the system. Alerts allow you to define situations that make sense to actively manage, while relying on the passive monitoring of the software to watch for changing conditions[0].

> Alerting is the responsive component of a monitoring system that performs actions based on changes in metric values. Alerts definitions are composed of two components: a metrics-based condition or threshold, and an action to perform when the values fall outside the acceptable conditions.

### Alerts
> An alert is an action that is triggered when a certain condition is met.

Dashboards are great for actively checking system metrics. However, in practice we don't want to waste (human) resources on that. Instead, we would set up an alerting system. An **alert** is an action that is triggered when a certain condition is met. A condition might be a carefully selected threshold, or a level which, when reached, triggers an event (alert).

Alerts are distributed and processed via **notifications** and **webhooks**.

**Notifications** are messages which are meant to bring attention to the event that triggered the alert. There are various communication channels which might be useful depending on the use case

A **webhook** is an endpoint exposed by your application which consumes an event (which may or may not already be an alert itself). Webhooks are primarily used to trigger automatic responses, like horizontal and vertical autoscaling, cluster scaling, or even air conditioning.

### Best practices
**1. Pick the right metrics**
**2. Choose the correct thresholds and severity**
 * Critical in situations that demand immediate human intervention
 * Secondary notifications for tickets and emails

**Effective alerts:**
 * alerts are triggered by events with real user impact
 * alerts contain appropriate context
 * alerts are sent to the right people
 * thresholds are set with graduated severity

<p align="center">
  <img src="/assets/kibana-alerti.png" alt="Kibana Alerting">
</p>

#### Pagerduty
<p align="center">
  <img src="/assets/pagerduty.png" alt="Pagerduty">
</p>

#### Alerta
<p align="center">
  <img src="/assets/alerta-opt.png" alt="Alerta">
</p>