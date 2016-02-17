#!/bin/bash

error() {
    echo $1 >&2
    exit 1
}

D2B_URL="https://raw.githubusercontent.com/imc-trading/dock2box/master"

set -eu

# Check that we're using CentOS 7.x
grep "CentOS Linux release 7" /etc/redhat-release &>/dev/null || error "This script was meant for CentOS 7"

# Install docker
if which docker &>/dev/null; then
    yum install -y docker
fi

# Install docker-compose
if which docker-compose &>/dev/null; then
    curl -L https://github.com/docker/compose/releases/download/1.5.2/docker-compose-`uname -s`-`uname -m` > /usr/bin/docker-compose
    chmod +x /usr/bin/docker-compose
fi

# Check Docker is running
if docker ps &>/dev/null; then
    systemctl start docker
fi

# Setup dock2box
mkdir /etc/dock2box/
mkdir -p /var/lib/dock2box/data/
curl "${D2B_URL}/docker-compose.yml" >/etc/dock2box/docker-compose.yml
curl "${D2B_URL}/dock2box.service" >/etc/systemd/system/dock2box.service
systemctl enable dock2box

# Setup dock2box-update.sh
curl "${D2B_URL}/scripts/dock2box-update.sh" >/usr/local/bin/dock2box-update.sh
chmod +x /usr/local/bin/dock2box-update.sh

# Setup backup
curl "${D2B_URL}/scripts/dock2box-backup.sh" >/usr/local/bin/dock2box-backup.sh
chmod +x /usr/local/bin/dock2box-backup.sh
curl "${D2B_URL}/scripts/dock2box-backup.service" >/etc/systemd/system/dock2box-backup.service
curl "${D2B_URL}/scripts/dock2box-backup.timer" > /etc/systemd/system/dock2box-backup.timer
systemctl enable dock2box-backup.service
systemctl enable dock2box-backup.timer
