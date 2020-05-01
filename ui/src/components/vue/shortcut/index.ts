/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import isArray from "lodash/isArray";
import {nanoid} from "nanoid";
import {DirectiveBinding} from "vue/types/options";

import {codeMap, macSymbols, platform} from "./vars";

export interface Shortcut {
    el: HTMLElement;
    value: string;
    keys: Array<string>;
    codes: Array<string>;
    modifiers: Record<string, boolean>;
    help: string | null;
}

interface ElementWithShortcut {
    _id: string | undefined;
}

export const shortcuts: Record<string, Shortcut> = {};

let registered = false;

function keyCode(key: string): string {
    let code: string = key;
    if (code.startsWith("Digit")) {
        code = code.substr(5);
    } else if (code.startsWith("Key")) {
        code = code.substr(3);
    }
    return code.toLowerCase();
}

function splitValue(value: string): Array<string> {
    return value === "+" ? ["+"] : value.split("+");
}

function handle(e: KeyboardEvent): void {
    // TODO: handle input-list?
    if (document.activeElement && document.activeElement.matches("input, select, textarea")) {
        return;
    }

    const shortcut: Shortcut | undefined = Object.values(shortcuts).find((shortcut: Shortcut): boolean => {
        return (
            // Handle codes having modifiers (e.g. "alt+,")
            (shortcut.codes.includes(keyCode(e.code)) &&
                e.altKey === shortcut.modifiers.alt &&
                e.ctrlKey === shortcut.modifiers.control &&
                e.metaKey === shortcut.modifiers.meta &&
                e.shiftKey === shortcut.modifiers.shift) ||
            // Handle simple modified codes (e.g. "?" with implicit Shift modifier)
            (e.key === shortcut.keys[0] &&
                !shortcut.modifiers.alt &&
                !e.altKey &&
                !shortcut.modifiers.control &&
                !e.ctrlKey &&
                !shortcut.modifiers.meta &&
                !e.metaKey &&
                !shortcut.modifiers.shift &&
                e.shiftKey)
        );
    });

    if (shortcut) {
        e.preventDefault();
        shortcut.el.dispatchEvent(new CustomEvent("shortcut", {bubbles: false}));
    }
}

function unbind(el: HTMLElement): void {
    const ws = el as HTMLElement & ElementWithShortcut;

    if (ws.id !== undefined && shortcuts[ws.id]) {
        delete shortcuts[ws.id];

        // Stop listening for shortcuts if none registered
        if (Object.keys(shortcuts).length === 0) {
            document.removeEventListener("keydown", handle);
            registered = false;
        }
    }
}

function bind(el: HTMLElement, binding: DirectiveBinding): void {
    // Unbind any previous bound shortcut
    unbind(el);

    if (!binding.value) {
        return;
    }

    // Start listening for shortcuts (i.e. listener not already registered)
    if (!registered) {
        document.addEventListener("keydown", handle, true);
        registered = true;
    }

    let value: string;
    let help: string | null = null;

    if (typeof binding.value === "string") {
        value = binding.value;
    } else if (isArray(binding.value)) {
        [value, help] = binding.value;
    } else {
        throw Error(`invalid shortcut value: "${binding.value}"`);
    }

    const origKeys: Array<string> = splitValue(value);
    const keys: Array<string> = [...origKeys];
    const modifiers: Record<string, boolean> = {};

    ["alt", "control", "meta", "shift"].forEach((key: string) => {
        const index = keys.indexOf(key);
        modifiers[key] = index !== -1;
        if (modifiers[key]) {
            keys.splice(index, 1);
        }
    });

    if (keys.length === 0) {
        throw Error("missing shortcut base key");
    } else if (keys.length > 1) {
        throw Error("too many shortcut base keys");
    }

    const ws = el as HTMLElement & ElementWithShortcut;
    ws.id = nanoid(8);

    shortcuts[ws.id] = {
        el,
        value,
        keys: origKeys,
        codes: codeMap[keys[0]] ?? [keys[0]],
        modifiers,
        help,
    };
}

export function shortcutLabel(value: string | Array<string>): string {
    if (isArray(value)) {
        value = value[0];
    } else if (typeof value !== "string") {
        throw Error(`invalid shortcut value: "${value}"`);
    }

    return splitValue(value)
        .map((key: string) => {
            switch (platform) {
                case "mac":
                    if (macSymbols[key]) {
                        return macSymbols[key] + " ";
                    }
                    break;

                case "win":
                    if (key === "meta") {
                        return "Win";
                    }
            }

            return key.charAt(0).toUpperCase() + key.substr(1);
        })
        .join(platform === "mac" ? "" : " + ")
        .trimEnd();
}

export default {
    bind,
    unbind,
    update: bind,
};
