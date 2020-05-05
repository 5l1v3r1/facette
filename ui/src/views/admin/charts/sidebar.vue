<template>
    <v-sidebar :active="sidebar">
        <v-toolbar clip="sidebar">
            <v-button
                icon="arrow-left"
                :shortcut="['alt+up', $t('labels.charts.goBack')]"
                :to="{name: 'charts-show'}"
                :tooltip="''"
                v-if="prevRoute === 'charts-show'"
            >
                {{ $t("labels.charts.goBack") }}
            </v-button>

            <v-button
                icon="arrow-left"
                :shortcut="['alt+up', $tc('labels.goto.charts', 2)]"
                :to="{name: 'admin-charts-list'}"
                :tooltip="''"
                v-else
            >
                {{ $tc("labels.goto.charts", 2) }}
            </v-button>
        </v-toolbar>

        <template v-if="!erred">
            <v-spinner v-if="loading"></v-spinner>

            <template v-else>
                <v-label>{{ $tc("labels.charts._", 1) }}</v-label>

                <v-button exact :to="{name: 'admin-charts-edit', params: {id: params.id}}">
                    {{ $t("labels.general") }}
                </v-button>

                <template v-if="!link">
                    <v-button
                        exact
                        :badge="series || null"
                        :to="{name: 'admin-charts-edit', params: {id: params.id}, hash: '#series'}"
                    >
                        {{ $tc("labels.series._", 2) }}
                    </v-button>

                    <v-button exact :to="{name: 'admin-charts-edit', params: {id: params.id}, hash: '#axes'}">
                        {{ $t("labels.charts.axes._") }}
                    </v-button>
                </template>

                <v-button
                    exact
                    :badge="variables || null"
                    :to="{name: 'admin-charts-edit', params: {id: params.id}, hash: '#variables'}"
                    v-else-if="variables !== null"
                >
                    {{ $t("labels.variables._") }}
                </v-button>
            </template>
        </template>
    </v-sidebar>
</template>

<script lang="ts">
import {Component, Mixins} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

@Component
export default class Sidebar extends Mixins<CustomMixins>(CustomMixins) {
    public link = false;

    public loading = true;

    public series: number | null = null;

    public variables: number | null = null;

    public mounted(): void {
        this.$parent.$on("chart-updated", this.onChartUpdated);
    }

    public beforeDestroy(): void {
        this.$parent.$off("chart-updated", this.onChartUpdated);
    }

    public get edit(): boolean {
        return this.params.id !== "link" && this.params.id !== "new";
    }

    public onChartUpdated(link: boolean, series: number | null, variables: number | null): void {
        if (this.loading) {
            this.loading = false;
        }
        Object.assign(this, {link, series, variables});
    }

    public get prevRoute(): string | null {
        return this.$store.getters.prevRoute;
    }
}
</script>

<style lang="scss" scoped>
.v-sidebar {
    .v-toolbar .v-button {
        flex-grow: 1;

        ::v-deep .v-button-content {
            justify-content: flex-start;
        }
    }
}
</style>
