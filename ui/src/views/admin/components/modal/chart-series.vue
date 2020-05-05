<template>
    <v-modal name="chart-series">
        <v-form slot-scope="modal">
            <v-label>{{ $tc("labels.series._", 1) }}</v-label>
            <v-textarea required spellcheck="false" v-autofocus v-model="series.expr"></v-textarea>

            <v-label>{{ $t("labels.color") }}</v-label>
            <v-color class="third" v-model="series.options.color"></v-color>

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

interface SeriesParams {
    series: ChartSeries;
}

const defaultSeries: ChartSeries = {
    expr: "",
    options: {
        color: "",
    },
};

@Component
export default class ModalChartSeriesComponent extends Vue {
    public series: ChartSeries = cloneDeep(defaultSeries);

    public onModalShow(params: SeriesParams): void {
        this.series = merge({}, defaultSeries, params.series);
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    .v-label:first-child {
        margin-top: 0;
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
