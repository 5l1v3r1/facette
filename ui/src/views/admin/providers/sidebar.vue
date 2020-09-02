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
                :to="{name: 'admin-providers-list'}"
                v-shortcut="{keys: 'alt+up', help: i18n.t('labels.goto.providers')}"
            >
                {{ i18n.t("labels.goto.providers") }}
            </v-button>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <template v-else>
            <v-label>{{ i18n.t("labels.providers._", 1) }}</v-label>

            <v-button
                :badge="routeData?.invalid.general ? '!' : null"
                :class="{invalid: routeData?.invalid.general}"
                :to="{name: 'admin-providers-edit', params: {id: $route.params.id}}"
            >
                {{ i18n.t("labels.general") }}
            </v-button>

            <v-button
                :badge="routeData?.filters ?? null"
                :to="{
                    name: 'admin-providers-edit',
                    params: {id: $route.params.id, section: 'filters'},
                    hash: '#filters',
                }"
            >
                {{ i18n.t("labels.filters._", 2) }}
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

        const {loading, sidebar} = common;

        const routeData = computed(() => store.state.routeData);

        return {
            i18n,
            loading,
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
