<template>
    <v-modal name="chart-preview">
        <v-form slot-scope="modal">
            <v-chart tooltip v-model="chart"></v-chart>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.close") }}</v-button>
                <v-button primary @click="createChart(modal.params.exprs)" v-autofocus>
                    {{ $t("labels.charts.create") }}
                </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator";

interface PreviewParams {
    exprs: Array<string>;
}

@Component
export default class ModalChartPreviewComponent extends Vue {
    public chart: Chart = {
        id: "",
        name: "",
        options: {
            type: "area",
        },
        series: [],
    };

    public createChart(exprs: Array<string>): void {
        this.$store.commit("seriesPrefill", exprs); // TODO: add prefill to store
        this.$router.push({name: "admin-charts-edit", params: {id: "new"}});
    }

    public onModalShow(params: PreviewParams): void {
        Object.assign(this.chart, {
            name: this.$t("labels.preview"),
            series: params.exprs.map(expr => ({expr})),
        });
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    ::v-deep .v-modal-content {
        height: 90vh;
        width: 90vw;

        .v-form {
            height: calc(100% - 3rem);
        }
    }

    .v-chart {
        height: 100%;
        width: 100%;
    }
}
</style>
