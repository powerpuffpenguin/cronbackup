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