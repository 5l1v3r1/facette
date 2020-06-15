<template>
    <v-modal name="chart-series">
        <v-form slot-scope="modal">
            <v-label>{{ $tc("labels.series._", 1) }}</v-label>
            <v-textarea required spellcheck="false" v-autofocus v-model="series.expr"></v-textarea>

            <v-flex>
                <v-flex class="third" direction="column">
                    <v-label>{{ $t("labels.color") }}</v-label>
                    <v-color v-model="series.options.color"></v-color>
                </v-flex>

                <v-flex class="third" direction="column">
                    <v-label>{{ $tc("labels.charts.axes._", 1) }}</v-label>
                    <v-select
                        :options="axes"
                        :placeholder="$t('labels.charts.axes.select')"
                        v-model="series.options.axis"
                    ></v-select>
                </v-flex>
            </v-flex>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.close") }}</v-button>
                <v-button primary @click="modal.close(series)">{{ $t("labels.ok") }} </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import cloneDeep from "lodash/merge";
import merge from "lodash/merge";
import {Component, Vue} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

import {resolveOption} from "@/src/helpers/select";

export interface ModalChartSeriesParams {
    series: ChartSeries;
}

const axes: Array<SelectOption> = [
    {label: "labels.charts.axes.left", value: "left"},
    {label: "labels.charts.axes.right", value: "right"},
];

const defaultSeries: ChartSeries = {
    expr: "",
    options: {
        color: "",
        axis: "left",
    },
};

@Component
export default class ModalChartSeriesComponent extends Vue {
    public series: ChartSeries = cloneDeep(defaultSeries);

    public get axes(): Array<SelectOption> {
        return axes.map(option => resolveOption(this, option));
    }

    public onModalShow(params: ModalChartSeriesParams): void {
        this.series = merge({}, defaultSeries, params.series);
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    .v-form > .v-label:first-child {
        margin-top: 0;
    }

    .v-color,
    .v-select {
        min-width: auto;
        width: 100%;
    }

    .v-textarea {
        font-family: "Roboto Mono", monospace;
        width: 30rem;

        ::v-deep textarea {
            min-height: 16rem;
        }
    }
}
</style>
