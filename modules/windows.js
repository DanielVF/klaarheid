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

	assert(config.uid !== undefined);
	assert(windobjects[config.uid] === undefined);

	let win_pixel_width = config.width;
	let win_pixel_height = config.height;

	// The config may or may not specify width and height in terms of a grid of boxes, with each box taking up a certain size...

	if (config.boxwidth !== undefined && config.boxheight !== undefined) {
		win_pixel_width *= config.boxwidth;
		win_pixel_height *= config.boxheight;
	}

	let win = new electron.BrowserWindow({
		title: config.name,
		width: win_pixel_width,
		height: win_pixel_height,
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

	if (config.nomenu === true) {
		win.setMenu(null);
	}

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

function update(content) {
	let windobject = windobjects[content.uid];
	send_or_queue(windobject, "update", content);
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
exports.update = update;
exports.handle_ready = handle_ready;
exports.resize = resize;
