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

            <v-button exact :to="{name: 'admin-providers-edit', params: {id: params.id}}">
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
import {Component, Mixins} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

@Component
export default class Sidebar extends Mixins<CustomMixins>(CustomMixins) {
    public filters: number | null = null;

    public loading = true;

    public mounted(): void {
        this.$parent.$on("provider-updated", this.onProviderUpdated);
    }

    public beforeDestroy(): void {
        this.$parent.$off("provider-updated", this.onProviderUpdated);
    }

    public onProviderUpdated(filters: number | null): void {
        if (this.loading) {
            this.loading = false;
        }
        this.filters = filters;
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
