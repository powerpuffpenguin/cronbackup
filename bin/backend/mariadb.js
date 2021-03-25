"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.backup = void 0;
var mariadb_mysql_1 = require("./mariadb-mysql");
Object.defineProperty(exports, "pack", { enumerable: true, get: function () { return mariadb_mysql_1.pack; } });
Object.defineProperty(exports, "removeExpired", { enumerable: true, get: function () { return mariadb_mysql_1.removeExpired; } });
Object.defineProperty(exports, "finish", { enumerable: true, get: function () { return mariadb_mysql_1.finish; } });
var mariadb_mysql_2 = require("./mariadb-mysql");
function backup(md) {
    mariadb_mysql_2.backup('mariabackup', md);
}
exports.backup = backup;
