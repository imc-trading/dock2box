#!/bin/bash

set -eu

echo "Disabling firewalld..."
systemctl disable firewalld
systemctl stop firewalld

