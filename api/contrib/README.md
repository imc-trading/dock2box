# Setup service on CentOS 7

```bash
mkdir /etc/dock2box-api
cp ../../docker-compose.yml /etc/dock2box-api
cp dock2box-api.service /etc/systemd/system
```

# Setup daily backup

Copy script.

```bash
cp dock2box-daily.sh /usr/local/bin/dock2box-daily.sh
```

Add the following to crontab.

```
# Daily backup of Dock2Box API
0 0 * * * /usr/local/bin/dock2box-daily.sh &>/var/log/dock2box-api-daily.log
```
