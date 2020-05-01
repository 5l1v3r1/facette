<template>
    <v-sidebar :active="sidebar">
        <v-toolbar clip="sidebar">
            <v-button icon="arrow-left" :to="{name: 'root'}">{{ $t("labels.goto.home") }}</v-button>
        </v-toolbar>

        <v-label>{{ $t("labels.settings.personal") }}</v-label>
        <v-button icon="desktop" :to="{name: 'settings-display'}">{{ $t("labels.settings.display._") }}</v-button>
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
