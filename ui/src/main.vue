<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <router-view name="toolbar"></router-view>
    <router-view name="sidebar"></router-view>
    <router-view></router-view>

    <v-notifier></v-notifier>
    <v-modal-confirm></v-modal-confirm>
    <v-modal-help></v-modal-help>
    <v-modal-prompt></v-modal-prompt>
    <v-modal-settings></v-modal-settings>
    <v-tooltip></v-tooltip>
</template>

<script lang="ts">
import {onBeforeUnmount, onMounted} from "vue";
import {useStore} from "vuex";

import {useUI} from "@/components/ui";
import {State} from "@/store";

export default {
    setup(): void {
        const store = useStore<State>();
        const ui = useUI();

        const onModifier = (ev: Event | KeyboardEvent): void => {
            switch (ev.type) {
                case "keydown":
                case "keyup": {
                    const ke = ev as KeyboardEvent;
                    let key: string;

                    if (ke.code === "AltLeft" || ke.code === "AltRight") {
                        key = "alt";
                    } else if (ke.code === "ShiftLeft" || ke.code === "ShiftRight") {
                        key = "shift";
                    } else {
                        break;
                    }

                    store.commit(
                        "modifiers",
                        Object.assign({}, store.state.modifiers, {
                            [key]: ev.type === "keydown",
                        }),
                    );

                    break;
                }

                case "visibilitychange": {
                    store.commit("modifiers", {alt: false, shift: false});

                    break;
                }
            }
        };

        onMounted(() => {
            document.addEventListener("keydown", onModifier);
            document.addEventListener("keyup", onModifier);
            document.addEventListener("visibilitychange", onModifier);

            // Trigger pending notification if found
            const notification = store.state.pendingNotification;
            if (notification !== null) {
                ui.notify(notification.text, notification.type, notification.icon);
                store.commit("pendingNotification", null);
            }
        });

        onBeforeUnmount(() => {
            document.removeEventListener("keydown", onModifier);
            document.removeEventListener("keyup", onModifier);
            document.removeEventListener("visibilitychange", onModifier);
        });
    },
};
</script>

<style lang="scss">
$fa-font-path: "~@fortawesome/fontawesome-free/webfonts";
@import "~@fortawesome/fontawesome-free/scss/regular";
@import "~@fortawesome/fontawesome-free/scss/solid";

body[data-v-theme="dark"] {
    --grid-item-background: #272727;
}

body[data-v-theme="light"] {
    --grid-item-background: #fafafa;
}
</style>
