/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {RouteRecordRaw} from "vue-router";

import home from "./home/home.vue";
import show from "./show/show.vue";
import sidebar from "./common/sidebar.vue";
import toolbar from "../common/toolbar.vue";

const routes: Array<RouteRecordRaw> = [
    {
        path: "/dashboards",
        name: "dashboards-home",
        components: {
            default: home,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        props: {
            default: (): Record<string, any> => ({type: "dashboards"}),
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
            default: (): Record<string, any> => ({type: "dashboards"}),
        },
    },
    {
        path: "/charts/:id",
        name: "charts-show",
        components: {
            default: show,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        props: {
            default: (): Record<string, any> => ({type: "charts"}),
        },
    },
    {
        path: "/basket",
        name: "basket-show",
        components: {
            default: show,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        props: {
            default: (): Record<string, any> => ({type: "basket"}),
        },
    },
];

export default routes;
