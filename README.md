# consul-backup
Dockerized Consul Backup and Restore tool.

This will use consul-api (Go library) to recursively backup and restore all your
key/value pairs.


```sh
docker run djenriquez/consul-backup

Usage:
  consul-backup [-i IP:PORT] [-t TOKEN] [--aclbackup] [--restore <filename>]
  consul-backup -h | --help
  consul-backup --version

Options:
  -h --help                          Show this screen.
  --version                          Show version.
  -i, --address=IP:PORT              The HTTP endpoint of Consul [default: 127.0.0.1:8500].
  -t, --token=TOKEN                  An ACL Token with proper permissions in Consul [default: ].
  -a, --aclbackup                    Backup ACLs, does nothing in restore mode. ACL restore not available at this time.
  -r, --restore                      Activate restore mode
```

## Creating backups:
```sh
CONSUL_ADDRESS=10.0.1.10
CONSUL_PORT=8500
BACKUP_FILE_NAME=backup

docker run --rm \
djenriquez/consul-backup \
-i $CONSUL_ADDRESS:$CONSUL_PORT > $BACKUP_FILE_NAME
```

## Restoring backups:
```sh
CONSUL_ADDRESS=10.0.1.10
CONSUL_PORT=8500
RESTORE_FILE_NAME=backup

docker run --rm \
-v `pwd`:/restore \
djenriquez/consul-backup \
-i $CONSUL_ADDRESS:$CONSUL_PORT $RESTORE_FILE_NAME
```