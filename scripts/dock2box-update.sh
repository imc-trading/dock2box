#!/bin/bash

set -eu

dock2box-backup.sh
cd /etc/dock2box
systemctl stop dock2box
docker-compose pull
systemctl start dock2box
