import { exec, join } from "os";
export { pack, removeExpired, finish } from "./mariadb-mysql";
import { isNotChanged } from "./mariadb-mysql";
export function backup(md: Metadata) {
    const output = join(md.Output, md.ID.toString())
    console.log('rm', output, '-rf')
    exec('rm', output, '-rf')
    const name = 'mariabackup'
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
        throw new Error(`${md.ID} data not changed`);
    }
}
