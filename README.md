[![Build Status](https://travis-ci.org/Spirals-Team/docker-machine-driver-g5k.svg)](https://travis-ci.org/Spirals-Team/docker-machine-driver-g5k)
[![License](https://img.shields.io/github/license/Spirals-Team/docker-machine-driver-g5k.svg)](https://opensource.org/licenses/Apache-2.0)

# docker-machine-driver-g5k
A Docker Machine driver for the Grid5000 testbed infrastructure. It can be used to provision a Docker machine on a node of the Grid5000 infrastructure.

## Requirements
* [Docker](https://www.docker.com/products/overview#/install_the_platform)
* [Docker Machine](https://docs.docker.com/machine/install-machine)
* [Go tools (Only for installation from sources)](https://golang.org/doc/install)

You need a Grid5000 account to use this driver. See [this page](https://www.grid5000.fr/mediawiki/index.php/Grid5000:Get_an_account) to create an account.

## VPN
**You need to be connected to the Grid5000 VPN to create and access your Docker node.**  
**Do not forget to configure your DNS or use OpenVPN DNS auto-configuration.**  
**Please follow the instructions from the [Grid5000 Wiki](https://www.grid5000.fr/mediawiki/index.php/VPN).**

## Installation from GitHub releases
Binary releases for Linux, MacOS and Windows using x86/x86_64 CPU architectures are available in the [releases page](https://github.com/Spirals-Team/docker-machine-driver-g5k/releases).
On Linux and MacOS, you can use the following commands to install or upgrade the driver:
```bash
# download the binary for your OS and CPU architecture :
sudo curl -L -o /usr/local/bin/docker-machine-driver-g5k "<link to release>"

# grant execution rigths to the driver for everyone :
sudo chmod +x /usr/local/bin/docker-machine-driver-g5k
```

## Installation from sources
To use the Go tools, you need to set your [GOPATH](https://golang.org/doc/code.html#GOPATH) variable environment.

To get the code and compile the binary, run:
```bash
go get -u github.com/Spirals-Team/docker-machine-driver-g5k
```

Then, either put the driver in a directory filled in your PATH environment variable, or run:
```bash
export PATH=$PATH:$GOPATH/bin
```

## How to use

### Driver-specific command line flags

#### Flags description
* **`--g5k-username` : Your Grid5000 account username (required)**
* **`--g5k-password` : Your Grid5000 account password (required)**
* **`--g5k-site` : Site where the reservation of the node will be made (required)**
* `--g5k-walltime` : Duration of the resource reservation (in "HH:MM:SS" format)
* `--g5k-image` : Name of the system image to deploy on the node
* `--g5k-resource-properties` : [Resource selection with OAR properties](#resource-properties)
* `--g5k-make-resource-reservation` : [Make a resource reservation for the given start date](#resource-reservation)
* `--g5k-use-resource-reservation` : [Use a resource reservation (need to be an existing job ID)](#resource-reservation)
* `--g5k-skip-vpn-checks` : Skip the VPN client connection and DNS configuration checks
* `--g5k-reuse-ref-environment` : [Reuse the Grid'5000 reference environment instead of re-deploying the node](#grid5000-reference-environment-reuse)
* `--g5k-job-queue` : [Specify the job queue (besteffort queue is NOT supported)](#job-queues)
* `--g5k-external-ssh-public-keys` : SSH public key(s) allowed to connect to the node (in authorized_keys format)

#### Flags usage
|             Option                |          Environment            |     Default value     |
|-----------------------------------|---------------------------------|-----------------------|
| `--g5k-username`                  | `G5K_USERNAME`                  |                       |
| `--g5k-password`                  | `G5K_PASSWORD`                  |                       |
| `--g5k-site`                      | `G5K_SITE`                      |                       |
| `--g5k-walltime`                  | `G5K_WALLTIME`                  | "1:00:00"             |
| `--g5k-image`                     | `G5K_IMAGE`                     | "debian9-x64-std"     |
| `--g5k-resource-properties`       | `G5K_RESOURCE_PROPERTIES`       |                       |
| `--g5k-make-resource-reservation` | `G5K_MAKE_RESOURCE_RESERVATION` |                       |
| `--g5k-use-resource-reservation`  | `G5K_USE_RESOURCE_RESERVATION`  |                       |
| `--g5k-skip-vpn-checks`           | `G5K_SKIP_VPN_CHECKS`           | False                 |
| `--g5k-reuse-ref-environment`     | `G5K_REUSE_REF_ENVIRONMENT`     | False                 |
| `--g5k-job-queue`                 | `G5K_JOB_QUEUE`                 | "default"             |
| `--g5k-external-ssh-public-keys`  | `G5K_EXTERNAL_SSH_PUBLIC_KEYS`  |                       |

#### Resource properties
You can use [OAR properties](http://oar.imag.fr/docs/2.5/user/usecases.html#using-properties) to only select a node that matches your hardware requirements.  
If you give incorrect properties or no resource matches your request, you will get this error:
```bash
...
Error with pre-create check: "Error when submitting new job: The server returned an error (code: 400) after sending Job submission: '400 Bad Request'"
```

More information about usage of OAR properties are available on the [Grid'5000 Wiki](https://www.grid5000.fr/mediawiki/index.php/Advanced_OAR#Other_examples_using_properties).

#### Resource reservation
You can either do a job submission to reserve resources as soon as possible (this is the default mode) or do an advance reservation for a specific date/time.

To do a resource reservation, you need to use the `--g5k-make-resource-reservation` flag and provide a starting date/time in either the 'YYYY-MM-DD HH:MM:SS' date format or an UNIX timestamp.  
Don't forget to save the job ID of your reservation in order to be able to create a machine when the resources will be available.

To use a resource reservation, set the `--g5k-use-resource-reservation` flag with the job ID of an existing reservation.  
This will create a machine, deploy an OS image and provision Docker on the node.  
Please note that the job must be in `running` state in order for the machine to be created, otherwise the driver will wait until the job start.

More information about the resources reservation are available on the [Grid'5000 Wiki](https://www.grid5000.fr/w/Grid5000:UsagePolicy#Resources_reservation).

#### Grid'5000 reference environment reuse
You can reuse the Grid'5000 reference environment (debian 9, stretch) instead of redeploying the machine.  
Doing so will skip the node deployment phase and will save a lot of time at the machine creation.  
If you don't need a tweaked environment or rely on Grid'5000 services (NFS for example), you should use this option.  
**Please note that, in this mode, if you reboot the machine by any mean, the reserved resource will be released and the node will be redeployed with the reference environment.**  

#### Job queues
You can specify the job queue of your reservation and access the resources of the production queue.  
The driver only support `default`, `production` and `testing` queues. The `besteffort` queue is **NOT** supported.  
If you use an incorrect queue for your site you will get the following error:
```bash
...
Error with pre-create check: "Error when submitting new job: The server returned an error (code: 500) after sending Job submission: '500 Internal Server Error'"
```

See [this page](https://www.grid5000.fr/mediawiki/index.php/Grid5000:UsagePolicy#Rules_for_the_production_queue) for more information about the production queue.

### Usage examples
An example of node provisioning reusing the Grid'5000 standard environment:
```bash
docker-machine create -d g5k \
--g5k-username "user" \
--g5k-password "********" \
--g5k-site "lille" \
--engine-storage-driver "overlay2" \
--g5k-reuse-ref-environment \
test-node
```

An example of node provisioning deploying the `debian9-x64-min` environment on the node:
```bash
docker-machine create -d g5k \
--g5k-username "user" \
--g5k-password "********" \
--g5k-site "lille" \
--g5k-image "debian9-x64-min" \
test-node
```

An example of node provisioning using environment variables:
```bash
export G5K_USERNAME="user"
export G5K_PASSWORD="********"
export G5K_SITE="lille"
docker-machine create -d g5k test-node
```

An example with resource properties (node in cluster `chimint` with more thant 8GB of RAM and at least 4 CPU cores):
```bash
docker-machine create -d g5k \
--g5k-username "user" \
--g5k-password "********" \
--g5k-site "lille" \
--g5k-resource-properties "cluster = 'chimint' and memnode > 8192 and cpucore >= 4" \
test-node
```

An example of doing a resource reservation of 1 node for `8 hours` starting the `2019-01-01` at `20:00:00`:
```bash
docker-machine create -d g5k \
--g5k-username "user" \
--g5k-password "********" \
--g5k-site "lille" \
--g5k-walltime "08:00:00" \
--g5k-make-resource-reservation "2019-01-01 20:00:00" \
test-node
``` 

An example of using a resource reservation:
```bash
docker-machine create -d g5k \
--g5k-username "user" \
--g5k-password "********" \
--g5k-site "lille" \
--g5k-use-resource-reservation 1234567 \
test-node
``` 

An example adding multiple external SSH keys (your keys can be of any supported format, and may be longer than the following example):
```bash
docker-machine create -d g5k \
--g5k-username "user" \
--g5k-password "********" \
--g5k-site "lille" \
--g5k-external-ssh-public-keys "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIFLs3JzUYn7LbHE+SzJNoMvYbasnhjlen0k6dFs801DT test-ed25519" \
--g5k-external-ssh-public-keys "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAAAgQC5qQt/nzGW19uCb9CDVEvP93LZ2mu3rd7drPP1nLf1pzLwlL2U2ksfwDCjMWU0P7KA6tB4scI+4dhxj07t0Z8g4TsMGYhbG0kjf7tWN73DombB4h/zobo2GvVoMg0NBLTP4peXLYAEofTYc0g7OWtJicAzLwcMzHsitDjjBwCKHQ== test-rsa"
```
