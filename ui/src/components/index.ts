/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import kebabCase from "lodash/kebabCase";
import {App, Component} from "vue";

import ChartComponent from "./chart/chart.vue";
import DatetimeComponent from "./datetime/datetime.vue";
import GridComponent from "./grid/grid.vue";
import GridItemComponent from "./grid/grid-item.vue";
import LabelsComponent from "./labels/labels.vue";
import MessageErrorComponent from "./message/error.vue";
import ModalConfirmComponent from "./modal/confirm.vue";
import ModalHelpComponent from "./modal/help.vue";
import ModalPromptComponent from "./modal/prompt.vue";
import ModalTimeRangeComponent from "./modal/time-range.vue";
import TextComponent from "./text/text.vue";

const components: Record<string, Component> = {
    ChartComponent,
    DatetimeComponent,
    GridComponent,
    GridItemComponent,
    LabelsComponent,
    MessageErrorComponent,
    ModalConfirmComponent,
    ModalHelpComponent,
    ModalPromptComponent,
    ModalTimeRangeComponent,
    TextComponent,
};

export default {
    install: (app: App): void => {
        Object.keys(components).forEach(key =>
            app.component(`v-${kebabCase(key.replace(/Component$/, ""))}`, components[key]),
        );
    },
};
