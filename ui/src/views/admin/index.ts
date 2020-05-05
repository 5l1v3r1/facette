/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {Route} from "vue-router";

import {defineParams} from "@/src/router";

import "./components";

import chartsEdit from "./charts/edit.vue";
import chartsSidebar from "./charts/sidebar.vue";
import info from "./info/info.vue";
import list from "./list.vue";
import metricsList from "./metrics/list.vue";
import providersEdit from "./providers/edit.vue";
import providersSidebar from "./providers/sidebar.vue";
import sidebar from "./sidebar.vue";
import toolbar from "./toolbar.vue";

export const namePattern = "^[a-zA-Z0-9](?:[a-zA-Z0-9-_]*[a-zA-Z0-9])?$";

export default [
    {
        name: "admin-root",
        path: "/admin",
        redirect: {
            name: "admin-dashboards-list",
        },
    },
    {
        components: {
            // default: editDashboards,
            // sidebar: sidebarDashboards,
            toolbar: toolbar,
        },
        name: "admin-dashboards-edit",
        path: "/admin/dashboards/:id",
        props: {
            default: (route: Route) => defineParams(route, {type: "dashboards"}),
        },
    },
    {
        components: {
            default: list,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        name: "admin-dashboards-list",
        path: "/admin/dashboards",
        props: {
            default: (route: Route) => defineParams(route, {type: "dashboards"}),
        },
    },
    {
        components: {
            default: chartsEdit,
            sidebar: chartsSidebar,
            toolbar: toolbar,
        },
        name: "admin-charts-edit",
        path: "/admin/charts/:id",
        props: {
            default: (route: Route) => defineParams(route, {type: "charts"}),
        },
    },
    {
        components: {
            default: list,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        name: "admin-charts-list",
        path: "/admin/charts",
        props: {
            default: (route: Route) => defineParams(route, {type: "charts"}),
        },
    },
    {
        components: {
            default: metricsList,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        name: "admin-metrics-list",
        path: "/admin/metrics",
    },
    {
        components: {
            default: providersEdit,
            sidebar: providersSidebar,
            toolbar: toolbar,
        },
        name: "admin-providers-edit",
        path: "/admin/providers/:id",
        props: {
            default: (route: Route) => defineParams(route, {type: "providers"}),
        },
    },
    {
        components: {
            default: list,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        name: "admin-providers-list",
        path: "/admin/providers",
        props: {
            default: (route: Route) => defineParams(route, {type: "providers"}),
        },
    },
    {
        components: {
            default: info,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        name: "admin-info",
        path: "/admin/info",
    },
];
