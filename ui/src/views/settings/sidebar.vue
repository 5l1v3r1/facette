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
                :to="{name: 'root'}"
                v-shortcut="{keys: 'alt+up', help: i18n.t('labels.goto.home')}"
            >
                {{ i18n.t("labels.goto.home") }}
            </v-button>
        </v-toolbar>

        <v-label>{{ i18n.t("labels.settings.personal") }}</v-label>

        <v-button icon="desktop" :to="{name: 'settings', params: {section: 'display'}, hash: '#display'}">
            {{ i18n.t("labels.display") }}
        </v-button>

        <v-button icon="keyboard" :to="{name: 'settings', params: {section: 'keyboard'}, hash: '#keyboard'}">
            {{ i18n.t("labels.keyboard._") }}
        </v-button>
    </v-sidebar>
</template>

<script lang="ts">
import {useI18n} from "vue-i18n";

import common from "@/common";

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();

        const {sidebar} = common;

        return {
            i18n,
            sidebar,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../mixins";

.v-sidebar {
    @include sidebar;
}
</style>
