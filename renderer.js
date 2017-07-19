"use strict";

const alert = require("./modules/alert");
const child_process = require("child_process");
const fs = require("fs");
const ipcRenderer = require("electron").ipcRenderer;
const readline = require("readline");

function id_from_xy(x, y) {
	return `td_${x}_${y}`;
}

function make_renderer() {

	let renderer = {
		ready_to_flip: false,
		postponed_flip_message: null,
	};

	renderer.init = (opts) => {
		let x;
		let y;
		let id;
		let html = `<table style="font-size: ${opts.fontpercent}%;">`;

		renderer.width = opts.width;
		renderer.height = opts.height;
		renderer.uid = opts.uid;

		for (y = 0; y < opts.height; y++) {
			html += "<tr>";
			for (x = 0; x < opts.width; x++) {
				id = id_from_xy(x, y);
				html += `<td id="${id}" style="width: ${opts.boxwidth}; height: ${opts.boxheight};"></td>`;
			}
			html += "</tr>"
		}

		html += "</table>"

		document.getElementsByTagName("body")[0].innerHTML = html;

		renderer.resize(opts.width * opts.boxwidth, opts.height * opts.boxheight);

		renderer.ready_to_flip = true;
		if (renderer.postponed_flip_message) {
			renderer.flip(postponed_flip_message);
			renderer.postponed_flip_message = null;
		}
	};

	renderer.flip = (charstring) => {

		if (renderer.ready_to_flip !== true) {
			renderer.postponed_flip_message = charstring;
			return;
		}

		for (let x = 0; x < renderer.width; x++) {
			for (let y = 0; y < renderer.height; y++) {
				let index = y * renderer.width + x;
				let id = id_from_xy(x, y);
				let element = document.getElementById(id);
				if (element) {
					element.innerHTML = charstring.charAt(index);
				}
			}
		}
	};

	renderer.resize = (xpixels, ypixels) => {
		let msg = {
			xpixels: xpixels,
			ypixels: ypixels,
		};
		ipcRenderer.send("resize", msg);
	};

	return renderer;
}

let renderer = make_renderer();

ipcRenderer.on("init", (event, opts) => {
	renderer.init(opts);
});

ipcRenderer.on("flip", (event, opts) => {
	renderer.flip(opts.chars);
});

ipcRenderer.send("ready", null);	// triggers an init message to be sent to us
