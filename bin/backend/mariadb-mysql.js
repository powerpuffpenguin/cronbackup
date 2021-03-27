"use strict";
var __spreadArrays = (this && this.__spreadArrays) || function () {
    for (var s = 0, i = 0, il = arguments.length; i < il; i++) s += arguments[i].length;
    for (var r = Array(s), k = 0, i = 0; i < il; i++)
        for (var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)
            r[k] = a[j];
    return r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.finish = exports.removeExpired = exports.pack = exports.backup = void 0;
var os_1 = require("os");
function loadLSN(filename) {
    var _a, _b, _c;
    var str = os_1.readFile(filename);
    var strs = str.split("\n");
    var keys = new Map();
    for (var i = 0; i < strs.length; i++) {
        var str_1 = strs[i].trim();
        var index = str_1.indexOf('=');
        if (index < 0) {
            continue;
        }
        var key = str_1.substring(0, index).trim();
        if (keys.has(key)) {
            continue;
        }
        if (key != 'from_lsn' && key != 'to_lsn' && key != 'last_lsn') {
            continue;
        }
        var n = str_1.substring(index + 1).trim();
        try {
            var val = parseInt(n);
            if (isNaN(val) || !isFinite(val) || val < 0) {
                continue;
            }
            keys.set(key, val);
        }
        catch (e) {
            console.log(e);
        }
    }
    if (keys.size != 3) {
        throw new Error("analyze xtrabackup_checkpoints error");
    }
    return {
        from: (_a = keys.get('from_lsn')) !== null && _a !== void 0 ? _a : 0,
        to: (_b = keys.get('to_lsn')) !== null && _b !== void 0 ? _b : 0,
        last: (_c = keys.get('last_lsn')) !== null && _c !== void 0 ? _c : 0,
    };
}
function backup(name, md) {
    var output = os_1.join(md.Output, md.ID.toString());
    console.log('rm', output, '-rf');
    os_1.exec('rm', output, '-rf');
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
    var current = loadLSN(os_1.join(output, 'xtrabackup_checkpoints'));
    if (current.from == current.to) {
        console.log('rm', output, '-rf');
        os_1.exec('rm', output, '-rf');
        throw new Error(md.ID + " data not changed");
    }
}
exports.backup = backup;
function pack(md) {
    var id = md.ID;
    var dir = os_1.join(md.Output, 'pack');
    console.log('mkdir', dir, '-p');
    os_1.exec('mkdir', dir, '-p');
    var source = id.toString();
    var cwd = md.Output;
    var dest = os_1.join('pack', id.toString() + '.tar.gz');
    console.log('rm', dest, '-rf');
    os_1.cwdExec(cwd, 'rm', dest, '-rf');
    var name = 'tar';
    var args = [
        "-zcvf",
        dest,
        source,
    ];
    console.log.apply(console, __spreadArrays([name], args));
    os_1.cwdExec.apply(void 0, __spreadArrays([cwd, name], args));
}
exports.pack = pack;
function removeExpired(md) {
    var id = md.ID;
    if (id < 1 + 2) {
        return;
    }
    var dest = os_1.join(md.Output, (id - 2).toString());
    console.log('rm', dest, '-rf');
    os_1.exec('rm', dest, '-rf');
}
exports.removeExpired = removeExpired;
function finish(md) {
}
exports.finish = finish;
