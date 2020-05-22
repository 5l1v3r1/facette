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

                <v-button exact :to="{name: 'admin-dashboards-edit', params: {id: params.id}}">
                    {{ $t("labels.general") }}
                </v-button>

                <v-button
                    exact
                    :badge="items || null"
                    :to="{name: 'admin-dashboards-edit', params: {id: params.id}, hash: '#layout'}"
                    v-if="!link"
                >
                    {{ $t("labels.layout") }}
                </v-button>

                <v-button
                    exact
                    :badge="variables || null"
                    :to="{name: 'admin-dashboards-edit', params: {id: params.id}, hash: '#variables'}"
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

    public items: number | null = null;

    public variables: number | null = null;

    public mounted(): void {
        this.$parent.$on("dashboard-updated", this.onDashboardUpdated);
    }

    public beforeDestroy(): void {
        this.$parent.$off("dashboard-updated", this.onDashboardUpdated);
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
