/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

declare module "*.vue" {
    import {defineComponent} from "vue";
    const component: ReturnType<typeof defineComponent>;
    export default component;
}

// FIXME: remove upon next release including typings
declare module "colorsys";
