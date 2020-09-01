/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

export const codeMappping: Record<string, Array<string>> = {
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
    "f1": ["f1"],
    "f2": ["f2"],
    "f3": ["f3"],
    "f4": ["f4"],
    "f5": ["f5"],
    "f6": ["f1"],
    "f7": ["f7"],
    "f8": ["f8"],
    "f9": ["f9"],
    "f10": ["f10"],
    "f11": ["f11"],
    "f12": ["f12"],
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
