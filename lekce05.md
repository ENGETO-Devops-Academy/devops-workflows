# Docker Essentials

## 1. Introduction and Setup
In this lab we'll setup access for DockerID and LinuxVMs.


> * [Task 1.1: Prerequisites](#Task1.1)
> * [Task 1.2: Simple Docker container](#Task1.2)
> * [Task 1.3: Building simple Docker container](#Task1.3)


> Note: Docker CLI Cheatsheet  https://devhints.io/docker

### <a name="Task1.1"></a>Task 1.1: Prerequisites

Before we start, you'll need to install Docker or use pre-defined Docker instance, clone a GitHub repo, and make sure you have a DockerID.

### Make sure you have a DockerID

If you do not have a DockerID (a free login used to access Docker Cloud, Docker Store, and Docker Hub), please visit [Docker Cloud](https://cloud.docker.com) to register for one.


### How to install Docker

>* [Linux](https://docs.docker.com/engine/install/)
>* [Windows](https://docs.docker.com/docker-for-windows/install/)
>* [Mac](https://docs.docker.com/docker-for-mac/install/)
>* [Without installation](https://hybrid.play-with-docker.com/)

### <a name="Task1.2"></a>Task 1.2: Run some simple Docker containers

In this section you'll deploy simple Docker container.

### Run a single-task Alpine Linux container

In this step we're going to start a new container and tell it to run the `hostname` command. The container will start, execute the `hostname` command, then exit.

1. Run the following command in your Linux console:

    ```
    $ docker container run alpine hostname
    Unable to find image 'alpine:latest' locally
    latest: Pulling from library/alpine
    88286f41530e: Pull complete
    Digest: sha256:f006ecbb824d87947d0b51ab8488634bf69fe4094959d935c0c103f4820a417d
    Status: Downloaded newer image for alpine:latest
    888e89a3b36b
    ```

    The output above shows that the `alpine:latest` image could not be found locally. When this happens, Docker automatically *pulls* it form Docker Hub.

    > Note: After the image is pulled, the container's hostname is displayed (`888e89a3b36b` in the example above).

### <a name="Task1.3"></a>Task 1.3: Build some simple Docker container
In this step we're going to crete your own Docker image from Docker file. In the end we'll run it.

1. Create file name 'Dockerfile'

    ```
    $ mkdir MyFirstDocker
    $ touch MyFirstDocker/Dockerfile
    ```

2. Use VIM to edit 'MyFirstDocker/Dockerfile'

    ```
    $ vim MyFirstDocker/Dockerfile
    ```

3. Let's add below to your Dockerfile

    ```
    FROM alpine
    MAINTAINER your@email.com
    RUN apk add --update \
	    busybox-extras
    CMD ["telnet", "towel.blinkenlights.nl"]
    ```

    > Note: Above we are adding base image 'Alpine'. More info about [ALPINE](https://alpinelinux.org/about/) distro.
    As well as latest update and [busybox](https://busybox.net/about.html). Last line is about executing telnet command.

4. Next step would be building and tagging our Docker image.

    ```
    $ cd MyFirstDocker
    $ docker build --tag myfirstdocker:1.0 .
    Sending build context to Docker daemon  2.048kB
    Step 1/4 : FROM alpine
     ---> a24bb4013296
    Step 2/4 : MAINTAINER your@email.com
     ---> Running in e713edd101d4
    Removing intermediate container e713edd101d4
     ---> 8f01ef63b2fd
    Step 3/4 : RUN apk add --update 	busybox-extras
     ---> Running in cd124ec06b4b
    fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/main/x86_64/APKINDEX.tar.gz
    fetch http://dl-cdn.alpinelinux.org/alpine/v3.12/community/x86_64/APKINDEX.tar.gz
    (1/1) Installing busybox-extras (1.31.1-r19)
    Executing busybox-extras-1.31.1-r19.post-install
    Executing busybox-1.31.1-r16.trigger
    OK: 6 MiB in 15 packages
    Removing intermediate container cd124ec06b4b
     ---> c17759bb9b00
    Step 4/4 : CMD ["telnet", "towel.blinkenlights.nl"]
     ---> Running in c12efd733c0e
    Removing intermediate container c12efd733c0e
     ---> f407ef35f1fb
    Successfully built f407ef35f1fb
    Successfully tagged myfirstdocker:1.0
    ```

5. As we now built our first Docker Image, let's make sure that we have it.

    ```
    $ docker image ls
    $ REPOSITORY     TAG    IMAGE ID        CREATED             SIZE
      myfirstdocker  1.0    f407ef35f1fb    19 minutes ago      7.48MB
    ```

6. Ok, so it's look like that we have, let's run it.

    ```
    $ docker run -ti docker run -ti myfirstdocker:1.0
    ```

## 2. Understanding the Docker File System and Volumes

We had an introduction to volumes earlier in , but let's take a practical look at the Docker file system and volumes.

The [Docker documentation](https://docs.docker.com/engine/userguide/storagedriver/imagesandcontainers/_) gives a great explanation of how storage works with Docker images and containers, but here are the high points.

The following exercises will help to illustrate those concepts in practice.

Let's start by looking at layers and how files written to a container are managed by something called *copy on write*.

> * [Task 2.1: Layers and Copy on Write](#Task2.1)
> * [Task 2.2: Anonymous Volumes](#Task2.2)
> * [Task 2.3: Named Volumes](#Task2.3)


### <a name="Task2.1"></a>Task 2.1: Layers and Copy on Write

1. Pull down the Alpine image

    ```
    $ docker pull alpine:3.11
      3.11: Pulling from library/alpine
      e6b0cf9c0882: Pull complete
      Digest: sha256:2171658620155679240babee0a7714f6509fae66898db422ad803b951257db78
      Status: Downloaded newer image for alpine:3.11
    ```

2. Pull down an Alpine example image with added command

    ```
    $ docker pull ivoklimsa/examples:alpine_example
      alpine_example: Pulling from ivoklimsa/examples
      e6b0cf9c0882: Already exists
      9a789c75f73c: Pull complete
      Digest: sha256:8118c555031d6e29eebd22194caf8be7d71d2acac44c576881ae26dc6435a210
      Status: Downloaded newer image for ivoklimsa/examples:alpine_example
    ```

    What do you notice about the output from the Docker pull request for MySQL?

    The first layer pulled says:

    `e6b0cf9c0882: Already exists`

    Notice that the layer id (`e6b0cf9c0882`) is the same for the first layer of the Alpine_example image and the only layer in the Alpine:3.11 image. And because we already had pulled that layer when we pulled the Alpine image, we didn't have to pull it again.

    So, what does that tell us about the Alpine_example image? Since each layer is created by a line in the image's *Dockerfile*, we know that the Alpine_example image is based on the Alpine:3.11 base image. We can confirm this by looking at the [Dockerfile ](https://github.com/ivoklimsa/DockerVSB/blob/master/alpine/Dockerfile).

    The first line in the the Dockerfile is: `FROM alpine:3.11` This will import that layer into the Alpine_example image.

    So layers are created by Dockerfiles and are shared between images. When you start a container, a writeable layer is added to the base image.

### <a name="Task2.2"></a>Task 2.2: Anonymous Volumes

[Docker volumes](https://docs.docker.com/engine/admin/volumes/volumes/) are directories on the host file system that are not managed by the storage driver. Since they are not managed by the storage drive they offer a couple of important benefits.

The next sections will cover both anonymous and named volumes.

1. Pull down a MySQL 5.7 image

    ```
    $ docker pull mysql:5.7
      5.7: Pulling from library/mysql
      804555ee0376: Pull complete
      c53bab458734: Pull complete
      ca9d72777f90: Pull complete
      2d7aad6cb96e: Pull complete
      8d6ca35c7908: Pull complete
      6ddae009e760: Pull complete
      327ae67bbe7b: Pull complete
      31f1f8385b27: Pull complete
      a5a3ad97e819: Pull complete
      48bede7828ac: Pull complete
      380afa2e6973: Pull complete
      Digest: sha256:b38555e593300df225daea22aeb104eed79fc80d2f064fde1e16e1804d00d0fc
      Status: Downloaded newer image for mysql:5.7
    ```

    If you look at the MySQL [Dockerfile](https://github.com/docker-library/mysql/blob/6659750146b7a6b91a96c786729b4d482cf49fe6/5.7/Dockerfile) you will find the following line:

    ```
    VOLUME /var/lib/mysql
    ```

    This line sets up an anonymous volume in order to increase database performance by avoiding sending a bunch of writes through the Docker storage driver.

    > Note: An anonymous volume is a volume that hasn't been explicitly named. This means that it's extremely difficult to use the volume later with a new container. Named volumes solve that problem, and will be covered later in this section.

2. Start a MySQL container

    ```
    $ docker run --name mysqldb -e MYSQL_USER=mysql -e MYSQL_PASSWORD=mysql -e MYSQL_DATABASE=sample -e MYSQL_ROOT_PASSWORD=supersecret -d mysql:5.7
      acf185dc16e274b2f332266a1bfc6d1df7d7b4f780e6a7ec6716b40cafa5b3c3
    ```

    When we start the container, the anonymous volume is created:

3. Use Docker inspect to view the details of the anonymous volume

    ```
    $ docker inspect -f 'in the {{.Name}} container {{(index .Mounts 0).Destination}} is mapped to {{(index .Mounts 0).Source}}' mysqldb
    in the /mysqldb container /var/lib/mysql is mapped to /var/lib/docker/volumes/cd79b3301df29d13a068d624467d6080354b81e34d794b615e6e93dd61f89628/_data
    ```

4. Change into the volume directory on the local host file system and list the contents

    ```
    $ cd $(docker inspect -f '{{(index .Mounts 0).Source}}' mysqldb)

    $ ls
      auto.cnf            ib_buffer_pool      mysql               server-cert.pem
      ca-key.pem          ib_logfile0         performance_schema  server-key.pem
      ca.pem              ib_logfile1         private_key.pem     sys
      client-cert.pem     ibdata1             public_key.pem
      client-key.pem      ibtmp1              sample
    ```

    Notice the directory name starts with `/var/lib/docker/volumes/` whereas for directories managed by the Overlay2 storage driver it was `/var/lib/docker/overlay2`

    As mentined, anonymous volumes will not persist data between containers, they are almost always used to increase performance.

5. Shell into your running MySQL container and log into MySQL

    ```
    $ docker exec --tty --interactive mysqldb bash

      root@132f4b3ec0dc:/# mysql --user=mysql --password=mysql
      mysql: [Warning] Using a password on the command line interface can be insecure.
      Welcome to the MySQL monitor.  Commands end with ; or \g.
      Your MySQL connection id is 3
      Server version: 5.7.19 MySQL Community Server (GPL)

      Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

      Oracle is a registered trademark of Oracle Corporation and/or its
      affiliates. Other names may be trademarks of their respective
      owners.

      Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
    ```

6. Create a new table

    ```
    mysql> show databases;
    +--------------------+
    | Database           |
    +--------------------+
    | information_schema |
    | sample             |
    +--------------------+
    2 rows in set (0.00 sec)

    mysql> connect sample;
    Connection id:    4
    Current database: sample

    mysql> show tables;
    Empty set (0.00 sec)

    mysql> create table user(name varchar(50));
    Query OK, 0 rows affected (0.01 sec)

    mysql> show tables;
    +------------------+
    | Tables_in_sample |
    +------------------+
    | user             |
    +------------------+
    1 row in set (0.00 sec)
    ```

7. Exit MySQL and the MySQL container.

    ```
    mysql> exit
    Bye

    root@132f4b3ec0dc:/# exit
    exit
    ```

8. Stop the container and restart it

    ```
    $ docker stop mysqldb
    mysqldb

    $ docker start mysqldb
    mysqldb
    ```

9. Shell back into the running container and log into MySQL

    ```
    $ docker exec --interactive --tty mysqldb bash

      root@132f4b3ec0dc:/# mysql --user=mysql --password=mysql
      mysql: [Warning] Using a password on the command line interface can be insecure.
      Welcome to the MySQL monitor.  Commands end with ; or \g.
      Your MySQL connection id is 3
      Server version: 5.7.19 MySQL Community Server (GPL)

      Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

      Oracle is a registered trademark of Oracle Corporation and/or its
      affiliates. Other names may be trademarks of their respective
      owners.

      Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
      ```

10. Ensure the table created previously still exists

    ```
    mysql> connect sample;
    Reading table information for completion of table and column names
    You can turn off this feature to get a quicker startup with -A

    Connection id:    4
    Current database: sample

    myslq> show tables;
    +------------------+
    | Tables_in_sample |
    +------------------+
    | user             |
    +------------------+
    1 row in set (0.00 sec)
    ```

11. Exit MySQL and the MySQL container.

    ```
    mysql> exit
    Bye

    root@132f4b3ec0dc:/# exit
    exit
    ```

    The table persisted across container restarts, which is to be expected. In fact, it would have done this whether or not we had actually used a volume as shown in the previous section.

12. Let's look at the volume again

    ```
    $ docker inspect -f 'in the {{.Name}} container {{(index .Mounts 0).Destination}} is mapped to {{(index .Mounts 0).Source}}' mysqldb
      in the /mysqldb container /var/lib/mysql is mapped to /var/lib/docker/volumes/cd79b3301df29d13a068d624467d6080354b81e34d794b615e6e93dd61f89628/_data
    ```

    We do see the volume was not affected by the container restart either.

    Where people often get confused is in expecting that the anonymous volume can be used to persist data BETWEEN containers.

    To examine that delete the old container, create a new one with the same command, and check to see if the table exists.

13. Remove the current MySQL container

    ```
    $ docker container rm --force mysqldb
      mysqldb
    ```

14. Start a new container with the same command that was used before

    ```
    $ docker run --name mysqldb -e MYSQL_USER=mysql -e MYSQL_PASSWORD=mysql -e MYSQL_DATABASE=sample -e MYSQL_ROOT_PASSWORD=supersecret -d mysql:5.7
      eb15eb4ecd26d7814a8da3bb27cee1a23304fab1961358dd904db37c061d3798
    ```

15. List out the volume details for the new container

    ```
    $ docker inspect -f 'in the {{.Name}} container {{(index .Mounts 0).Destination}} is mapped to {{(index .Mounts 0).Source}}' mysqldb
      in the /mysqldb container /var/lib/mysql is mapped to /var/lib/docker/volumes/e0ffdc6b4e0cfc6e795b83cece06b5b807e6af1b52c9d0b787e38a48e159404a/_data
    ```

    Notice this directory is different than before.

16. Shell back into the running container and log into MySQL

    ```
    $ docker exec --interactive --tty mysqldb bash

    root@132f4b3ec0dc:/# mysql --user=mysql --password=mysql
    mysql: [Warning] Using a password on the command line interface can be insecure.
    Welcome to the MySQL monitor.  Commands end with ; or \g.
    Your MySQL connection id is 3
    Server version: 5.7.19 MySQL Community Server (GPL)

    Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

    Oracle is a registered trademark of Oracle Corporation and/or its
    affiliates. Other names may be trademarks of their respective
    owners.

    Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
    ```

17. Check to see if the table created previously still exists

    ```
    mysql> connect sample;
    Connection id:    4
    Current database: sample

    mysql> show tables;
    Empty set (0.00 sec)
    ```

18. Exit MySQL and the MySQL container.

    ```
    mysql> exit
    Bye

    root@132f4b3ec0dc:/# exit
    exit
    ```

19. Remove the container

    ```
    docker container rm --force mysqldb
    mysqldb
    ```

So while a volume was used to store the new table in the original container, because it wasn't a named volume, the data could not be persisted between containers.

To achieve persistence, a named volume should be used.

### <a name="Task2.3"></a>Task 2.3: Named Volumes

A named volume (as the name implies) is a volume that's been explicitly named and can easily be referenced.

A named volume can be created on the command line, in a docker-compose file, and when you start a new container. They [CANNOT be created as part of the image's dockerfile](https://github.com/moby/moby/issues/30647).

1. Start a MySQL container with a named volume (`dbdata`)

    ```
    $ docker run --name mysqldb \
      -e MYSQL_USER=mysql \
      -e MYSQL_PASSWORD=mysql \
      -e MYSQL_DATABASE=sample \
      -e MYSQL_ROOT_PASSWORD=supersecret \
      --detach \
      --mount type=volume,source=mydbdata,target=/var/lib/mysql \
      mysql:5.7
    ```

    Because the newly created volume is empty, Docker will copy over whatever existed in the container at `/var/lib/mysql` when the container starts.

    Docker volumes are simple, just like images and containers. As such, they can be listed and removed in the same way.

2. List the volumes on the Docker host

    ```
    $ docker volume ls
      DRIVER              VOLUME NAME
      local               55c322b9c4a644a5284ccb5e4d7b6b466a0534e26d57c9ef4221637d39cf9a88
      local               cc44059d23e0a914d4390ea860fd35b2acdaa480e83c025fb381da187b652a66
      local               e0ffdc6b4e0cfc6e795b83cece06b5b807e6af1b52c9d0b787e38a48e159404a
      local               mydbdata
    ```

3. Inspect the volume

    ```
    $ docker inspect mydbdata
      [
          {
              "CreatedAt": "2017-10-13T19:55:10Z",
              "Driver": "local",
              "Labels": null,
              "Mountpoint": "/var/lib/docker/volumes/mydbdata/_data",
              "Name": "mydbdata",
              "Options": {},
              "Scope": "local"
          }
      ]
    ```

    Any data written to `/var/lib/mysql` in the container will be rerouted to `/var/lib/docker/volumes/mydbdata/_data` instead.

4. Shell into your running MySQL container and log into MySQL

    ```
    $ docker exec --tty --interactive mysqldb bash

    root@132f4b3ec0dc:/# mysql --user=mysql --password=mysql
    mysql: [Warning] Using a password on the command line interface can be insecure.
    Welcome to the MySQL monitor.  Commands end with ; or \g.
    Your MySQL connection id is 3
    Server version: 5.7.19 MySQL Community Server (GPL)

    Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

    Oracle is a registered trademark of Oracle Corporation and/or its
    affiliates. Other names may be trademarks of their respective
    owners.

    Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
    ```

5. Create a new table

    ```
    mysql> connect sample;
    Connection id:    4
    Current database: sample

    mysql> show tables;
    Empty set (0.00 sec)

    mysql> create table user(name varchar(50));
    Query OK, 0 rows affected (0.01 sec)

    mysql> show tables;
    +------------------+
    | Tables_in_sample |
    +------------------+
    | user             |
    +------------------+
    1 row in set (0.00 sec)
    ```

6. Exit MySQL and the MySQL container.

    ```
    mysql> exit
    Bye

    root@132f4b3ec0dc:/# exit
    exit
    ```

7. Remove the MySQL container

    ```
    $ docker container rm --force mysqldb
    ```

    Because the MySQL was writing out to a named volume, we can start a new container with the same data.

    When the container starts, it will not overwrite existing data in a volume. So the data created in the previous steps will be left intact and mounted into the new container.

8. Start a new MySQL container

    ```
    $ docker run --name new_mysqldb \
      -e MYSQL_USER=mysql \
      -e MYSQL_PASSWORD=mysql \
      -e MYSQL_DATABASE=sample \
      -e MYSQL_ROOT_PASSWORD=supersecret \
      --detach \
      --mount type=volume,source=mydbdata,target=/var/lib/mysql \
      mysql:5.7
    ```

9. Shell into your running MySQL container and log into MySQL

    ```
    $ docker exec --tty --interactive new_mysqldb bash

      root@132f4b3ec0dc:/# mysql --user=mysql --password=mysql
      mysql: [Warning] Using a password on the command line interface can be insecure.
      Welcome to the MySQL monitor.  Commands end with ; or \g.
      Your MySQL connection id is 3
      Server version: 5.7.19 MySQL Community Server (GPL)

      Copyright (c) 2000, 2017, Oracle and/or its affiliates. All rights reserved.

      Oracle is a registered trademark of Oracle Corporation and/or its
      affiliates. Other names may be trademarks of their respective
      owners.

      Type 'help;' or '\h' for help. Type '\c' to clear the current input statement.
    ```

10. Check to see if the previously created table exists in your new container.

    ```
    mysql> connect sample;
    Reading table information for completion of table and column names
    You can turn off this feature to get a quicker startup with -A

    Connection id:    4
    Current database: sample

    mysql> show tables;
    +------------------+
    | Tables_in_sample |
    +------------------+
    | user             |
    +------------------+
    1 row in set (0.00 sec)
    ```

    The data will exist until the volume is explicitly deleted.

11. Exit MySQL and the MySQL container.

    ```
    mysql> exit
    Bye

    root@132f4b3ec0dc:/# exit
    exit
    ```

12. Remove the new MySQL container and volume

    ```
    $ docker container rm --force new_mysqldb
      new_mysqldb

    $ docker volume rm mydbdata
      mydbdata
    ```

    If a new container was started with the previous command, it would create a new empty volume.

### 3. Docker Container Networking

### Networking

Docker supports several different networking options, but this lab will cover only bridge .

Bridge networks are only available on the local host, and can be created on hosts. Bridge networks only work on the host on which they were created.

> * [Task 3.1: Bridge networking](#Task3.1)

### <a name="Task3.1"></a>Task 3.1: Bridge networking

1. Create a bridge network (`mybridge`)

    ```
    $ docker network create mybridge
      52fb9de4ad1cbe505e451599df2cb62c53e56893b0c2b8d9b8715b5e76947551
    ```

2. List the networks on your host

    ```
    $ docker network ls
      NETWORK ID          NAME                DRIVER              SCOPE
      6ce9bcd89540        bridge              bridge              local
      2edacbd4d472        host                host                local
      f6ee43f99cb        mybridge            bridge              local
      1a3768df786f        none                null                local
    ```

    The newly created `mybridge` network is listed.

    > Note: Docker creates several networks by default, however the purpose of those networks is outside the scope of this workshop.

3. Create an Alpine container named `alpine_host` running the `top` process in `detached` mode and connecit it to the `mybridge` network.

    ```
    $ docker container run \
      --detach \
      --network mybridge \
      --name alpine_host \
      alpine top
      Unable to find image 'alpine:latest' locally
      latest: Pulling from library/alpine88286f41530e: Pull complete
      Digest: sha256:f006ecbb824d87947d0b51ab8488634bf69fe4094959d935c0c103f4820a417d
      Status: Downloaded newer image for alpine:latest
      974903580c3e452237835403bf3a210afad2ad1dff3e0b90f6d421733c2e05e6
    ```
    > Note: We run the `top` process to keep the container from exiting as soon as it's created.

4. Start another Alpine container named `alpine client`

    ```
    $ docker container run \
      --detach \
      --name alpine_client \
      alpine top
      c81a3a14f43fed93b6ce2eb10338c1749fde0fe7466a672f6d45e11fb3515536
    ```

5. Attempt to PING `alpine_host` from `alpine_client`

    ```
    $ docker exec alpine_client ping alpine_host
      ping: bad address 'alpine_host'
    ```

    Because the two containers are not on the same network they cannot reach each other.

6. Inspect `alpine_host` and `alpine_client` to see which networks they are attached to.

    ```
    $ docker inspect -f {{.NetworkSettings.Networks}} alpine_host
      map[mybridge:0xc420466000]

    $ docker inspect -f {{.NetworkSettings.Networks}} alpine_client
      map[bridge:0xc4204420c0]
    ```

    `alpine_host` is, as expected, attached to the `mybridge` network.

    `alpine_client` is attached to the default bridge network `bridge`

7. Stop and remove `alpine_client`

    ```
    $ docker container rm --force alpine_client
      alpine_client
    ```

8. Start another container called `alpine_client` but attach it to the `mybridge` network this time.

    ```
    $ docker container run \
      --detach \
      --network mybridge \
      --name alpine_client \
      alpine top
      8cf39f89560fa8b0f6438222b4c5e3fe53bdeab8133cb59038650231f3744a79
    ```

9. Verify via `inspect` that `alpine_client` is on the `mybridge` network

    ```
    $ docker inspect -f {{.NetworkSettings.Networks}} alpine_client
      map[mybridge:0xc42043e0c0]
    ```

10. PING `alpine_host` from `alpine_client`

    ```
    docker exec alpine_client ping -c 5 alpine_host
    PING alpine_host (172.20.0.2): 56 data bytes
    64 bytes from 172.20.0.2: seq=0 ttl=64 time=0.102 ms
    64 bytes from 172.20.0.2: seq=1 ttl=64 time=0.108 ms
    64 bytes from 172.20.0.2: seq=2 ttl=64 time=0.088 ms
    64 bytes from 172.20.0.2: seq=3 ttl=64 time=0.113 ms
    64 bytes from 172.20.0.2: seq=4 ttl=64 time=0.122 ms

    --- alpine_host ping statistics ---
    5 packets transmitted, 5 packets received, 0% packet loss
    round-trip min/avg/max = 0.088/0.106/0.122 ms
    ```

    Something to notice is that it was not necessary to specify an IP address.  Docker has a built in DNS that resolved `alpine_client` to the correct address.
