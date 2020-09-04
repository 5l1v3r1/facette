/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import isEqual from "lodash/isEqual";
import {NavigationGuardNext, RouteLocationNormalized, RouteRecordRaw, createRouter, createWebHistory} from "vue-router";

import api from "@/lib/api";
import common from "@/common";
import store from "@/store";

import admin from "@/views/admin";
import dashboards from "@/views/dashboards";

const routes = Array<RouteRecordRaw>().concat(
    {
        path: "/",
        name: "root",
        redirect: {
            name: "dashboards-home",
        },
    },
    ...admin,
    ...dashboards,
);

const router = createRouter({
    history: createWebHistory(process.env.BASE_URL),
    routes,
    linkActiveClass: "",
    linkExactActiveClass: "active",
});

export default router;

router.beforeEach((to: RouteLocationNormalized, from: RouteLocationNormalized, next: NavigationGuardNext) => {
    // Each time route is changing, store previous route location and reset
    // error from store.
    if (from.name && to.name !== from.name && !isEqual(to.params, from.params)) {
        store.commit("prevRoute", from);
    }

    store.commit("error", null);

    next();
});

router.isReady().then(() => {
    const {onFetchRejected} = common;

    // Router is ready, thus get options from API and store retrieved data into
    // store.
    api.options().then(response => {
        if (response.data) {
            store.commit("apiOptions", response.data);
        }
    }, onFetchRejected);
});
