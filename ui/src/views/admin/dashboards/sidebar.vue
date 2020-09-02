<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-sidebar :active="sidebar">
        <v-toolbar clip="sidebar">
            <v-button
                icon="arrow-left"
                :to="prevRoute"
                v-shortcut="{keys: 'alt+up', help: i18n.t('labels.goto.dashboardBack')}"
                v-if="prevRoute.name === 'dashboards-show'"
            >
                {{ i18n.t("labels.goto.dashboardBack") }}
            </v-button>

            <v-button
                icon="arrow-left"
                :to="{name: 'admin-dashboards-list'}"
                v-shortcut="{keys: 'alt+up', help: i18n.t('labels.goto.dashboards', 2)}"
                v-else
            >
                {{ i18n.t("labels.goto.dashboards", 2) }}
            </v-button>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <template v-else>
            <v-label>{{ i18n.t("labels.dashboards._", 1) }}</v-label>

            <v-button
                :badge="routeData?.invalid?.general ? '!' : null"
                :class="{invalid: routeData?.invalid?.general}"
                :to="{name: 'admin-dashboards-edit', params: {id: $route.params.id}, query: $route.query}"
            >
                {{ i18n.t("labels.general") }}
            </v-button>

            <v-button
                :badge="routeData?.items ?? null"
                :to="{
                    name: 'admin-dashboards-edit',
                    params: {id: $route.params.id, section: 'layout'},
                    query: $route.query,
                    hash: '#layout',
                }"
            >
                {{ i18n.t("labels.layout") }}
            </v-button>

            <v-button
                :badge="routeData?.variables"
                :to="{
                    name: 'admin-dashboards-edit',
                    params: {id: $route.params.id, section: 'variables'},
                    query: $route.query,
                    hash: '#variables',
                }"
                v-if="routeData?.variables"
            >
                {{ i18n.t("labels.variables._", 2) }}
            </v-button>
        </template>
    </v-sidebar>
</template>

<script lang="ts">
import {computed} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";

import common from "@/common";
import {State} from "@/store";

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();
        const store = useStore<State>();

        const {loading, prevRoute, sidebar} = common;

        const routeData = computed(() => store.state.routeData);

        return {
            i18n,
            loading,
            prevRoute,
            routeData,
            sidebar,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../../mixins";

.v-sidebar {
    @include sidebar;
}
</style>
