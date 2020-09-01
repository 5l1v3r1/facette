<template>
    <v-modal name="chart-series" :title="i18n.t(`labels.series.${edit ? 'edit' : 'add'}`)" @show="onShow">
        <template v-slot="modal">
            <v-form>
                <v-label>{{ i18n.t("labels.series._", 1) }}</v-label>
                <v-input required spellcheck="false" type="textarea" v-autofocus v-model:value="series.expr"></v-input>

                <v-flex class="columns">
                    <v-flex direction="column">
                        <v-label>{{ i18n.t("labels.color") }}</v-label>
                        <v-color class="half" v-model:value="series.options.color"></v-color>
                    </v-flex>

                    <v-flex direction="column">
                        <v-label>{{ i18n.t("labels.charts.axes._", 1) }}</v-label>
                        <v-select
                            class="half"
                            required
                            :options="axes"
                            :placeholder="i18n.t('labels.charts.axes.select')"
                            v-model:value="series.options.axis"
                        ></v-select>
                    </v-flex>
                </v-flex>

                <template v-slot:bottom>
                    <v-button @click="modal.close(false)">
                        {{ i18n.t("labels.cancel") }}
                    </v-button>

                    <v-button :disabled="!series.expr || !series.options.axis" primary @click="modal.close(series)">
                        {{ i18n.t("labels.series.set") }}
                    </v-button>
                </template>
            </v-form>
        </template>
    </v-modal>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import merge from "lodash/merge";
import {ref} from "vue";
import {useI18n} from "vue-i18n";

import {SelectOption} from "types/ui";

export interface ModalChartSeriesParams {
    edit: boolean;
    series: ChartSeries;
}

const defaultSeries: ChartSeries = {
    expr: "",
    options: {
        color: "",
        axis: "left",
    },
};

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();

        const edit = ref(false);
        const series = ref<ChartSeries>(cloneDeep(defaultSeries));

        const axes = ref<Array<SelectOption>>([
            {label: i18n.t("labels.charts.axes.left"), value: "left"},
            {label: i18n.t("labels.charts.axes.right"), value: "right"},
        ]);

        const onShow = (params: ModalChartSeriesParams): void => {
            edit.value = params.edit;
            series.value = merge({}, defaultSeries, params.series);
        };

        return {
            axes,
            edit,
            i18n,
            onShow,
            series,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../../../mixins";
@import "../../../../components/ui/components/form/mixins";

.v-modal {
    ::v-deep(.v-modal-content) {
        @include content;

        min-width: 35vw;

        .columns .column {
            @include form;
        }
    }
}
</style>
