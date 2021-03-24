"use strict";
var __spreadArrays = (this && this.__spreadArrays) || function () {
    for (var s = 0, i = 0, il = arguments.length; i < il; i++) s += arguments[i].length;
    for (var r = Array(s), k = 0, i = 0; i < il; i++)
        for (var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)
            r[k] = a[j];
    return r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.finish = exports.removeExpired = exports.pack = exports.checkChanged = void 0;
var os_1 = require("os");
function checkChanged(filename) {
    var str = os_1.readFile(filename);
    var strs = str.split("\n");
    for (var i = 0; i < strs.length; i++) {
        var str_1 = strs[i].trim();
        console.log(str_1);
    }
    return true;
}
exports.checkChanged = checkChanged;
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
