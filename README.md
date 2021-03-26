# cronbackup
crontab incremental backup database tools

1. cronbackup periodically backs up the database by calling different back-end scripts
2. cronbackup will take incremental backups and pack them into compressed packages, so you can set up a daily backup plan so that you can restore the data to any day's state
3. cronbackup will handle the errors in the backup and restore it in the picture above during the next backup to ensure that the backup will not be damaged by accidents such as a power outage during preparation.
4. Currently mariabackup and xtrabackup are supported, but you can define your own backend script to support other databases

# How To
```
$ ./cronbackup backup -h
crontab incremental backup

Usage:
  cronbackup backup [flags]

Flags:
  -b, --backend string    backend script (default "backend/mariadb.js")
  -c, --contab strings    minute hour dom mon dow
  -d, --description       generate description.json
  -h, --help              help for backup
  -H, --host string       hostname for connecting to the server (default "localhost")
  -i, --immediate         perform a backup immediately
  -o, --output string     output path (default "output")
  -p, --password string   password for connecting to the server
  -P, --port uint16       port for connecting to the server (default 3306)
  -u, --user string       username for connecting to the server (default "root")
```

## Execute immediately and exit after the backup is complete

```
./cronbackup backup --host=127.0.0.1 --port=3306 -i -d
```

## Perform a backup every night at 3 o'clock

```
./cronbackup backup --host=127.0.0.1 --port=3306 --contab="0 3 * * *" -d
```

## Execute immediately and Perform a backup every night at 3 o'clock

```
./cronbackup backup --host=127.0.0.1 --port=3306 --contab="0 3 * * *" -i -d
```

# backend

cronbackup is only responsible for scheduling work. Operations such as backup are provided by third-party tools. The back-end script defines how to call these tools.

mariadb.js and mysql.js are predefined back-end scripts, you can refer to them to implement your own back-end scripts.

the backend script is a js file, it is compatible with nodejs modules, but it does not include nodejs runtime. the script only provides the following modules:

```
declare module "os" {
    function exec(name: string, ...args: Array<string>): void
    function cwdExec(cwd: string, name: string, ...args: Array<string>): void
    function join(...args: Array<string>): string
    function readFile(filename: string): string
}
interface Metadata {
    ID: number

    Host: string
    Port: number

    Username: string
    Password: string

    Output: string
}
```