export { pack, removeExpired, finish } from "./mariadb-mysql";
import { backup as _backup } from "./mariadb-mysql";
export function backup(md: Metadata) {
    _backup('xtrabackup', md)
}