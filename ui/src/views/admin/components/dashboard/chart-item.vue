<template>
    <div>
        <v-spinner :size="24" v-if="loading"></v-spinner>

        <v-message type="error" v-else-if="erred">{{ $t("messages.error.load") }}</v-message>

        <template v-else>
            <v-label>{{ $tc("labels.charts._", 1) }}</v-label>
            <v-flex>
                <v-select
                    :options="charts"
                    :placeholder="$t('labels.charts.select')"
                    v-autofocus
                    v-model="options.id"
                ></v-select>

                <v-button
                    icon="pencil-alt"
                    :to="{name: 'admin-charts-edit', params: {id: options.id}}"
                    v-show="options.id"
                >
                    {{ $t("labels.edit") }}
                </v-button>
            </v-flex>

            <v-checkbox toggle v-model="options.legend">{{ $t("labels.charts.showLegend") }}</v-checkbox>
        </template>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

@Component
export default class ChartItemComponent extends Vue {
    @Prop({required: true, type: Object})
    public options!: Record<string, unknown>;

    public charts: Array<SelectOption> = [];

    public erred = false;

    public loading = true;

    public mounted(): void {
        this.getCharts();
    }

    private getCharts(): void {
        this.$http
            .get("/api/v1/charts", {params: {kind: "plain"}})
            .then(response => response.json())
            .then(
                (response: APIResponse<Array<Chart>>) => {
                    if (response.data) {
                        this.charts = response.data.map(chart => ({label: chart.name, value: chart.id}));
                    }
                    this.loading = false;
                },
                () => {
                    this.erred = true;
                    this.loading = false;
                },
            );
    }
}
</script>

<style lang="scss" scoped>
.v-checkbox {
    margin-top: 1rem;
}
</style>
