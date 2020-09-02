/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {TooltipState} from "types/ui";
import {DirectiveBinding} from "vue";

type ElementWithTooltip = HTMLElement & {_binding: DirectiveBinding};

export interface TooltipEvent {
    state: TooltipState;
}

const anchors = ["bottom", "left", "right", "top"];

let tooltipTimeout: number | null = null;

function handle(this: ElementWithTooltip, ev: MouseEvent): void {
    if (tooltipTimeout !== null) {
        clearTimeout(tooltipTimeout);
        tooltipTimeout = null;
    }

    if (
        ev.type === "mouseleave" ||
        (ev.type === "mouseup" && this.closest(".v-button") !== null) ||
        this.closest(".v-button.active") !== null
    ) {
        document.dispatchEvent(new CustomEvent<TooltipEvent>("tooltip-hide"));
        return;
    }

    tooltipTimeout = setTimeout(
        () =>
            document.dispatchEvent(
                new CustomEvent<TooltipEvent>("tooltip-show", {
                    detail: {
                        state: {
                            anchor:
                                (Object.keys(this._binding.modifiers).filter(mod =>
                                    anchors.includes(mod),
                                )?.[0] as TooltipState["anchor"]) ?? "bottom",
                            domRect: this.getBoundingClientRect(),
                            nowrap: this._binding.modifiers.nowrap,
                            content: this._binding.value,
                            shortcut: (ev.target as HTMLElement).dataset.vShortcut,
                        },
                    },
                }),
            ),
        500,
    );
}

export default {
    beforeMount(el: ElementWithTooltip, binding: DirectiveBinding): void {
        el._binding = binding;

        el.addEventListener("mouseenter", handle.bind(el));
        el.addEventListener("mouseleave", handle.bind(el));
        el.addEventListener("mouseup", handle.bind(el));
    },

    beforeUnmount(el: ElementWithTooltip): void {
        el.removeEventListener("mouseenter", handle.bind(el));
        el.removeEventListener("mouseleave", handle.bind(el));
        el.removeEventListener("mouseup", handle.bind(el));
    },
};
