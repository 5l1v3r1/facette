/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import get from "lodash/get";
import kebabCase from "lodash/kebabCase";
import set from "lodash/set";
import {DirectiveOptions, VueConstructor} from "vue";
import {Component, Vue} from "vue-property-decorator";

import {ComponentsOptions, ComponentsState, ModalHandler, ShortcutRecord} from "@/types/components";

import AppComponent from "./app/app.vue";
import ButtonComponent from "./button/button.vue";
import CheckboxComponent from "./checkbox/checkbox.vue";
import ColumnsComponent from "./columns/columns.vue";
import ContentComponent from "./content/content.vue";
import DividerComponent from "./divider/divider.vue";
import DropdownComponent from "./dropdown/dropdown.vue";
import FlexComponent from "./flex/flex.vue";
import FormComponent from "./form/form.vue";
import IconComponent from "./icon/icon.vue";
import InputComponent from "./input/input.vue";
import LabelComponent from "./label/label.vue";
import MessageComponent from "./message/message.vue";
import ModalComponent from "./modal/modal.vue";
import NotifierComponent from "./notifier/notifier.vue";
import PagingComponent from "./paging/paging.vue";
import RadioComponent from "./radio/radio.vue";
import RadioGroupComponent from "./radio/radio-group.vue";
import SelectComponent from "./select/select.vue";
import SidebarComponent from "./sidebar/sidebar.vue";
import SpacerComponent from "./spacer/spacer.vue";
import SpinnerComponent from "./spinner/spinner.vue";
import TabBarComponent from "./tab/tab-bar.vue";
import TabComponent from "./tab/tab.vue";
import TableCellComponent from "./table/cell.vue";
import TableComponent from "./table/table.vue";
import ToolbarComponent from "./toolbar/toolbar.vue";
import TooltipComponent from "./tooltip/tooltip.vue";

import ShortcutDirective, {Shortcut, shortcutLabel, shortcuts} from "./shortcut";
import TooltipDirective from "./tooltip";

const components: Record<string, VueConstructor> = {
    AppComponent,
    ButtonComponent,
    CheckboxComponent,
    ColumnsComponent,
    ContentComponent,
    DividerComponent,
    DropdownComponent,
    FlexComponent,
    FormComponent,
    IconComponent,
    InputComponent,
    LabelComponent,
    MessageComponent,
    ModalComponent,
    NotifierComponent,
    PagingComponent,
    RadioComponent,
    RadioGroupComponent,
    SelectComponent,
    SidebarComponent,
    SpacerComponent,
    SpinnerComponent,
    TabBarComponent,
    TabComponent,
    TableComponent,
    TableCellComponent,
    ToolbarComponent,
    TooltipComponent,
};

const directives: Record<string, DirectiveOptions> = {
    ShortcutDirective,
    TooltipDirective,
};

@Component
class Components extends Vue {
    private data: ComponentsState = {
        sidebar: {
            active: false,
            mode: "static",
        },
        title: "",
        toolbars: {},
        tooltip: null,
    };

    public modal(name: string, params?: unknown, handler?: ModalHandler): void {
        this.$emit("modal-open", name, params, handler);
    }

    public notify(text: string, type?: string, icon?: string): void {
        this.$emit("notify", text, type, icon);
    }

    public set(key: string, payload: unknown): void {
        const data = get(this.data, key);
        set(
            this.data,
            key,
            data instanceof Object && payload instanceof Object ? Object.assign({}, data, payload) : payload,
        );
    }

    public shortcuts(): Array<ShortcutRecord> {
        return Object.values(shortcuts).reduce((records: Array<ShortcutRecord>, shortcut: Shortcut) => {
            if (shortcut.help !== null) {
                records.push({
                    label: shortcutLabel(shortcut.value),
                    help: shortcut.help,
                });
            }
            return records;
        }, []);
    }

    public get state(): ComponentsState {
        return this.data;
    }
}

let installed = false;

export default function Install(vue: typeof Vue, options?: ComponentsOptions): void {
    if (installed) {
        return;
    }
    installed = true;

    const comp: Components = new Components();
    if (options?.title) {
        comp.set("title", options.title);
    }

    Object.assign(vue.prototype, {
        $components: comp,
    });

    Object.keys(components).forEach((key: string) =>
        vue.component(`v-${kebabCase(key.replace(/Component$/, ""))}`, components[key]),
    );

    Object.keys(directives).forEach((key: string) =>
        vue.directive(kebabCase(key.replace(/Directive$/, "")), directives[key]),
    );
}
