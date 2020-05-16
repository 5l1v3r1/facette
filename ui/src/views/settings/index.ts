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
        path: "/settings",
        name: "settings-root",
        redirect: {
            name: "settings-display",
        },
    },
    {
        path: "/settings/display",
        name: "settings-display",
        components: {
            default: display,
            sidebar: sidebar,
            toolbar: toolbar,
        },
    },
];
