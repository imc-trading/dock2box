# In development

## Run dock2box

## Setup script

If you are running CentOS 7 there is a setup script.

```bash
curl https://raw.githubusercontent.com/imc-trading/dock2box/master/setup_centos7.sh | bash -
```

## Vagrant

You can start a VM using Vagrant, this will provision a CentOS 7 with the setup script.

```bash
vagrant up
```

## Docker Compose

First comment out or change the etcd data path.

```bash
docker-compose up
```

## Manual instructions

These are the steps in order to run dock2box.

```bash
mkdir /etc/dock2box/
mkdir -p /var/lib/dock2box/data/
cp <checkout>/docker-compose.yml /etc/dock2box/
cp <checkout>/dock2box.service /etc/systemd/system/
systemctl start dock2box
systemctl enable dock2box
```

## Update dock2box

Setup update script.

```bash
cp <checkout>/dock2box-update.sh /usr/local/bin/
chmod +x /usr/local/bin/dock2box-update.sh
```

Run update script.

```bash
dock2box-update.sh
```

## Setup backup

Backup is setup using a SystemD timer. Backup relies on [etcdtool](https://github.com/mickep76/etcdtool).

```bash
cp <checkout>/dock2box-backup.sh /usr/local/bin/
cp <checkout>/dock2box-backup.service /etc/systemd/system/
cp <checkout>/dock2box-backup.timer /etc/systemd/system/
systemctl start dock2box-backup.service
systemctl start dock2box-backup.timer
systemctl enable dock2box-backup.service
systemctl enable dock2box-backup.timer
```

# Docker Components

## cfg (imctrading/etcdrest-config:latest)

This contains the configuration, schemas and templates for etcdrest.

For more information look in the [cfg](cfg/) folder.

## api (mickep76/etcdrest:latest)

This contains the software etcdrest to run a REST API with etcd as a database.

For more informaton look in [etcdrest](https://github.com/mickep76/etcdrest).

## etcd (quay.io/coreos/etcd:latest)

This contains etcd database from CoreOS.

## tftp (imctrading/dock2box-tftp-dev:latest)

This contains TFTP daemon that serves "undionly.kpxe" iPXE.

For more information look in the [tftp](tftp/) folder.

# What?

Dock2Box is a tool that allows you to apply Docker or File images to bare-metal hardware using PXE boot.

# Why?

## Late vs early binding

Traditional host provisioning uses tools such as Cobbler, Kickstart and Configuration Management tools such as SaltStack, Ansible, Puppet or Chef to
finalize the host configuration. This means most of the complexity is late-binding i.e. while the host is being provisioned.

By moving the complexity to the build step on a CI server, you move most of the complexity where it will fail earlier rather then later.

Obviously not everything can be early-binding, but you can keep the things that are late-binding to a minimum and thus uncover errors sooner.

## Faster

At our current setup we can provision a fully configured server in less than 3 min. The main limitation is bandwidth and the size of the host image.

## Idempotent (repeatable, predictable)

Since we're dealing with a host image it's much more repeatable and predictable.

## Testable

One of the main issues with Configuration Management is testability: it's very hard to test since each deployment is slightly different.

With images they are basically immutable i.e. easier to test.

## DevOps

DevOps is all about breaking the barriers in-between Dev and Ops. This starts by adopting the same toolchain.

In using Docker for host provisioning you can use the same toolchain for software deployment and host provisioning.

## Running software on Bare-Metal

There are situations when you want to run software on bare-metal, but keep the deployment process close to what you
already do with Docker. Now it's just a matter of changing your Base Image "FROM ..." in your Dockerfile and rebuild.

# Workflow

![Workflow](img/workflow.png?raw=true)

# Overview

![Overview](img/overview.png?raw=true)

# PXE menu

![PXE menu](img/pxe_menu.png?raw=true)
