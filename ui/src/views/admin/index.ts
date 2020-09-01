/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {Component} from "vue";
import {RouteRecordRaw} from "vue-router";

import chartsEdit from "./charts/edit.vue";
import chartsSidebar from "./charts/sidebar.vue";
import dashboardsEdit from "./dashboards/edit.vue";
import dashboardsSidebar from "./dashboards/sidebar.vue";
import database from "./database/database.vue";
import info from "./info/info.vue";
import list from "./common/list.vue";
import metricsList from "./metrics/list.vue";
import providersEdit from "./providers/edit.vue";
import providersSidebar from "./providers/sidebar.vue";
import sidebar from "./common/sidebar.vue";
import toolbar from "../common/toolbar.vue";

const routes: Array<RouteRecordRaw> = [
    {
        path: "/admin",
        name: "admin-root",
        redirect: {
            name: "admin-dashboards-list",
        },
    },
    ...objectRoutes("charts", chartsEdit, chartsSidebar),
    ...objectRoutes("dashboards", dashboardsEdit, dashboardsSidebar),
    {
        path: "/admin/metrics",
        name: "admin-metrics-list",
        components: {
            default: metricsList,
            sidebar: sidebar,
            toolbar: toolbar,
        },
    },
    ...objectRoutes("providers", providersEdit, providersSidebar),
    {
        path: "/admin/database",
        name: "admin-database",
        components: {
            default: database,
            sidebar: sidebar,
            toolbar: toolbar,
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

export default routes;

function objectRoutes(type: string, editDefault: Component, editSidebar: Component): Array<RouteRecordRaw> {
    return [
        {
            path: `/admin/${type}`,
            name: `admin-${type}-list`,
            components: {
                default: list,
                sidebar: sidebar,
                toolbar: toolbar,
            },
            props: {
                default: (): Record<string, any> => ({type}),
            },
        },
        {
            path: `/admin/${type}/:id`,
            name: `admin-${type}-edit`,
            components: {
                default: editDefault,
                sidebar: editSidebar,
                toolbar: toolbar,
            },
            props: {
                default: (): Record<string, any> => ({type}),
            },
        },
    ];
}
