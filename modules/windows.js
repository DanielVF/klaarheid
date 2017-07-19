"use strict";

const alert = require("./alert");
const assert = require("assert");
const electron = require("electron");
const ipcMain = require("electron").ipcMain;
const url = require("url");

let windobjects = Object.create(null);		// {win, config, ready, queue}

function get_windobject_from_event(event) {
	for (let uid in windobjects) {
		let val = windobjects[uid];
		if (val.win.webContents === event.sender) {
			return val;
		}
	}
}

function get_win_from_event(event) {
	for (let uid in windobjects) {
		let val = windobjects[uid];
		if (val.win.webContents === event.sender) {
			return val.win;
		}
	}
}

function get_config_from_event(event) {
	for (let uid in windobjects) {
		let val = windobjects[uid];
		if (val.win.webContents === event.sender) {
			return val.config;
		}
	}
}

function resize(win, opts) {
	if (win) {
		win.setContentSize(opts.xpixels, opts.ypixels);
	}
};

// Commands from launcher.js...

exports.new = (config) => {

	assert(windobjects[config.uid] === undefined);

	let win = new electron.BrowserWindow({
		title: config.name,
		width: 200,
		height: 200,
		backgroundColor: "#000000",
		useContentSize: true,
		resizable: config.resizable
	});

	win.loadURL(url.format({
		protocol: "file:",
		pathname: config.page,
		slashes: true
	}));

	win.on("closed", () => {
		// TODO
	});

	windobjects[config.uid] = {win: win, config: config, ready: false, queue: []};
};

function maybe_add_to_queue(windobject, retry_func) {
	if (windobject.ready === false) {
		windobject.queue.push(retry_func);
		return true;
	}
	return false;
}

exports.flip = (content) => {

	let windobject = windobjects[content.uid];
	if (maybe_add_to_queue(windobject, () => {exports.flip(content);})) {
		return;
	}

	let win = windobjects[content.uid].win;
	win.webContents.send("flip", content);
}

// Messages from a window...

ipcMain.on("resize", (event, opts) => {
	let win = get_win_from_event(event);
	resize(win, opts);
});

ipcMain.on("ready", (event, opts) => {

	let windobject = get_windobject_from_event(event);
	windobject.ready = true;

	let config = get_config_from_event(event);
	event.sender.send("init", config);

	// Now resend things which we queued up because the window wasn't ready.
	// This might not be good enough (because the init message is async?)
	// therefore the renderer also implements its own ability to delay
	// a flip message if needed.

	for (let n = 0; n < windobject.queue.length; n++) {
		windobject.queue[n]();
	}
	windobject.queue = [];
});
