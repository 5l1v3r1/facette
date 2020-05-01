/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import display from "./display.vue";
import sidebar from "./sidebar.vue";
import toolbar from "./toolbar.vue";

export default [
    {
        name: "settings-root",
        path: "/settings",
        redirect: {
            name: "settings-display",
        },
    },
    {
        components: {
            default: display,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        name: "settings-display",
        path: "/settings/display",
    },
];
