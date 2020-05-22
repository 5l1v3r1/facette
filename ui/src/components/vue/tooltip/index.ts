/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import Vue, {VNode} from "vue";
import {DirectiveBinding} from "vue/types/options";

import {Components, TooltipState} from "@/types/components";

interface ElementWithTooltip {
    _binding: DirectiveBinding;
    _components: Components;
}

const anchors = ["bottom", "left", "right", "top"];

let timeout: number | null = null;

function toggle(this: HTMLElement, e: MouseEvent | null = null): void {
    if (timeout !== null) {
        clearTimeout(timeout);
        timeout = null;
    }

    const wt = this as HTMLElement & ElementWithTooltip;

    if (e === null) {
        wt._components.set("tooltip", null);
        return;
    }

    switch (e.type) {
        case "mouseenter": {
            // Skip tooltip if element is active (useful for dropdown menus)
            if ((e.target as HTMLElement).classList.contains("active")) {
                wt._components.set("tooltip", null);
                return;
            }

            let message = "";
            let shortcut: string | null = null;

            switch (typeof wt._binding.value) {
                case "string":
                    message = wt._binding.value;
                    break;

                case "object":
                    message = wt._binding.value.message ?? null;
                    shortcut = wt._binding.value.shortcut?.[0] ?? null;
                    break;

                default:
                    return;
            }

            const modifiers: Array<string> = Object.keys(wt._binding.modifiers);

            timeout = setTimeout(() => {
                wt._components.set("tooltip", {
                    anchor: modifiers.filter(mod => anchors.includes(mod))?.[0] ?? "bottom",
                    rect: this.getBoundingClientRect(),
                    message,
                    shortcut,
                    wrap: !modifiers.includes("nowrap"),
                } as TooltipState);
            }, 750);

            break;
        }

        case "mouseleave": {
            timeout = setTimeout(() => {
                wt._components.set("tooltip", null);
            }, 250);
        }
    }
}

function bind(el: HTMLElement, binding: DirectiveBinding, vnode: VNode): void {
    if (!binding.value) {
        return;
    }

    const wt = el as HTMLElement & ElementWithTooltip;

    if (binding.oldValue === undefined) {
        wt._binding = binding;
        wt._components = (vnode.context as Vue).$components;

        el.addEventListener("mouseenter", toggle.bind(el));
        el.addEventListener("mouseleave", toggle.bind(el));
    }
}

function unbind(el: HTMLElement): void {
    const wt = el as HTMLElement & ElementWithTooltip;

    if (wt._components) {
        toggle.bind(el)(null);
    }

    el.removeEventListener("mouseenter", toggle);
    el.removeEventListener("mouseleave", toggle);
}

export default {
    bind,
    unbind,
    update: bind,
};
