## Let's improve what he have already

### Packaging & push should be bound to a tag:

```yaml
push my super awesome app:
  stage: push
  scripts:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
    - docker load -i image.tar
    - docket tag awesomeApp $CI_REGISTRY_IMAGE:$CI_COMMIT_TAG
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_TAG
  dependencies:
    - test my super awesome app
  rules:
    - if: $CI_COMMIT_TAG
```

Sometimes, you can see instead of `rules` directives `only` & `except`, see https://docs.gitlab.com/ee/ci/yaml/#onlyexcept-basic
those are not recommended anymore.

* Utilize gitlab registry, dockerhub have strict limits, see 
  https://www.docker.com/blog/scaling-docker-to-serve-millions-more-developers-network-egress/
* Gitlab supports also package registry: npm, pypi, Maven, ...

### Let's pipeline only if something changed

```yaml
build my super awesome app:
  stage: build
  scripts:
    - docker build . -t awesomeApp
    - docker save awesomeApp > image.tar
  artifacts:
    paths:
      - image.tar
    expire_in: "1 hour"
  rules:
    - changes:
        - main.go
```

### Deploy the application into K8s

You will study k8s in next lecture, let's take this as blackbox.

```yaml
deploy my super awesome app:
  stage: deploy
  script:
    - kubectl set image -n default deployment/deploy-0 deploy-0-container=$CI_REGISTRY_IMAGE:$CI_COMMIT_TAG
  rules:
    - if: $CI_COMMIT_TAG
  environment:
    name: dev
    url: https://example.com
```

#### Environments 

https://docs.gitlab.com/ee/ci/environments/

Environments allow control of the continuous deployment of your software, all within GitLab.

It's very helpful once you need to:
 * deploy the app for QA team,
 * run (smoke, integration, ...) tests,
 * show the application at demo.

Once you are done, define `on_stop` action to stop the environment:

```yaml
deploy my super awesome app:
  stage: deploy
  script:
    - kubectl set image -n default deployment/deploy-0 deploy-0-container=$CI_REGISTRY_IMAGE:$CI_COMMIT_TAG
  rules:
    - if: $CI_COMMIT_TAG
  environment:
    name: dev
    url: https://example.com
    on_stop: stop_env

stop_env:
  script:
    - kubectl delete deployment -n default deploy-0
  when: manual
  environment:
    name: https://example.com
    action: stop
```

## What have we learned

 * What's CI/CD?
 * What are the basic stages of build, test, deploy
 * Gitlab
 * Runners and their meaning
 * Artifacts vs Caches

## What's next

* Kubernetes integration
* Packages in project
* DAG

### DAG - Directed Acyclic Graph 

 * Stages are barriers, once job in a stage is running, you need to wait
 * What if you can define relations between the jobs
 * Run the job once the related job finnish 

```yaml
stages:
  - build
  - test
  - deploy

arm build:
  stage: build
  script: sleep 10 && echo Done

amd64 build:
  stage: build
  script: sleep 20 && echo Done

arm tests:
  stage: test
  needs:
    - arm build
  script: echo Done

amd64 tests:
  stage: test
  needs:
    - amd64 build
  script: sleep 10 && echo Done

deploy:
  stage: test
  needs:  # try to remove build steps
    - arm build
    - amd64 build
    - amd64 tests
    - arm tests
  script: echo Done

```

## Unrelated

Create javascript exploit on our service.
