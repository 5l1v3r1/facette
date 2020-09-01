/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {Component} from "vue";

import "./style.scss";

import ButtonComponent from "./button/button.vue";
import CheckboxComponent from "./checkbox/checkbox.vue";
import ColorComponent from "./color/color.vue";
import ContentComponent from "./content/content.vue";
import DividerComponent from "./divider/divider.vue";
import DropdownComponent from "./dropdown/dropdown.vue";
import FlexComponent from "./flex/flex.vue";
import FormComponent from "./form/form.vue";
import IconComponent from "./icon/icon.vue";
import InputComponent from "./input/input.vue";
import LabelComponent from "./label/label.vue";
import MarkdownComponent from "./markdown/markdown.vue";
import MessageComponent from "./message/message.vue";
import ModalComponent from "./modal/modal.vue";
import NotifierComponent from "./notifier/notifier.vue";
import PagingComponent from "./paging/paging.vue";
import SelectComponent from "./select/select.vue";
import SidebarComponent from "./sidebar/sidebar.vue";
import SpacerComponent from "./spacer/spacer.vue";
import SpinnerComponent from "./spinner/spinner.vue";
import TablistComponent from "./tablist/tablist.vue";
import TableCellComponent from "./table/table-cell.vue";
import TableComponent from "./table/table.vue";
import ToolbarComponent from "./toolbar/toolbar.vue";

const components: Record<string, Component> = {
    ButtonComponent,
    CheckboxComponent,
    ColorComponent,
    ContentComponent,
    DividerComponent,
    DropdownComponent,
    FlexComponent,
    FormComponent,
    IconComponent,
    InputComponent,
    LabelComponent,
    MarkdownComponent,
    MessageComponent,
    ModalComponent,
    NotifierComponent,
    PagingComponent,
    SelectComponent,
    SidebarComponent,
    SpacerComponent,
    SpinnerComponent,
    TablistComponent,
    TableCellComponent,
    TableComponent,
    ToolbarComponent,
};

export default components;
