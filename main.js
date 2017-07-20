"use strict";

const alert = require("./modules/alert");
const child_process = require("child_process");
const electron = require("electron");
const fs = require('fs');
const ipcMain = require("electron").ipcMain;
const readline = require("readline");
const windows = require("./modules/windows");

electron.app.on("ready", () => {
	main();
});

electron.app.on("window-all-closed", () => {
	electron.app.quit();
});

function main() {

	// main() contains stuff that deals with direct communication with the go program.

	let exe = child_process.spawn("test.exe");

	let scanner = readline.createInterface({
		input: exe.stdout,
		output: undefined,
		terminal: false
	});

	scanner.on("line", (line) => {
		let j = JSON.parse(line);

		if (j.command === "new") {
			windows.new_window(j.content);
		}

		// Other messages can fail if the window isn't ready...

		if (j.command === "flip") {
			windows.flip(j.content);
		}
	});

	ipcMain.on("keydown", (event, key) => {

		let windobject = windows.get_windobject_from_event(event);

		if (windobject === undefined) {
			return
		}

		let msg = {
			type: "key",
			content: {
				down: true,
				uid: windobject.uid,
				key: key
			}
		};

		exe.stdin.write(JSON.stringify(msg) + "\n");
	});

	ipcMain.on("keyup", (event, key) => {

		let windobject = windows.get_windobject_from_event(event);

		if (windobject === undefined) {
			return
		}

		let msg = {
			type: "key",
			content: {
				down: false,
				uid: windobject.uid,
				key: key
			}
		};

		exe.stdin.write(JSON.stringify(msg) + "\n");
	});
}
