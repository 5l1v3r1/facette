/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import Vue from "vue";

declare module "vue/types/vue" {
    interface Vue {
        $components: Components;
    }
}

declare type AppTheme = Record<string, string>;

declare class Components extends Vue {
    public state: ComponentsState;
    public modal(name: string, params?: unknown, handler?: ModalHandler): void;
    public notify(text: string, type?: string, icon?: string): void;
    public set(key: string, payload: unknown): void;
    public shortcuts(): Array<ShortcutRecord>;
}

declare interface ComponentsOptions {
    title?: string;
}

declare interface ComponentsState {
    shortcuts: boolean;
    sidebar: SidebarState;
    title: string;
    toolbars: ToolbarState;
    tooltip: TooltipState | null;
}

declare type ModalHandler = (value: any) => void; // eslint-disable-line @typescript-eslint/no-explicit-any

declare interface Notification {
    text: string;
    type: string;
    icon?: string | null;
}

declare interface SelectOption {
    label: string;
    value: unknown;
    icon?: string | null;
    [key: string]: unknown;
}

declare interface ShortcutRecord {
    label: string;
    help: string;
}

declare interface SidebarState {
    active: boolean;
    mode: string;
}

declare type ToolbarState = Record<string, "horizontal" | "vertical" | null>;

declare interface TooltipState {
    anchor: "bottom" | "left" | "right" | "top";
    rect: ClientRect | DOMRect;
    message: string;
    shortcut: string | null;
}
