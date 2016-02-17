#!/bin/bash

error() {
    echo $1 >&2
    exit 1
}

set -eu

# Check that we're using CentOS 7.x
grep "CentOS Linux release 7" /etc/redhat-release &>/dev/null || error "This script was meant for CentOS 7.x"

# Check Docker and Docker Compose
which docker &>/dev/null || error "Missing dependency: docker"
which docker-compose &>/dev/null || error "Missing dependency: docker-compose"
which etcdtool &>/dev/null || error "Missing dependency: etcdtool"

# Check Docker is running
docker ps &>/dev/null || error "Docker isn't running"

# Setup backend
mkdir /etc/dock2box/
mkdir -p /var/lib/dock2box/data/
cp ./docker-compose.yml /etc/dock2box/
cp ./dock2box.service /etc/systemd/system/
systemctl enable dock2box

# Setup update
cp ./dock2box-update.sh /usr/local/bin/
chmod +x /usr/local/bin/dock2box-update.sh

# Setup backup
cp ./dock2box-backup.sh /usr/local/bin/
cp ./dock2box-backup.service /etc/systemd/system/
cp ./dock2box-backup.timer /etc/systemd/system/
systemctl enable dock2box-backup.service
systemctl enable dock2box-backup.timer
