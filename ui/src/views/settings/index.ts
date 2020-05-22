/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import settings from "./settings.vue";
import sidebar from "./sidebar.vue";
import toolbar from "./toolbar.vue";

export default [
    {
        path: "/settings",
        name: "settings",
        components: {
            default: settings,
            sidebar: sidebar,
            toolbar: toolbar,
        },
    },
];
