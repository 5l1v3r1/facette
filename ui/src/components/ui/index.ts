/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import kebabCase from "lodash/kebabCase";
import {App, UnwrapRef, inject, nextTick, reactive} from "vue";

import {Shortcut, UIOptions, UIState} from "types/ui";

import components from "./components";
import directives from "./directives";
import themes, {detectTheme} from "./themes";

const uiKey = "components";

export class UI {
    public state: UnwrapRef<UIState>;

    private options: UIOptions;

    constructor(options?: UIOptions) {
        this.options = options ?? {};

        this.state = reactive<UIState>({
            modals: {},
            notifier: null,
            shortcuts: {
                enabled: true,
                entries: {},
            },
            sidebar: {
                active: false,
                mode: "static",
            },
            title: "",
            toolbars: {},
            tooltip: null,
        });
    }

    public install(app: App): void {
        app.config.globalProperties.$ui = this;
        app.provide(uiKey, this);

        Object.keys(components).forEach(key =>
            app.component(`v-${kebabCase(key.replace(/Component$/, ""))}`, components[key]),
        );

        Object.keys(directives).forEach(key =>
            app.directive(kebabCase(key.replace(/Directive$/, "")), directives[key]),
        );

        // Initialize document title
        this.state.title = this.options?.title ?? "";
        this.title();

        // Define CSS custom properties from theme records. Once done, remove
        // "preload" class from parent container.
        //
        // FIXME: find a way to not use "nextTick"
        nextTick(() => {
            const container = app._container as HTMLElement;

            const props = themes[this.options?.theme ?? detectTheme()];
            if (!props) {
                return;
            }

            Object.keys(props).forEach(key => {
                if (key !== "name") {
                    container.style.setProperty(`--${key}`, props[key]);
                }
            });

            container.classList.replace("preload", "v-app");
        });
    }

    public modal<T = unknown>(name: string, params?: unknown): Promise<T> {
        return this.state.modals[name]?.(params) as Promise<T>;
    }

    public notify(text: string, type: "error" | "success" | "warning", icon?: string): void {
        this.state.notifier?.({text, type, icon});
    }

    public shortcuts(): Array<Shortcut> {
        return Object.values(this.state.shortcuts.entries);
    }

    public title(text?: string): void {
        document.title = text ? `${text} – ${this.state.title}` : this.state.title;
    }
}

export function createUI(options?: UIOptions): UI {
    return new UI(options);
}

export function useUI(): UI {
    return inject(uiKey) as UI;
}

export default {
    UI,
    createUI,
    useUI,
};
