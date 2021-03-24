import { exec, cwdExec, join, readFile } from "os";
export function checkChanged(filename: string): boolean {
    const str = readFile(filename)
    const strs = str.split("\n")
    for (let i = 0; i < strs.length; i++) {
        let str = strs[i].trim();
        console.log(str)
    }
    return true
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
