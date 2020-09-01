<template>
    <v-sidebar :active="sidebar">
        <v-toolbar clip="sidebar">
            <v-button
                icon="arrow-left"
                :to="prevRoute"
                v-shortcut="{keys: 'alt+up', help: i18n.t('labels.goto.chartBack')}"
                v-if="prevRoute.name === 'charts-show'"
            >
                {{ i18n.t("labels.goto.chartBack") }}
            </v-button>

            <v-button
                icon="arrow-left"
                :to="{name: 'admin-charts-list'}"
                v-shortcut="{keys: 'alt+up', help: i18n.t('labels.goto.charts', 2)}"
                v-else
            >
                {{ i18n.t("labels.goto.charts", 2) }}
            </v-button>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <template v-else>
            <v-label>{{ i18n.t("labels.charts._", 1) }}</v-label>

            <v-button
                :badge="routeData?.invalid?.general ? '!' : null"
                :class="{invalid: routeData?.invalid?.general}"
                :to="{name: 'admin-charts-edit', params: {id: $route.params.id}}"
            >
                {{ i18n.t("labels.general") }}
            </v-button>

            <template v-if="!routeData?.link">
                <v-button
                    :badge="routeData?.series ?? null"
                    :to="{
                        name: 'admin-charts-edit',
                        params: {id: $route.params.id, section: 'series'},
                        hash: '#series',
                    }"
                >
                    {{ i18n.t("labels.series._", 2) }}
                </v-button>

                <v-button
                    :to="{
                        name: 'admin-charts-edit',
                        params: {id: $route.params.id, section: 'axes'},
                        hash: '#axes',
                    }"
                >
                    {{ i18n.t("labels.charts.axes._", 2) }}
                </v-button>

                <v-button
                    :badge="routeData?.markers ?? null"
                    :to="{
                        name: 'admin-charts-edit',
                        params: {id: $route.params.id, section: 'markers'},
                        hash: '#markers',
                    }"
                >
                    {{ i18n.t("labels.markers._", 2) }}
                </v-button>
            </template>

            <v-button
                :badge="routeData?.variables"
                :to="{
                    name: 'admin-charts-edit',
                    params: {id: $route.params.id, section: 'variables'},
                    hash: '#variables',
                }"
                v-else-if="routeData?.variables"
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
