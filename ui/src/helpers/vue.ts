/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import kebabCase from "lodash/kebabCase";
import Vue, {VueConstructor} from "vue";

export function registerComponents(components: Record<string, VueConstructor>): void {
    Object.keys(components).forEach((key: string) =>
        Vue.component(`v-${kebabCase(key.replace(/Component$/, ""))}`, components[key]),
    );
}
