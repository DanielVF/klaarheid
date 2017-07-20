"use strict";

const alert = require("./modules/alert");
const child_process = require("child_process");
const fs = require("fs");
const ipcRenderer = require("electron").ipcRenderer;
const readline = require("readline");

// --------------------------------------------------------------

const colour_dict = JSON.parse(fs.readFileSync("colours.json"));

// --------------------------------------------------------------

function id_from_xy(x, y) {
	return `td_${x}_${y}`;
}

// --------------------------------------------------------------

function make_renderer() {

	let renderer = {
		ready_to_flip: false,
		postponed_actions: [],
	};

	renderer.init = (opts) => {

		renderer.width = opts.width;
		renderer.height = opts.height;
		renderer.uid = opts.uid;

		// Make the table...

		let html = `<table style="font-size: ${opts.fontpercent}%;">`;

		for (let y = 0; y < opts.height; y++) {
			html += "<tr>";
			for (let x = 0; x < opts.width; x++) {
				let id = id_from_xy(x, y);
				html += `<td id="${id}" style="width: ${opts.boxwidth}; height: ${opts.boxheight};"></td>`;
			}
			html += "</tr>"
		}

		html += "</table>"

		document.getElementsByTagName("body")[0].innerHTML = html;

		renderer.resize(opts.width * opts.boxwidth, opts.height * opts.boxheight);

		// Set the "ready" flag and then deal with any postponed items that were queued
		// due to us not being ready earlier...

		renderer.ready_to_flip = true;

		for (let n = 0; n < renderer.postponed_actions.length; n++) {
			renderer.postponed_actions[n]();
		}
		renderer.postponed_actions = [];

		// Input handlers...

		document.addEventListener("keydown", (evt) => {
			ipcRenderer.send("keydown", {key: evt.key});
		});

		document.addEventListener("keyup", (evt) => {
			ipcRenderer.send("keyup", {key: evt.key});
		});

		for (let x = 0; x < renderer.width; x++) {
			for (let y = 0; y < renderer.height; y++) {
				let id = id_from_xy(x, y);
				let element = document.getElementById(id)
				element.addEventListener("mousedown", (evt) => {
					evt.preventDefault();								// Prevent selecting text with the mouse
					ipcRenderer.send("mousedown", {x: x, y: y});		// x and y work despite closures because we use "let" in the loops
				});
				element.addEventListener("mouseup", (evt) => {
					ipcRenderer.send("mouseup", {x: x, y: y});			// Works despite closures because we use "let" in the loops
				});
			}
		}
	};

	renderer.flip = (opts) => {

		if (renderer.ready_to_flip !== true) {
			renderer.postponed_actions.push(() => renderer.flip(opts));
			return;
		}

		let charstring = opts.chars;
		let colourstring = opts.colours;
		let highlight = opts.highlight;

		for (let x = 0; x < renderer.width; x++) {
			for (let y = 0; y < renderer.height; y++) {
				let index = y * renderer.width + x;
				let id = id_from_xy(x, y);
				let element = document.getElementById(id);
				if (element) {
					element.innerHTML = charstring.charAt(index);
					let colour_key = colourstring.charAt(index);
					let colour = colour_dict[colour_key];
					if (colour) {
						element.style["color"] = colour;
					}
					if (x === highlight.x && y === highlight.y) {
						element.style["background-color"] = "#333333";
					} else {
						element.style["background-color"] = "black";
					}
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

// --------------------------------------------------------------

let renderer = make_renderer();

// --------------------------------------------------------------

ipcRenderer.on("init", (event, opts) => {
	renderer.init(opts);
});

ipcRenderer.on("flip", (event, opts) => {
	renderer.flip(opts);
});

// --------------------------------------------------------------

ipcRenderer.send("ready", null);	// triggers an init message to be sent to us
