## Our first pipeline

```yaml
default:
  image: ackee/gitlab-builder

services:
  - docker:dind

stages:
  - build
  - test

build my super awesome app:
  stage: build
  script:
    - docker build . -t awesome_app
    - docker save awesome_app > image.tar
  artifacts:
    paths:
      - image.tar
    expire_in: "1 hour"

test my super awesome app:
  stage: test
  script:
    - docker load -i image.tar
    - docker run -d -p 8080:8080 awesome_app  # this might not work
    - curl localhost:8080
  dependencies:
    - build my super awesome app
```

 * Customizing your runtime can result into longer pipeline runtime due to docker pull operations

Let's see what were the problems during development of pipeline:
https://gitlab.com/beranm/simple_app/-/pipelines

 * Why exporting a port on runner might not work?
 * What are the alternatives?

## Gitab CI building blocks

Everything could be found at https://docs.gitlab.com/ce/ci/yaml/ you can start
at https://docs.gitlab.com/ce/ci/quick_start/README.html

### Runtime separation

 * stages - contain jobs (https://docs.gitlab.com/ce/ci/yaml/README.html#stages):
    * Jobs of the same stage are run concurrently,
    * Jobs of the next stage are run after the jobs from the previous stage complete successfully.
 * job - runtime execution, there is no other place where you can execute your commands

### Runners:

https://docs.gitlab.com/runner/ written in golang ;-)

 * Runtime of the jobs
 * can have their own runtime, but!

Possible use-cases:
 * iOS apps need to be build at Mac OS due to legal policies, special runner have to be deployed at Mac OS.

Which runners are picked for a job are managed by tags added to a job:

```yaml
tests:
  stage: tests
  interruptible: true
  script:
    - if [ ! -f Gemfile ]; then fail "Missing Gemfile"; fi
    - bundle config set path ${GEM_CACHE_DIR}/gems
    - bundle install
    - source helper_functions.sh
    - read_branch_config
    - bundle exec fastlane prepare environment:$ENVIRONMENT
    - bundle exec fastlane test type:$TESTS_TYPE
  artifacts:
    reports:
      junit:
        - fastlane/test_output/report.junit
  cache: &cache
    key: ${CI_COMMIT_SHORT_SHA}
    paths:
      - ./Pods
      - ./Carthage
  tags:
    - ios
```

This job will run on runners tagged `ios`.

### Installing your own runner

To install you own runner, you can just:

```bash
$ curl -LJO https://gitlab-runner-downloads.s3.amazonaws.com/latest/deb/gitlab-runner_amd64.deb 
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 57.0M  100 57.0M    0     0   623k      0  0:01:33  0:01:33 --:--:--  987k

# sudo su
# dpkg -i gitlab-runner_amd64.deb 
Selecting previously unselected package gitlab-runner.
(Reading database ... 381224 files and directories currently installed.)
Preparing to unpack gitlab-runner_amd64.deb ...
Unpacking gitlab-runner (13.5.0) ...
Setting up gitlab-runner (13.5.0) ...
GitLab Runner: creating gitlab-runner...
Runtime platform                                    arch=amd64 os=linux pid=28932 revision=ece86343 version=13.5.0
gitlab-runner: Service is not running.
Runtime platform                                    arch=amd64 os=linux pid=28941 revision=ece86343 version=13.5.0
gitlab-ci-multi-runner: Service is not running.
Runtime platform                                    arch=amd64 os=linux pid=28972 revision=ece86343 version=13.5.0
Runtime platform                                    arch=amd64 os=linux pid=29028 revision=ece86343 version=13.5.0
Clearing docker cache...

# gitlab-runner register
Runtime platform                                    arch=amd64 os=linux pid=29391 revision=ece86343 version=13.5.0
Running in system-mode.                            
                                                   
Please enter the gitlab-ci coordinator URL (e.g. https://gitlab.com/):
https://gitlab.ack.ee/
Please enter the gitlab-ci token for this runner:
KCKyKbn3q1uLXWg6xbUb
Please enter the gitlab-ci description for this runner:
[Nemo64]: my_pc  
Please enter the gitlab-ci tags for this runner (comma separated):
home,amd64,mint,linux
Registering runner... succeeded                     runner=KCKyKbn3
Please enter the executor: shell, ssh, virtualbox, kubernetes, docker-ssh, parallels, docker+machine, docker-ssh+machine, custom, docker:
docker
Please enter the default Docker image (e.g. ruby:2.6):
ubuntu:latest
Runner registered successfully. Feel free to start it, but if it's running already the config should be automatically reloaded! 
```

Token is at your project `Settings -> CI/CD -> Runners`.

Configuration is saved at `/etc/gitlab-runner/config.toml`.

Let's see an example of config, which is responsible for a little bit more:

```
concurrent = 50   # All registered Runners can run up to 50 concurrent jobs

[[runners]]
  name = "blupidupy-tester"
  url = "https://gitlab.ack.ee/"
  token = "YyKNmNAdamELSxJ4hEmp" 
  executor = "docker+machine"
  limit = 10
    [runners.docker]
    tls_verify = false
    image = "ubuntu:latest"
    privileged = true
    disable_entrypoint_overwrite = false
    oom_kill_disable = false
    disable_cache = false
    volumes = ["/cache"]
    shm_size = 0
  [runners.machine]
    IdleCount = 5                    # There must be 5 machines in Idle state - when Off Peak time mode is off
    IdleTime = 600                   # Each machine can be in Idle state up to 600 seconds (after this it will be removed) - when Off Peak time mode is off
    MaxBuilds = 100                  # Each machine can handle up to 100 jobs in a row (after this it will be removed)
    MachineName = "auto-scale-%s"    # Each machine will have a unique name ('%s' is required)
    MachineDriver = "google" # Refer to Docker Machine docs on how to authenticate: https://docs.docker.com/machine/drivers/gce/#credentials
    MachineOptions = [
      "google-project=terraform-test-hejda",
      "google-zone=europe-west1-d",
      "google-machine-type=n1-standard-1",
      "google-machine-image=ubuntu-os-cloud/global/images/family/ubuntu-1804-lts",
      "google-username=root",
      "google-use-internal-ip",
      "engine-registry-mirror=https://mirror.gcr.io"
    ]
    [[runners.machine.autoscaling]]  # Define periods with different settings
      Periods = ["* * 9-17 * * mon-fri *"] # Every workday between 9 and 17 UTC
      IdleCount = 50
      IdleTime = 3600
      Timezone = "UTC"
    [[runners.machine.autoscaling]]
      Periods = ["* * * * * sat,sun *"] # During the weekends
      IdleCount = 5
      IdleTime = 60
      Timezone = "UTC"
  [runners.cache]
    Type = "gcs"
    [runners.cache.gcs]
      CredentialsFile = "/etc/gitlab-runner/service-account.json"
      BucketName = "beranm-test"
```

It's actually adjusted example from
https://docs.gitlab.com/runner/configuration/autoscale.html#a-complete-example-of-configtoml

It does:
 * scaling of runners in Google Compute Engine,
 * save cache data to Google Cloud Storage bucket.

### Sharing options in a pipeline:

Variables are not shared in between the jobs, only files.

 * artifacts - sharing trough http transfer to gitlab server
    * browsable at gitlab server
 * cache - sharing on whenever: AWS S3 bucket, GCS, local file system
    * saving something into cache in one job does not always mean it will be in the next job

 * workspaces (not in gitlab ci) - Jenkins 
 
### The power of variables

See https://docs.gitlab.com/ee/ci/variables/

Precedence https://docs.gitlab.com/ee/ci/variables/#priority-of-environment-variables 

The order of precedence for variables is (from highest to lowest):
 * Trigger variables, scheduled pipeline variables, and manual pipeline run variables.
 * Project-level variables or protected variables.
 * Group-level variables or protected variables.
 * Instance-level variables or protected variables.
 * Inherited environment variables.
 * YAML-defined job-level variables.
 * YAML-defined global variables.
 * Deployment variables.
 * Predefined environment variables. 


 * Define constraints related to project
  * Deployment numbers, hashes of commits, ...

 * `CI_DEBUG_TRACE` - if set `true` you'll see all the variables passed to jobs

#### Example in our pipeline

##### Cache

Beware, cache content do not need to be there every time you need it:
https://docs.gitlab.com/ee/ci/caching/#availability-of-the-cache


```yaml
default:
  image: ackee/gitlab-builder

services:
  - docker:dind

variables:
  - DOCKER_APP_NAME: awesome_app

stages:
  - build
  - test

build my super awesome app:
  stage: build
  script:
    - docker build . -t $DOCKER_APP_NAME
    - docker save $DOCKER_APP_NAME > image.tar
  cache:
    paths:
      - image.tar

test my super awesome app:
  stage: test
  script:
    - |
      if [ ! -f image.tar ]; then
        # paranoid cache check of image
        docker build . -t $DOCKER_APP_NAME
        docker save $DOCKER_APP_NAME > image.tar
      fi
    - docker load -i image.tar
    - docker run -d -p 8080:8080 $DOCKER_APP_NAME  # this might not work
    - curl localhost:8080
  dependencies:
    - build my super awesome app
```

Do not forget: Be careful if you use cache and artifacts to store the same path in your jobs as caches are restored
before artifacts and the content could be overwritten. 
https://docs.gitlab.com/ee/ci/caching/#cache-vs-artifacts

## Let us fix that check

Let's create an infrastructure in docker-compose:

```yaml
version: "3.3"

services:
  app:
    build:
      context: ./
      dockerfile: ./Dockerfile

  client:
    image: curlimages/curl
    command: sleep 600
```

and we change the pipeline to:

```yaml
default:
  image: ackee/gitlab-builder

services:
  - docker:dind

stages:
  - build
  - test

build my super awesome app:
  stage: build
  script:
    - docker build . -t awesome_app
    - docker save awesome_app > image.tar
  cache:
    paths:
      - image.tar

test my super awesome app:
  stage: test
  script:
    - docker-compose up -d
    - docker exec simple_app_client_1 curl http://app:8080/hello
    - docker-compose rm -s -f
  dependencies:
    - build my super awesome app
```

## Docker in docker and it's consequences

See https://jpetazzo.github.io/2015/09/03/do-not-use-docker-in-docker-for-ci/

All our pipelines are DIND. Gitlab does not all mounting `/var/run/docker.sock`, we use service description for that:
```
services:
  - docker:dind
```
Try to remove the services and run pipeline again.

Reasons why not to:
 * it takes much resources, every image is pulled with new pipeline
 * unrelated file system issues (which never happened to me)
 
Reasons why you should consider it:
 * You do not have permission to mount `/var/run/docker.sock`
 * Your runtime is unique
 * Speed is not an issue
 
Mounting `/var/run/docker.sock` is not using docker in docker, see following example:
```
cd /root
touch secret_file.txt
echo "secret" > secret_file.txt 
cd /tmp 
mkdir test
echo "test test, testy test" > test/file
docker run -ti -v '/var/run/docker.sock:/var/run/docker.sock' -v '/tmp:/container_tmp' docker:latest sh
cat /container_tmp/dind_test/test/file 
test test, testy test
```
And now run docker container mounting the volume `/container_tmp` in docker:
```
docker run -ti -v '/container_tmp:/another_container_tmp' docker:latest sh
ls another_container_tmp/
```
But, when we do following:
```
docker run -ti -v '/tmp:/another_container_tmp' docker:latest sh
cat another_container_tmp/dind_test/test/file 
test test, testy test
```
And even worse:
```
docker run -ti -v '/root:/root' docker:latest sh
cat root/secret_file.txt 
secret
```
DO NOT MOUNT `/var/run/docker.sock` UNLESS YOU HAVE NOTHING TO WORRY ABOUT

or you already set correct tools like apparmor or SELinux

Let's run the example as DIND:
```
docker run -ti -v '/tmp:/container_tmp' docker:latest sh
cat /container_tmp/dind_test/test/file 
test test, testy test
docker run -ti -v '/container_tmp:/another_container_tmp' docker:latest sh
docker: error during connect: Post http://docker:2375/v1.40/containers/create: dial tcp: lookup docker on 10.0.1.5:53: no such host.
```
That's why we have:
```
services:
  - docker:dind
```
in the pipeline.

We can:
```
docker run -ti -e DOCKER_TLS_CERTDIR="" --privileged docker:dind sh
dockerd-entrypoint.sh &
# will show a lot of output
echo '127.0.0.1 docker' >> /etc/hosts

mkdir /tmp/mount_into_container
echo "test content" > /tmp/mount_into_container/test_file
docker run -v '/tmp/mount_into_container/:/tmp/mount_into_container/' ubuntu cat /tmp/mount_into_container/test_file
test content
docker run -v '/tmp/mount_into_container/:/tmp/mount_into_container/' ubuntu cat /root/secret_file.txt
cat: /root/secret_file.txt: No such file or directory
```
