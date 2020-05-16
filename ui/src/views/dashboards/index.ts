/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {Route} from "vue-router";

import {defineParams} from "@/src/router";

import "./components";

import home from "./home.vue";
import show from "./show.vue";
import sidebar from "./sidebar.vue";
import toolbar from "./toolbar.vue";

export default [
    {
        path: "/dashboards",
        name: "dashboards-home",
        components: {
            default: home,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        props: {
            default: (route: Route) => defineParams(route, {type: "dashboards"}),
        },
    },
    {
        path: "/dashboards/:id",
        name: "dashboards-show",
        components: {
            default: show,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        props: {
            default: (route: Route) => defineParams(route, {type: "dashboards"}),
        },
    },
];
