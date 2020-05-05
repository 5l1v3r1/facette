/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {registerComponents} from "@/src/helpers/vue";

import ColorComponent from "./color/color.vue";
import FormChartYaxisComponent from "./form/chart-yaxis.vue";
import FormTemplateVariablesComponent from "./form/template-variables.vue";
import LabelsComponent from "./labels/labels.vue";
import MessageErrorComponent from "./message/error.vue";
import ModalChartPreviewComponent from "./modal/chart-preview.vue";
import ModalChartSeriesComponent from "./modal/chart-series.vue";
import ModalProviderFilterComponent from "./modal/provider-filter.vue";
import ModalTemplateVariableComponent from "./modal/template-variable.vue";

registerComponents({
    ColorComponent,
    FormChartYaxisComponent,
    FormTemplateVariablesComponent,
    LabelsComponent,
    MessageErrorComponent,
    ModalChartPreviewComponent,
    ModalChartSeriesComponent,
    ModalProviderFilterComponent,
    ModalTemplateVariableComponent,
});
