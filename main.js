"use strict";

const alert = require("./modules/alert");
const child_process = require("child_process");
const electron = require("electron");
const fs = require('fs');
const ipcMain = require("electron").ipcMain;
const readline = require("readline");
const windows = require("./modules/windows");

const STDERR_LOG_WINDOW_ID = -1

electron.app.on("ready", () => {
	menu_build();
	main();
});

electron.app.on("window-all-closed", () => {
	electron.app.quit();
});

function main() {

	// Communications with the compiled app....................................

	let exe = child_process.spawn("game.exe");

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

		if (j.command === "update") {
			windows.update(j.content);
		}

		if (j.command === "alert") {
			alert(j.content);
		}
	});

	windows.new_window({
		uid: STDERR_LOG_WINDOW_ID,
		page: "log.html",
		name: "Stderr",
		width: 500,
		height: 500,
		resizable: true,
		nomenu: true
	});

	let stderr_scanner = readline.createInterface({
		input: exe.stderr,
		output: undefined,
		terminal: false
	});

	stderr_scanner.on("line", (line) => {
		windows.update({
			uid: STDERR_LOG_WINDOW_ID,
			msg: line + "\n"
		});
	});

	// Messages from the renderer..............................................

	ipcMain.on("keydown", (event, msg) => {

		let windobject = windows.get_windobject_from_event(event);

		if (windobject === undefined) {
			return
		}

		let output = {
			type: "key",
			content: {
				down: true,
				uid: windobject.uid,
				key: msg.key
			}
		};

		exe.stdin.write(JSON.stringify(output) + "\n");
	});

	ipcMain.on("keyup", (event, msg) => {

		let windobject = windows.get_windobject_from_event(event);

		if (windobject === undefined) {
			return
		}

		let output = {
			type: "key",
			content: {
				down: false,
				uid: windobject.uid,
				key: msg.key
			}
		};

		exe.stdin.write(JSON.stringify(output) + "\n");
	});

	ipcMain.on("mousedown", (event, msg) => {

		let windobject = windows.get_windobject_from_event(event);

		if (windobject === undefined) {
			return
		}

		let output = {
			type: "mouse",
			content: {
				down: true,
				uid: windobject.uid,
				x: msg.x,
				y: msg.y
			}
		}

		exe.stdin.write(JSON.stringify(output) + "\n");
	});

	ipcMain.on("mouseup", (event, msg) => {

		let windobject = windows.get_windobject_from_event(event);

		if (windobject === undefined) {
			return
		}

		let output = {
			type: "mouse",
			content: {
				down: false,
				uid: windobject.uid,
				x: msg.x,
				y: msg.y
			}
		}

		exe.stdin.write(JSON.stringify(output) + "\n");
	});

	ipcMain.on("request_resize", (event, opts) => {
		let windobject = windows.get_windobject_from_event(event);
		windows.resize(windobject, opts);
	});

	ipcMain.on("ready", (event, opts) => {
		let windobject = windows.get_windobject_from_event(event);
		windows.handle_ready(windobject, opts);
	});
}

function menu_build() {
	const template = [
		{
			label: "Menu",
			submenu: [
				{
					role: "quit"
				},
				{
					type: "separator"
				},
				{
					role: "toggledevtools"
				}
			]
		}
	];

	const menu = electron.Menu.buildFromTemplate(template);
	electron.Menu.setApplicationMenu(menu);
}
