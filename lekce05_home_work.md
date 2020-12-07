# Setting up automated Docker image build

> * [Task 0: Prerequisites](#Task0)
> * [Task 1: Link to a GitHub user account](#Task1)
> * [Task 2: Configure automated build settings](#Task2)
> * [Task 3: Testing of automated build](#Task3)

### <a name="Task0"></a>Task 0: Prerequisites
Create new repository under your GitHub account.

1. Open [GitHub new](https://github.com/new).

2. Add name & description.

  > Note: For name try to use lower case and no spaces.

3. Keep Public option.

4. Tick "Initialize this repository with a README".

5. Click "Create repository".

6. Login to [Docker Hub](https://hub.docker.com/).

7. Create new repository.

8. Add name & description.

   > Note: It might be good idea used same name as you used for GitHub repository.

9. Click create.

### <a name="Task1"></a>Task 1: Link to a GitHub user account

To automate building and testing of your images, you can link your GitHub account to Docker Hub.

1. Log in to Docker Hub using your Docker ID.

2. Click Account Settings in the top-right dropdown navigation, then open Linked Accounts.

3. Click Connect for the source provider you want to link.

4. Review the settings for the Docker Hub Builder OAuth application.

5. Click Authorize docker** to save the link.



### <a name="Task2"></a>Task 2: Configure automated build settings

You can configure repositories in Docker Hub so that they automatically
build an image each time you push new code to your source provider.

1. From the Repositories section, click into a repository which you created in [previous steps](#Task1) to view its details.

2. Click the Builds tab.

3. If you are setting up automated builds for the first time, select
the code repository service (GitHub) where the image's source code is stored.

    Otherwise, if you are editing the build settings for an existing automated
    build, click Configure automated builds.

4. Select GitHub as the source repository to build the Docker images from.

5. Select your GitHub account and repository which you created in [previous steps](#Task0).

6. Review the default Build Rules.

    Build rules control what Docker Hub builds into images from the contents
    of the source code repository, and how the resulting images are tagged
    within the Docker repository.

    A default build rule is set up for you, which you can edit or delete. This
    default set builds from the Branch in your source code repository called
    master, and creates a Docker image tagged with latest.

7. Click Save to save the settings.

    A webhook is automatically added to your source code repository to notify
    Docker Hub on every push. Only pushes to branches that are listed as the
    source for one or more tags trigger a build.


### <a name="Task3"></a>Task 3: Testing of automated build

Push your code and build automatically your image from Docker file.

1. Clone your Github repository which you created in [previous steps](#Task0).
    Open terminal and run

    ```
    $ git clone git@github.com:_youraccount_/_yourrepo_.git
    $ cd _yourrepo_
    $ vim Dockerfile
    ```

    Past below and save your Dockerfile
    ```
    FROM alpine
    RUN apk add --update \
	       busybox-extras
    CMD ["echo", "Hello World"]
    ```
2. Push your changes to Github

    In your terminal run
    ```
    $ git add .
    $ git commit -m "Add Dockerfile"
    $   [master 2164a0f] Add Dockerfile
    $    1 file changed, 4 insertions(+)
    $    create mode 100644 Dockerfile
    $ git push
    $   Enumerating objects: 4, done.
    $   Counting objects: 100% (4/4), done.
    $   Delta compression using up to 16 threads
    $   Compressing objects: 100% (3/3), done.
    $   Writing objects: 100% (3/3), 365 bytes | 365.00 KiB/s, done.
    $   Total 3 (delta 0), reused 0 (delta 0)
    $   To github.com:ivoklimsa/docker_test.git
    $     ee3c0fd..2164a0f  master -> master
    ```

    Open your GitHub reposiry and check that your Dockerfile has been uploaded.

3. Open your Docker Hub repository and click on Builds tab.
   You should see that push to GitHub, trigger build of Docker Image.

4. Once your build is finished. Anyone should be able to spin container from your image.

    ```
    $ docker run _youraccount_/_yourrepo_:latest
    ```
    > Note: You have to reference your DockerHub account & repository.
