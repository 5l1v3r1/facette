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

        <v-label>{{ i18n.t("labels.library") }}</v-label>

        <v-button icon="th-large" :to="{name: 'admin-dashboards-list'}">
            {{ i18n.t("labels.dashboards._", 2) }}
        </v-button>

        <v-button icon="chart-area" :to="{name: 'admin-charts-list'}">
            {{ i18n.t("labels.charts._", 2) }}
        </v-button>

        <v-label>{{ i18n.t("labels.catalog") }}</v-label>

        <v-button icon="dot-circle" :to="{name: 'admin-metrics-list'}">
            {{ i18n.t("labels.metrics._", 2) }}
        </v-button>

        <v-label>{{ i18n.t("labels.system") }}</v-label>

        <v-button icon="plug" :to="{name: 'admin-providers-list'}">
            {{ i18n.t("labels.providers._", 2) }}
        </v-button>

        <v-button icon="database" :to="{name: 'admin-database'}">
            {{ i18n.t("labels.database._") }}
        </v-button>

        <v-button icon="info-circle" :to="{name: 'admin-info'}">
            {{ i18n.t("labels.info._") }}
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
@import "../../mixins";

.v-sidebar {
    @include sidebar;
}
</style>
