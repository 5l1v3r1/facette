<template>
    <v-sidebar :active="sidebar">
        <v-toolbar clip="sidebar">
            <v-button
                icon="arrow-left"
                :shortcut="['alt+up', $t('labels.charts.goBack')]"
                :to="prevRoute"
                v-tooltip="{shortcut: 'alt+up'}"
                v-if="prevRoute.name === 'charts-show'"
            >
                {{ $t("labels.charts.goBack") }}
            </v-button>

            <v-button
                icon="arrow-left"
                :shortcut="['alt+up', $tc('labels.goto.charts', 2)]"
                :to="{name: 'admin-charts-list'}"
                v-tooltip="{shortcut: 'alt+up'}"
                v-else
            >
                {{ $tc("labels.goto.charts", 2) }}
            </v-button>
        </v-toolbar>

        <template v-if="!erred">
            <v-spinner v-if="loading"></v-spinner>

            <template v-else>
                <v-label>{{ $tc("labels.charts._", 1) }}</v-label>

                <v-button
                    exact
                    :badge="visited.general && !validity.general ? '!' : null"
                    :class="{invalid: visited.general && !validity.general}"
                    :to="{name: 'admin-charts-edit', params: {id: params.id}}"
                >
                    {{ $t("labels.general") }}
                </v-button>

                <template v-if="!link">
                    <v-button
                        exact
                        :badge="visited.series && !validity.series ? '!' : series || null"
                        :class="{invalid: visited.series && !validity.series}"
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
import {Component, Mixins, Watch} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

@Component
export default class Sidebar extends Mixins<CustomMixins>(CustomMixins) {
    public link = false;

    public loading = true;

    public series: number | null = null;

    public validity: Record<string, boolean> = {};

    public variables: number | null = null;

    public visited: Record<string, boolean> = {};

    public mounted(): void {
        this.$parent.$on("chart-updated", this.onChartUpdated);
        this.$parent.$on("chart-validity", this.onChartValidity);
    }

    public beforeDestroy(): void {
        this.$parent.$off("chart-updated", this.onChartUpdated);
        this.$parent.$off("chart-validity", this.onChartValidity);
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

    public onChartValidity(validity: Record<string, boolean>): void {
        this.validity = validity;
    }

    @Watch("$route.hash")
    public onRouteHash(to: string, from: string): void {
        this.visited[from.substr(1) || "general"] = true;
    }
}
</script>

<style lang="scss" scoped>
@import "../mixins";

.v-sidebar {
    @include sidebar;
}
</style>
