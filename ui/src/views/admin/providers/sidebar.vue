<template>
    <v-sidebar :active="sidebar">
        <v-toolbar clip="sidebar">
            <v-button
                icon="arrow-left"
                :shortcut="['alt+up', $t('labels.goto.providers')]"
                :to="{name: 'admin-providers-list'}"
                v-tooltip="{shortcut: 'alt+up'}"
            >
                {{ $t("labels.goto.providers") }}
            </v-button>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <template v-else>
            <v-label>{{ $tc("labels.providers._", 1) }}</v-label>

            <v-button
                exact
                :badge="visited.general && !validity.general ? '!' : null"
                :class="{invalid: visited.general && !validity.general}"
                :to="{name: 'admin-providers-edit', params: {id: params.id}}"
            >
                {{ $t("labels.general") }}
            </v-button>

            <v-button
                exact
                :badge="filters || null"
                :to="{name: 'admin-providers-edit', params: {id: params.id}, hash: '#filters'}"
            >
                {{ $tc("labels.filters._", 2) }}
            </v-button>
        </template>
    </v-sidebar>
</template>

<script lang="ts">
import {Component, Mixins, Watch} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

@Component
export default class Sidebar extends Mixins<CustomMixins>(CustomMixins) {
    public filters: number | null = null;

    public loading = true;

    public validity: Record<string, boolean> = {};

    public visited: Record<string, boolean> = {};

    public mounted(): void {
        this.$parent.$on("provider-updated", this.onProviderUpdated);
        this.$parent.$on("provider-validity", this.onProviderValidity);
    }

    public beforeDestroy(): void {
        this.$parent.$off("provider-updated", this.onProviderUpdated);
        this.$parent.$off("provider-validity", this.onProviderValidity);
    }

    public onProviderUpdated(filters: number | null): void {
        if (this.loading) {
            this.loading = false;
        }
        this.filters = filters;
    }

    public onProviderValidity(validity: Record<string, boolean>): void {
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
