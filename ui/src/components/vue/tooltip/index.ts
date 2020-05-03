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
            const modifiers: Array<string> = Object.keys(wt._binding.modifiers);

            timeout = setTimeout(() => {
                wt._components.set("tooltip", {
                    anchor: modifiers.length > 0 ? modifiers[modifiers.length - 1] : "bottom",
                    rect: this.getBoundingClientRect(),
                    value: wt._binding.value,
                } as TooltipState);
            }, 750);

            break;
        }

        case "mouseleave": {
            wt._components.set("tooltip", null);
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
