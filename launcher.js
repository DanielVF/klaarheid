"use strict";

const alert = require("./modules/alert");
const child_process = require("child_process");
const electron = require("electron");
const fs = require('fs');
const readline = require("readline");
const windows = require("./modules/windows");

electron.app.on("ready", () => {
	main();
});

electron.app.on("window-all-closed", () => {
	electron.app.quit();
});

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

function main() {

	let exe = child_process.spawn("test.exe");

	let scanner = readline.createInterface({
		input: exe.stdout,
		output: undefined,
		terminal: false
	});

	scanner.on("line", (line) => {
		let j = JSON.parse(line);

		if (j.command === "new") {
			windows.new(j.content);
		}

		// Other messages can fail if the window isn't ready...

		if (j.command === "flip") {
			windows.flip(j.content);
		}
	});
}
