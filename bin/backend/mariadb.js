"use strict";
var __spreadArrays = (this && this.__spreadArrays) || function () {
    for (var s = 0, i = 0, il = arguments.length; i < il; i++) s += arguments[i].length;
    for (var r = Array(s), k = 0, i = 0; i < il; i++)
        for (var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)
            r[k] = a[j];
    return r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.backup = void 0;
var os_1 = require("os");
var mariadb_mysql_1 = require("./mariadb-mysql");
Object.defineProperty(exports, "pack", { enumerable: true, get: function () { return mariadb_mysql_1.pack; } });
Object.defineProperty(exports, "removeExpired", { enumerable: true, get: function () { return mariadb_mysql_1.removeExpired; } });
Object.defineProperty(exports, "finish", { enumerable: true, get: function () { return mariadb_mysql_1.finish; } });
var mariadb_mysql_2 = require("./mariadb-mysql");
function backup(md) {
    var output = os_1.join(md.Output, md.ID.toString());
    console.log('rm', output, '-rf');
    os_1.exec('rm', output, '-rf');
    var name = 'mariabackup';
    var args = [
        "--user=" + md.Username,
        "--password=" + md.Password,
        '--backup',
        "--host=" + md.Host,
        "--port=" + md.Port,
        "--target-dir=" + output,
    ];
    if (md.ID > 1) {
        var incremental = os_1.join(md.Output, (md.ID - 1).toString());
        args.push("--incremental-basedir=" + incremental);
    }
    var logs = new Array();
    for (var i = 0; i < args.length; i++) {
        if (args[i].startsWith("--user=") || args[i].startsWith("--password=")) {
            continue;
        }
        logs.push(args[i]);
    }
    console.log.apply(console, __spreadArrays([name], logs));
    os_1.exec.apply(void 0, __spreadArrays([name], args));
    if (!mariadb_mysql_2.checkChanged(os_1.join(output, 'xtrabackup_checkpoints'))) {
        throw new Error(md.ID + " data not changed");
    }
    throw new Error("test " + md.ID + " data not changed");
}
exports.backup = backup;
