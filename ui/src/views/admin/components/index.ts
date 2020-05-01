/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {registerComponents} from "@/src/helpers/vue";

import LabelsComponent from "./labels/labels.vue";
import MessageErrorComponent from "./message/error.vue";
import ModalChartPreviewComponent from "./modal/chart-preview.vue";
import ModalProviderFilterComponent from "./modal/provider-filter.vue";

registerComponents({
    LabelsComponent,
    MessageErrorComponent,
    ModalChartPreviewComponent,
    ModalProviderFilterComponent,
});
