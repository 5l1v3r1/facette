/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {Directive} from "vue";

import AutofocusDirective from "./autofocus";
import ShortcutDirective from "./shortcut";
import TooltipDirective from "./tooltip";

const directives: Record<string, Directive> = {
    AutofocusDirective,
    ShortcutDirective,
    TooltipDirective,
};

export default directives;
