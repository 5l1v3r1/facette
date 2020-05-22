<template>
    <v-toolbar clip="app">
        <v-button
            icon="bars"
            :shortcut="['t', $t('labels.toggleSidebar')]"
            :tooltip="$t('labels.toggleSidebar')"
            @click="toggleSidebar"
        ></v-button>

        <v-button class="logo" :to="{name: 'root'}">
            <img alt="Facette" src="../../assets/facette.svg" />
        </v-button>

        <v-spacer></v-spacer>

        <v-button class="icon" dropdown-anchor="bottom-right" icon="ellipsis-v" :tooltip="$t('labels.openMenu')">
            <template slot="dropdown">
                <v-button
                    icon="cog"
                    :shortcut="['alt+shift+a', $t('labels.goto.adminPanel')]"
                    :to="{name: 'admin-root'}"
                >
                    {{ $t("labels.adminPanel") }}
                </v-button>
                <v-divider></v-divider>
                <template v-if="fullscreenSupport">
                    <v-button
                        :icon="fullscreen ? 'compress' : 'expand'"
                        :shortcut="['meta+enter', $t('labels.fullscreen.toggle')]"
                        @click="toggleFullscreen"
                    >
                        {{ $t(`labels.fullscreen.${fullscreen ? "leave" : "enter"}`) }}
                    </v-button>
                    <v-divider></v-divider>
                </template>
                <v-button icon="info-circle" :shortcut="['?', $t('labels.displayHelp')]" @click="displayHelp">
                    {{ $t("labels.help") }}
                </v-button>
                <v-button
                    icon="sliders-h"
                    :shortcut="['alt+,', $t('labels.goto.settings')]"
                    :to="{name: 'settings', hash: '#display'}"
                >
                    {{ $t("labels.goto.settingsAlt") }}
                </v-button>
            </template>
        </v-button>
    </v-toolbar>
</template>

<script lang="ts">
import {Component, Mixins} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

@Component
export default class ToolbarMainComponent extends Mixins<CustomMixins>(CustomMixins) {
    public fullscreen = false;

    public mounted(): void {
        document.addEventListener("fullscreenchange", this.onFullscreen);
    }

    public beforeDestroy(): void {
        document.removeEventListener("fullscreenchange", this.onFullscreen);
    }

    public displayHelp(): void {
        this.$components.modal("help");
    }

    public get fullscreenSupport(): boolean {
        return document.fullscreenEnabled;
    }

    public onFullscreen(): void {
        this.fullscreen = document.fullscreen;
    }

    public toggleFullscreen(): void {
        if (document.fullscreen) {
            document.exitFullscreen();
        } else {
            document.documentElement.requestFullscreen();
        }
    }
}
</script>

<style lang="scss" scoped>
.v-toolbar {
    .v-button:first-child {
        margin-right: 0.15rem;
    }

    .logo {
        img {
            max-height: 1.5rem;
        }

        ::v-deep .v-button-content {
            padding: 0 0.25rem;
        }
    }
}
</style>
