/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

declare module "@vue/runtime-core" {
    interface ComponentCustomProperties {
        $ui: UI;
    }
}

import {ComponentPublicInstance, UnwrapRef} from "vue";

export declare type AppTheme = Record<string, string>;

export declare class UI {
    public options: UIOptions;

    public state: UnwrapRef<UIState>;

    constructor(options?: UIOptions);

    public modal<T = unknown>(name: string, params?: unknown): Promise<T>;

    public notify(text: string, type: "error" | "success" | "warning", icon?: string): void;

    public shortcuts(): Array<Shortcut>;
}

export declare interface UIOptions {
    theme?: string;
    title?: string;
}

export declare interface UIState {
    modals: Record<string, (params?: unknown) => Promise<unknown>>;
    notifier: ((notification: Notification) => void) | null;
    shortcuts: {
        enabled: boolean;
        entries: Record<string, Shortcut>;
    };
    sidebar: SidebarState;
    title: string;
    toolbars: ToolbarState;
    tooltip: TooltipState | null;
}

export declare type ToolbarState = Record<string, "horizontal" | "vertical" | null>;

export declare interface ButtonComponent extends ComponentPublicInstance {
    focus: () => void;
}

export declare type CheckboxComponent = ComponentPublicInstance;

export declare type ContentComponent = ComponentPublicInstance;

export declare type DividerComponent = ComponentPublicInstance;

export declare type DropdownComponent = ComponentPublicInstance;

export declare type FlexComponent = ComponentPublicInstance;

export declare interface FormComponent extends ComponentPublicInstance {
    checkValidity: () => boolean;
    reset: () => void;
}

export declare type IconComponent = ComponentPublicInstance;

export declare interface InputComponent extends ComponentPublicInstance {
    focus: (select: boolean) => void;
}

export declare type LabelComponent = ComponentPublicInstance;

export declare type MarkdownComponent = ComponentPublicInstance;

export declare type MessageComponent = ComponentPublicInstance;

export declare type ModalComponent = ComponentPublicInstance;

export declare type NotifierComponent = ComponentPublicInstance;

export declare interface Notification {
    text: string;
    type: "error" | "success" | "warning";
    icon?: string;
}

export declare type PagingComponent = ComponentPublicInstance;

export declare interface SelectComponent extends ComponentPublicInstance {
    focus: (select: boolean) => void;
}

export declare interface SelectOption {
    label: string;
    value: unknown;
    icon?: string;
}

export declare type SidebarComponent = ComponentPublicInstance;

export declare interface SidebarState {
    active: boolean;
    mode: string;
}

export declare type SpacerComponent = ComponentPublicInstance;

export declare type SpinnerComponent = ComponentPublicInstance;

export declare type TablistComponent = ComponentPublicInstance;

export declare interface Tab {
    label: string;
    value: unknown;
    icon?: string;
}

export declare type TableCellComponent = ComponentPublicInstance;

export declare type TableComponent = ComponentPublicInstance;

export declare type ToolbarComponent = ComponentPublicInstance;

export declare interface TooltipState {
    anchor: "bottom" | "left" | "right" | "top";
    message?: string;
    shortcut?: string;
    wrap: boolean;
}

export interface Shortcut {
    keys: string;
    help?: string;
}
