"use strict";

const electron = require("electron");

function object_to_string(o) {
    let msg = "{\n"
    let keys = Object.keys(o);
    for (let n = 0; n < keys.length; n++) {
        let key = keys[n];
        let val = o[key];
        msg += `   ${key}: ${typeof(val)}` + "\n"
    }
    msg += "}"
    return msg;
}

function alert_main(msg) {
    electron.dialog.showMessageBox({
        message: msg.toString(),
        title: "Alert",
        buttons: ["OK"]
    });
}

function alert_renderer(msg) {
    electron.remote.dialog.showMessageBox({
        message: msg.toString(),
        title: "Alert",
        buttons: ["OK"]
    });
}

module.exports = (msg) => {
    if (typeof(msg) === "object") {
        msg = object_to_string(msg);
    }
    if (process.type === "renderer") {
        alert_renderer(msg);
    } else {
        alert_main(msg);
    }
}
