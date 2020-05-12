/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import isEqual from "lodash/isEqual";
import Vue from "vue";
import {HttpResponse} from "vue-resource/types/vue_resource";
import VueRouter, {Route, RouteConfig} from "vue-router";
import {Dictionary} from "vue-router/types/router";

Vue.use(VueRouter);

import store from "./store";

import admin from "./views/admin";
import settings from "./views/settings";

interface Options {
    connectors: Array<string>;
}

const routes: Array<RouteConfig> = [
    {
        name: "root",
        path: "/",
    },
].concat(admin, settings);

const router = new VueRouter({
    routes,
    mode: "history",
    base: process.env.BASE_URL,
    linkActiveClass: "active",
    linkExactActiveClass: "active",
});

export default router;

// eslint-disable-next-line @typescript-eslint/no-explicit-any
router.beforeEach((to: Route, from: Route, next: any) => {
    if (from.name !== null && !isEqual(to, from)) {
        store.commit("prevRoute", from);
    }

    next();
});

router.onReady(() => {
    (Vue as any) // eslint-disable-line @typescript-eslint/no-explicit-any
        .http({method: "OPTIONS", url: "/api/v1"})
        .then((response: HttpResponse) => response.json())
        .then(
            (response: APIResponse<Options>) => {
                if (response.data?.connectors) {
                    store.commit("connectors", response.data?.connectors);
                }
            },
            (response: HttpResponse) => {
                console.error(`cannot fetch options: ${response.statusText ?? "unknown error"}`);
            },
        );
});

export function defineParams(route: Route, params: object): void {
    Object.assign(route.params, params);
}

export function updateRouteQuery(route: Route, q: Dictionary<string>, push = false): void {
    let url: string = route.path;

    if (Object.keys(q).length > 0) {
        url += `?${Object.keys(q)
            .map((k: string) => `${k}=${q[k]}`)
            .join("&")}`;
    }

    if (route.hash) {
        url += route.hash;
    }

    history[push ? "pushState" : "replaceState"](null, "", url);
}
