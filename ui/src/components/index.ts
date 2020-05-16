/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {registerComponents} from "@/src/helpers/vue";

import ChartComponent from "./chart/chart.vue";
import DatetimeComponent from "./datetime/datetime.vue";
import GridComponent from "./grid/grid.vue";
import GridItemComponent from "./grid/item.vue";
import MessageErrorComponent from "./message/error.vue";
import ModalConfirmComponent from "./modal/confirm.vue";
import ModalHelpComponent from "./modal/help.vue";
import ModalPromptComponent from "./modal/prompt.vue";
import ToolbarMainComponent from "./toolbar/main.vue";

registerComponents({
    DatetimeComponent,
    GridComponent,
    GridItemComponent,
    ChartComponent,
    MessageErrorComponent,
    ModalConfirmComponent,
    ModalHelpComponent,
    ModalPromptComponent,
    ToolbarMainComponent,
});
