/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import Vue from "vue";
import {DirectiveBinding} from "vue/types/options";

interface ElementWithFocus {
    __vue__?: Vue & {
        focus: (select?: boolean) => void;
    };
}

function inserted(el: HTMLElement & ElementWithFocus, binding: DirectiveBinding): void {
    const vue = el.__vue__;

    if (vue?.focus) {
        vue.$nextTick(() => {
            vue.focus(binding.modifiers.select ?? false);
        });
    }
}

export default {
    inserted,
};
