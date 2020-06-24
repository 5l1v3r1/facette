<template>
    <v-sidebar :active="sidebar">
        <v-toolbar clip="sidebar">
            <v-button
                icon="arrow-left"
                :shortcut="['alt+up', $t('labels.dashboards.goBack')]"
                :to="prevRoute"
                v-tooltip="{shortcut: 'alt+up'}"
                v-if="prevRoute.name === 'dashboards-show'"
            >
                {{ $t("labels.dashboards.goBack") }}
            </v-button>

            <v-button
                icon="arrow-left"
                :shortcut="['alt+up', $tc('labels.goto.dashboards', 2)]"
                :to="{name: 'admin-dashboards-list'}"
                v-tooltip="{shortcut: 'alt+up'}"
                v-else
            >
                {{ $tc("labels.goto.dashboards", 2) }}
            </v-button>
        </v-toolbar>

        <template v-if="!erred">
            <v-spinner v-if="loading"></v-spinner>

            <template v-else>
                <v-label>{{ $tc("labels.dashboards._", 1) }}</v-label>

                <v-button
                    exact
                    :badge="visited.general && !validity.general ? '!' : null"
                    :class="{invalid: visited.general && !validity.general}"
                    :to="{name: 'admin-dashboards-edit', params: {id: params.id}}"
                >
                    {{ $t("labels.general") }}
                </v-button>

                <v-button
                    exact
                    :badge="visited.layout && !validity.layout ? '!' : items || null"
                    :class="{invalid: visited.layout && !validity.layout}"
                    :to="{name: 'admin-dashboards-edit', params: {id: params.id}, hash: '#layout'}"
                >
                    {{ $t("labels.layout") }}
                </v-button>

                <v-button
                    exact
                    :badge="variables || null"
                    :to="{name: 'admin-dashboards-edit', params: {id: params.id}, hash: '#variables'}"
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

    public items: number | null = null;

    public validity: Record<string, boolean> = {};

    public variables: number | null = null;

    public visited: Record<string, boolean> = {};

    public mounted(): void {
        this.$parent.$on("dashboard-updated", this.onDashboardUpdated);
        this.$parent.$on("dashboard-validity", this.onDashboardValidity);
    }

    public beforeDestroy(): void {
        this.$parent.$off("dashboard-updated", this.onDashboardUpdated);
        this.$parent.$off("dashboard-validity", this.onDashboardValidity);
    }

    public get edit(): boolean {
        return this.params.id !== "link" && this.params.id !== "new";
    }

    public onDashboardUpdated(link: boolean, items: number | null, variables: number | null): void {
        if (this.loading) {
            this.loading = false;
        }
        Object.assign(this, {link, items, variables});
    }

    public onDashboardValidity(validity: Record<string, boolean>): void {
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
