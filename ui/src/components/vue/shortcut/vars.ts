/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

export const codeMap: Record<string, Array<string>> = {
    "backspace": ["backspace"],
    "tab": ["tab"],
    "enter": ["enter"],
    "escape": ["escape"],
    "space": ["space"],
    "left": ["arrowleft"],
    "up": ["arrowup"],
    "right": ["arrowright"],
    "down": ["arrowdown"],
    "delete": ["delete"],
    "insert": ["insert"],
    "home": ["home"],
    "end": ["end"],
    "pageup": ["pageup"],
    "pagedown": ["pagedown"],
    "capslock": ["capslock"],
    ",": ["comma"],
    ".": ["period"],
    "/": ["slash"],
    "`": ["backquote", "intlbackslash"],
    "-": ["minus"],
    "=": ["equal"],
    ";": ["semicolon"],
    "'": ["quote"],
    "[": ["bracketleft"],
    "]": ["bracketright"],
    "\\": ["backslash"],
};

for (let i = 1; i <= 12; i++) {
    codeMap[`f${i}`] = [`F${i}`];
}

export const macSymbols: Record<string, string> = {
    alt: "⌥",
    backspace: "⌫",
    capslock: "⇪",
    control: "⌃",
    delete: "⌫",
    down: "↓",
    end: "↘",
    enter: "↩",
    escape: "⎋",
    home: "↖",
    left: "←",
    meta: "⌘",
    pagedown: "⇟",
    pageup: "⇞",
    right: "→",
    shift: "⇧",
    tab: "⇥",
    up: "↑",
};

export const platform: "linux" | "mac" | "win" | "other" = navigator.platform.startsWith("Linux")
    ? "linux"
    : navigator.platform.startsWith("Mac")
    ? "mac"
    : navigator.platform.startsWith("Win")
    ? "win"
    : "other";
