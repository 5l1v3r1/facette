/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {SelectOption} from "@/types/components";

export function resolveOption(component: Vue, option: SelectOption): SelectOption {
    return Object.assign({}, option, {label: component.$t(option.label)});
}
