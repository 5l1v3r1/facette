<template>
    <v-toolbar clip="app">
        <v-button
            icon="bars"
            @click="toggleSidebar"
            v-shortcut="{keys: 't', help: i18n.t('labels.toggleSidebar')}"
        ></v-button>

        <v-button class="logo" :to="{name: 'root'}">
            <img alt="Facette" src="@/assets/facette.svg" />
        </v-button>

        <v-spacer></v-spacer>

        <v-button
            class="icon"
            dropdown-anchor="bottom-right"
            icon="ellipsis-v"
            v-tooltip.nowrap="i18n.t('labels.openMenu')"
        >
            <template v-slot:dropdown>
                <v-button
                    icon="cog"
                    :to="{name: 'admin-root'}"
                    v-shortcut="{keys: 'alt+shift+a', help: i18n.t('labels.goto.adminPanel')}"
                >
                    {{ i18n.t("labels.adminPanel") }}
                </v-button>

                <v-divider></v-divider>

                <template v-if="fullscreenSupport">
                    <v-button
                        :icon="fullscreen ? 'compress' : 'expand'"
                        @click="toggleFullscreen"
                        v-shortcut="{keys: 'meta+enter', help: i18n.t('labels.fullscreen.toggle')}"
                    >
                        {{ i18n.t(`labels.fullscreen.${fullscreen ? "leave" : "enter"}`) }}
                    </v-button>

                    <v-divider></v-divider>
                </template>

                <v-button
                    icon="info-circle"
                    @click="displayHelp"
                    v-shortcut="{keys: '?', help: i18n.t('labels.displayHelp')}"
                >
                    {{ i18n.t("labels.help") }}
                </v-button>

                <v-button
                    icon="sliders-h"
                    :to="{name: 'settings', params: {section: 'display'}, hash: '#display'}"
                    v-shortcut="{keys: 'alt+,', help: i18n.t('labels.goto.settings')}"
                >
                    {{ i18n.t("labels.goto.settingsAlt") }}
                </v-button>
            </template>
        </v-button>
    </v-toolbar>
</template>

<script lang="ts">
import {computed, onBeforeUnmount, onMounted, ref} from "vue";
import {useI18n} from "vue-i18n";

import common from "@/common";
import {useUI} from "@/components/ui";

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();
        const ui = useUI();

        const {toggleSidebar} = common;

        const fullscreen = ref(false);

        const fullscreenSupport = computed((): boolean => document.fullscreenEnabled);

        const toggleFullscreen = (): void => {
            if (document.fullscreen) {
                document.exitFullscreen();
            } else {
                document.documentElement.requestFullscreen();
            }
        };

        const displayHelp = (): void => {
            ui.modal("help");
        };

        const onFullscreen = (): void => {
            fullscreen.value = document.fullscreen;
        };

        onMounted(() => document.addEventListener("fullscreenchange", onFullscreen));

        onBeforeUnmount(() => document.removeEventListener("fullscreenchange", onFullscreen));

        return {
            displayHelp,
            fullscreen,
            fullscreenSupport,
            i18n,
            toggleFullscreen,
            toggleSidebar,
        };
    },
};
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

        ::v-deep(.v-button-content) {
            padding: 0 0.25rem;
        }
    }
}
</style>
