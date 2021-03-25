import { exec, cwdExec, join, readFile } from "os";
export function isNotChanged(filename: string): boolean {
    const str = readFile(filename)
    const strs = str.split("\n")
    const keys = new Map<string, number>()
    for (let i = 0; i < strs.length; i++) {
        let str = strs[i].trim()
        const index = str.indexOf('=')
        if (index < 0) {
            continue
        }
        const key = str.substring(0, index).trim()
        if (keys.has(key)) {
            continue
        }
        if (key != 'from_lsn' && key != 'to_lsn' && key != 'last_lsn') {
            continue
        }
        const n = str.substring(index + 1).trim()
        try {
            const val = parseInt(n)
            if (isNaN(val) || !isFinite(val) || val < 0) {
                continue
            }
            keys.set(key, val)
        } catch (e) {
            console.log(e)
        }
    }
    if (keys.size != 3) {
        throw new Error("analyze xtrabackup_checkpoints error");
    }
    const from_lsn = keys.get('from_lsn')
    const to_lsn = keys.get('to_lsn')
    const last_lsn = keys.get('last_lsn')
    return from_lsn === to_lsn && to_lsn === last_lsn
}
export function backup(name: string, md: Metadata) {
    const output = join(md.Output, md.ID.toString())
    console.log('rm', output, '-rf')
    exec('rm', output, '-rf')
    const args = [
        `--user=${md.Username}`,
        `--password=${md.Password}`,
        '--backup',
        `--host=${md.Host}`,
        `--port=${md.Port}`,
        `--target-dir=${output}`,
    ]

    if (md.ID > 1) {
        const incremental = join(md.Output, (md.ID - 1).toString())
        args.push(`--incremental-basedir=${incremental}`)
    }
    const logs = new Array<string>()
    for (let i = 0; i < args.length; i++) {
        if (args[i].startsWith("--user=") || args[i].startsWith("--password=")) {
            continue
        }
        logs.push(args[i])
    }
    console.log(name, ...logs)
    exec(name, ...args)

    if (isNotChanged(join(output, 'xtrabackup_checkpoints'))) {
        console.log('rm', output, '-rf')
        exec('rm', output, '-rf')
        throw new Error(`${md.ID} data not changed`);
    }
}
export function pack(md: Metadata) {
    const id = md.ID

    const dir = join(md.Output, 'pack')
    console.log('mkdir', dir, '-p')
    exec('mkdir', dir, '-p')

    const source = id.toString()
    const cwd = md.Output
    const dest = join('pack', id.toString() + '.tar.gz')
    console.log('rm', dest, '-rf')
    cwdExec(cwd, 'rm', dest, '-rf')

    const name = 'tar'
    const args = [
        `-zcvf`,
        dest,
        source,
    ]
    console.log(name, ...args)
    cwdExec(cwd, name, ...args)
}
export function removeExpired(md: Metadata) {
    const id = md.ID
    if (id < 1 + 2) {
        return
    }
    const dest = join(md.Output, (id - 2).toString())
    console.log('rm', dest, '-rf')
    exec('rm', dest, '-rf')
}
export function finish(md: Metadata) {

}
