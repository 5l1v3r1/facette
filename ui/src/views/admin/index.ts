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
import dashboardsEdit from "./dashboards/edit.vue";
import dashboardsSidebar from "./dashboards/sidebar.vue";
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
        path: "/admin",
        name: "admin-root",
        redirect: {
            name: "admin-dashboards-list",
        },
    },
    {
        path: "/admin/dashboards/:id",
        name: "admin-dashboards-edit",
        components: {
            default: dashboardsEdit,
            sidebar: dashboardsSidebar,
            toolbar: toolbar,
        },
        props: {
            default: (route: Route): void => defineParams(route, {type: "dashboards"}),
        },
    },
    {
        path: "/admin/dashboards",
        name: "admin-dashboards-list",
        components: {
            default: list,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        props: {
            default: (route: Route): void => defineParams(route, {type: "dashboards"}),
        },
    },
    {
        path: "/admin/charts/:id",
        name: "admin-charts-edit",
        components: {
            default: chartsEdit,
            sidebar: chartsSidebar,
            toolbar: toolbar,
        },
        props: {
            default: (route: Route): void => defineParams(route, {type: "charts"}),
        },
    },
    {
        path: "/admin/charts",
        name: "admin-charts-list",
        components: {
            default: list,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        props: {
            default: (route: Route): void => defineParams(route, {type: "charts"}),
        },
    },
    {
        path: "/admin/metrics",
        name: "admin-metrics-list",
        components: {
            default: metricsList,
            sidebar: sidebar,
            toolbar: toolbar,
        },
    },
    {
        path: "/admin/providers/:id",
        name: "admin-providers-edit",
        components: {
            default: providersEdit,
            sidebar: providersSidebar,
            toolbar: toolbar,
        },
        props: {
            default: (route: Route): void => defineParams(route, {type: "providers"}),
        },
    },
    {
        path: "/admin/providers",
        name: "admin-providers-list",
        components: {
            default: list,
            sidebar: sidebar,
            toolbar: toolbar,
        },
        props: {
            default: (route: Route): void => defineParams(route, {type: "providers"}),
        },
    },
    {
        path: "/admin/info",
        name: "admin-info",
        components: {
            default: info,
            sidebar: sidebar,
            toolbar: toolbar,
        },
    },
];
