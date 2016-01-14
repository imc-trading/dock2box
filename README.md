# In development

Currently we're focusing on migrating to a new backend, once this is finished we should have a more complete version of this project.

# What?

Dock2Box is a tool that allows you to apply Docker or File images to Bare-Metal hardware using PXE boot.

# Why?

## Late vs early binding

Traditional host provisioning uses tools such as Cobbler, Kickstart and Configuration Management such as SaltStack, Ansible, Puppet or Chef to
finalize the host configuration. This means most of the complexity is late binding i.e. while it's being provisioned.

By moving the complexity to the build step on a CI server, you move most of the complexity where it will fail early rather then late.

Obviously not everything can be early binding, but you can keep the things that are late binding to a minimum.

## Faster

At our current setup we can provision a fully configured server in less then 3 min. the main limitation is bandwidth and the size of the host image.

## Idempotent (repeatable, predictable)

Since we're dealing with a host image it's much more repeatable and predictable.

## Testable

One of the main issues with Configuration Management is testability, it's very hard to test since each deployment is slightly different.

With images they are basically immutable i.e. easier to test.

## DevOps

DevOps is all about breaking the barriers in-between Dev. and Ops., this start's by adopting the same toolchain.

In using Docker for Host provisioning you can use the same toolchain for Software deployment and Host provisioning.

## Running software on Bare-Metal

There are situations when you want to run software on Bare-Metal, but keep the deployment process close to what you
already do with Docker. Now it's just a matter of changing your Base Image "FROM ..." in your Dockerfile and rebuild.

# Overview

![Host Provisioning Overview](img/host_provisioning.png?raw=true)
