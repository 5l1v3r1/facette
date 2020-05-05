/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import Vue from "vue";
import Vuex from "vuex";
import VuexPersistedState from "vuex-persistedstate";

Vue.use(Vuex);

class State {
    public connectors: Array<string> = [];

    public data: unknown = null;

    public locale = "en";

    public modifierPressed = false;

    public prevRoute: string | null = null;

    public routeGuarded = false;

    public routeGuardFn: ((e: Event) => void) | null = null;

    public sidebar = true;

    public theme = "dark";

    public timezoneUTC = false;
}

const store = new Vuex.Store({
    state: new State(),
    getters: {
        connectors(state: State): Array<string> {
            return state.connectors;
        },

        data(state: State): unknown {
            return state.data;
        },

        locale(state: State): string {
            return state.locale;
        },

        modifierPressed(state: State): boolean {
            return state.modifierPressed;
        },

        prevRoute(state: State): string | null {
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
        connectors(state: State, value: Array<string>): void {
            state.connectors = value;
        },

        data(state: State, value: unknown): void {
            state.data = value;
        },

        locale(state: State, value: string): void {
            state.locale = value;
        },

        modifierPressed(state: State, value: boolean): void {
            state.modifierPressed = value;
        },

        prevRoute(state: State, value: string | null): void {
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
                locale: state.locale,
                prevRoute: state.prevRoute,
                sidebar: state.sidebar,
                theme: state.theme,
                timezoneUTC: state.timezoneUTC,
            }),
        }),
    ],
});

export default store;
