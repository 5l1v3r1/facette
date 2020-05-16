/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import Vue from "vue";
import {Route} from "vue-router";
import Vuex from "vuex";
import VuexPersistedState from "vuex-persistedstate";

Vue.use(Vuex);

class State {
    public autoPropagate = true;

    public connectors: Array<string> = [];

    public data: unknown = null;

    public error: APIError = null;

    public locale = "en";

    public modifiers: Modifiers = {
        alt: false,
        shift: false,
    };

    public prevRoute: Route | null = null;

    public routeGuarded = false;

    public routeGuardFn: ((e: Event) => void) | null = null;

    public sidebar = true;

    public theme = "dark";

    public timezoneUTC = false;
}

const store = new Vuex.Store({
    state: new State(),
    getters: {
        autoPropagate(state: State): boolean {
            return state.autoPropagate;
        },

        connectors(state: State): Array<string> {
            return state.connectors;
        },

        data(state: State): unknown {
            return state.data;
        },

        error(state: State): APIError {
            return state.error;
        },

        locale(state: State): string {
            return state.locale;
        },

        modifiers(state: State): Modifiers {
            return state.modifiers;
        },

        prevRoute(state: State): Route | null {
            return state.prevRoute;
        },

        routeGuarded(state: State): boolean {
            return state.routeGuarded;
        },

        routeGuardFn(state: State): ((e: Event) => void) | null {
            return state.routeGuardFn;
        },

        sidebar(state: State): boolean {
            return state.sidebar;
        },

        theme(state: State): string {
            return state.theme;
        },

        timezoneUTC(state: State): boolean {
            return state.timezoneUTC;
        },
    },
    mutations: {
        autoPropagate(state: State, value: boolean): void {
            state.autoPropagate = value;
        },

        connectors(state: State, value: Array<string>): void {
            state.connectors = value;
        },

        data(state: State, value: unknown): void {
            state.data = value;
        },

        error(state: State, value: APIError): void {
            state.error = value;
        },

        locale(state: State, value: string): void {
            state.locale = value;
        },

        modifiers(state: State, value: Modifiers): void {
            state.modifiers = value;
        },

        prevRoute(state: State, value: Route | null): void {
            state.prevRoute = value;
        },

        routeGuarded(state: State, value: boolean): void {
            state.routeGuarded = value;
        },

        routeGuardFn(state: State, value: ((e: Event) => void) | null): void {
            state.routeGuardFn = value;
        },

        sidebar(state: State, value: boolean): void {
            state.sidebar = value;
        },

        theme(state: State, value: string): void {
            state.theme = value;
        },

        timezoneUTC(state: State, value: boolean): void {
            state.timezoneUTC = value;
        },
    },
    plugins: [
        VuexPersistedState({
            key: "facette",
            reducer: (state: State): object => ({
                autoPropagate: state.autoPropagate,
                locale: state.locale,
                prevRoute:
                    state.prevRoute !== null
                        ? {
                              name: state.prevRoute.name,
                              params: state.prevRoute.params,
                          }
                        : null,
                sidebar: state.sidebar,
                theme: state.theme,
                timezoneUTC: state.timezoneUTC,
            }),
        }),
    ],
});

export default store;
