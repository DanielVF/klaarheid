"use strict";

const alert = require("./alert");
const assert = require("assert");
const electron = require("electron");
const ipcMain = require("electron").ipcMain;
const url = require("url");

// The windobject is our fundamental object, containing fields:
//			{uid, win, config, ready, queue}

let windobjects = Object.create(null);		// dict: uid --> windobject

function get_windobject_from_event(event) {
	for (let uid in windobjects) {
		let val = windobjects[uid];
		if (val.win.webContents === event.sender) {
			return val;
		}
	}
	return undefined;
}

function resize(windobject, opts) {
	if (windobject) {
		windobject.win.setContentSize(opts.xpixels, opts.ypixels);
	}
};

function new_window(config) {

	assert(windobjects[config.uid] === undefined);

	let win = new electron.BrowserWindow({
		title: config.name,
		width: config.width * config.boxwidth,
		height: config.height * config.boxheight,
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
		delete windobjects[config.uid];
	});

	windobjects[config.uid] = {
		uid: config.uid,
		win: win,
		config: config,
		ready: false,
		queue: [],

		send: (channel, msg) => {
			win.webContents.send(channel, msg);
		}
	};
};

function flip(content) {
	let windobject = windobjects[content.uid];
	send_or_queue(windobject, "flip", content);
}

function send_or_queue(windobject, channel, msg) {
	if (windobject === undefined) {
		return;
	}
	if (windobject.ready !== true) {
		windobject.queue.push(() => windobject.send(channel, msg))
		return;
	}
	windobject.send(channel, msg);
}

function handle_ready(windobject, opts) {

	if (windobject === undefined) {
		return;
	}

	windobject.ready = true;

	let config = windobject.config;
	windobject.send("init", config);

	// Now carry out whatever actions were delayed because the window wasn't ready...

	for (let n = 0; n < windobject.queue.length; n++) {
		windobject.queue[n]();
	}

	windobject.queue = [];
}


exports.get_windobject_from_event = get_windobject_from_event;
exports.new_window = new_window;
exports.flip = flip;
exports.handle_ready = handle_ready;
exports.resize = resize;
