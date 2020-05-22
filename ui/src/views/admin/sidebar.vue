<template>
    <v-sidebar :active="sidebar">
        <v-toolbar clip="sidebar">
            <v-button
                icon="arrow-left"
                :shortcut="['alt+up', $t('labels.goto.home')]"
                :to="{name: 'root'}"
                v-tooltip="{shortcut: 'alt+up'}"
            >
                {{ $t("labels.goto.home") }}
            </v-button>
        </v-toolbar>

        <v-label>{{ $t("labels.library") }}</v-label>
        <v-button icon="th-large" :to="{name: 'admin-dashboards-list'}">{{ $tc("labels.dashboards._", 2) }}</v-button>
        <v-button icon="chart-area" :to="{name: 'admin-charts-list'}">{{ $tc("labels.charts._", 2) }}</v-button>

        <v-label>{{ $t("labels.catalog") }}</v-label>
        <v-button icon="dot-circle" :to="{name: 'admin-metrics-list'}">{{ $tc("labels.metrics._", 2) }}</v-button>

        <v-label>{{ $t("labels.system") }}</v-label>
        <v-button icon="plug" :to="{name: 'admin-providers-list'}">{{ $tc("labels.providers._", 2) }}</v-button>
        <v-button icon="info-circle" :to="{name: 'admin-info'}">{{ $t("labels.info._") }}</v-button>
    </v-sidebar>
</template>

<script lang="ts">
import {Component, Mixins} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

@Component
export default class Sidebar extends Mixins<CustomMixins>(CustomMixins) {}
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
