/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {DirectiveBinding} from "vue";

interface AutofocusEvent {
    select: boolean;
}

export default {
    beforeMount(el: HTMLElement): void {
        el.dataset.vAutofocus = "true";
    },
    mounted(el: HTMLElement, binding: DirectiveBinding): void {
        el.dispatchEvent(
            new CustomEvent<AutofocusEvent>("autofocus", {
                detail: {
                    select: binding.modifiers.select,
                },
            }),
        );
    },
};
