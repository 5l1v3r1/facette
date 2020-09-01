/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {RouteLocationNormalized} from "vue-router";
import {MutationPayload, Store, createStore} from "vuex";

import {Notification} from "types/ui";

const persistKey = "facette";

export class State {
    public apiOptions: Options = {
        connectors: [],
        driver: {
            name: "",
            version: "",
        },
    };

    public autoPropagate = true;

    public basket: Array<DashboardItem> = [];

    public error: APIError = null;

    public loading = true;

    public locale = "en";

    public modifiers: Modifiers = {
        alt: false,
        shift: false,
    };

    public pendingNotification: Notification | null = null;

    public routeData: Record<string, unknown> | null = null;

    public routeGuarded = false;

    public shortcuts = true;

    public prevRoute: RouteLocationNormalized | null = null;

    public sidebar = true;

    public theme: string | null = null;

    public timeRange: TimeRange | null = null;

    public timezoneUTC = false;
}

const store = createStore({
    state: new State(),
    mutations: {
        apiOptions(state: State, value: Options): void {
            state.apiOptions = value;
        },
        autoPropagate(state: State, value: boolean): void {
            state.autoPropagate = value;
        },
        basket(state: State, value: Array<DashboardItem>): void {
            state.basket = value;
        },
        error(state: State, value: APIError): void {
            state.error = value;
        },
        locale(state: State, value: string): void {
            state.locale = value;
        },
        loading(state: State, value: boolean): void {
            state.loading = value;
        },
        modifiers(state: State, value: Modifiers): void {
            state.modifiers = value;
        },
        pendingNotification(state: State, value: Notification | null): void {
            state.pendingNotification = value;
        },
        prevRoute(state: State, value: RouteLocationNormalized | null): void {
            state.prevRoute = value;
        },
        routeData(state: State, value: Record<string, unknown> | null): void {
            state.routeData = value;
        },
        routeGuarded(state: State, value: boolean): void {
            state.routeGuarded = value;
        },
        shortcuts(state: State, value: boolean): void {
            state.shortcuts = value;
        },
        sidebar(state: State, value: boolean): void {
            state.sidebar = value;
        },
        theme(state: State, value: string | null): void {
            state.theme = value;
        },
        timeRange(state: State, value: TimeRange | null): void {
            state.timeRange = value;
        },
        timezoneUTC(state: State, value: boolean): void {
            state.timezoneUTC = value;
        },
    },
    plugins: [persist()],
});

function persist() {
    const save = (state: State): void => {
        localStorage.setItem(
            persistKey,
            JSON.stringify({
                autoPropagate: state.autoPropagate,
                basket: state.basket,
                locale: state.locale,
                pendingNotification: state.pendingNotification,
                prevRoute:
                    state.prevRoute !== null
                        ? {
                              name: state.prevRoute.name,
                              params: state.prevRoute.params,
                          }
                        : null,
                shortcuts: state.shortcuts,
                sidebar: state.sidebar,
                theme: state.theme,
                timezoneUTC: state.timezoneUTC,
            }),
        );
    };

    let saveTimeout: number | null = null;

    const saveDebounce = (mutation: MutationPayload, state: State) => {
        if (saveTimeout !== null) {
            clearTimeout(saveTimeout);
        }

        saveTimeout = setTimeout(() => save(state), 250);
    };

    return (store: Store<State>) => {
        // Restore state keys from local storage, then subscribe to mutations
        // to keep it in-sync with live changes.
        try {
            const value = localStorage.getItem(persistKey);
            if (value !== null) {
                store.replaceState(Object.assign({}, store.state, JSON.parse(value)));
            }
        } catch (err) {}

        store.subscribe(saveDebounce);

        // Ensure local storage is in-sync before unloading
        window.addEventListener("beforeunload", () => {
            if (saveTimeout !== null) {
                clearTimeout(saveTimeout);
            }

            save(store.state);
        });
    };
}

export default store;
