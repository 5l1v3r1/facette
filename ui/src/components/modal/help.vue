<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-modal name="help" :title="i18n.t('labels.help')">
        <template v-slot="modal">
            <v-form>
                <v-tablist :tabs="tabs" v-model:value="tab"></v-tablist>

                <template v-if="tab === 'keyboard'">
                    <v-spinner :size="24" v-if="shortcuts.length === 0"></v-spinner>

                    <table v-else>
                        <tr :key="index" v-for="(shortcut, index) in shortcuts">
                            <th>{{ shortcutLabel(shortcut.keys) }}</th>
                            <td>{{ shortcut.help }}</td>
                        </tr>
                    </table>

                    <v-message type="warning" v-if="!shortcutsEnabled">
                        {{ i18n.t("messages.keyboard.shortcutsDisabled") }}
                    </v-message>
                </template>

                <template v-else-if="tab === 'documentation'">
                    <p>{{ i18n.t("messages.documentation") }}</p>

                    <v-button href="https://docs.facette.io" icon="external-link-alt" target="_blank">
                        {{ i18n.t("labels.visit.documentation") }}
                    </v-button>
                </template>

                <template v-slot:bottom>
                    <v-button primary @click="modal.close()">{{ i18n.t("labels.ok") }}</v-button>
                </template>
            </v-form>
        </template>
    </v-modal>
</template>

<script lang="ts">
import {computed, ref} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";

import {Tab} from "types/ui";

import {useUI} from "@/components/ui";
import {shortcutLabel} from "@/components/ui/directives/shortcut";
import {State} from "@/store";

const defaultTab = "keyboard";

function sortableShortcut(input: string): string {
    return input.split(" ").reverse().join("");
}

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();
        const ui = useUI();
        const store = useStore<State>();

        const tab = ref(defaultTab);

        const tabs = ref<Array<Tab>>([
            {label: i18n.t("labels.keyboard.shortcuts._"), value: "keyboard"},
            {label: i18n.t("labels.documentation"), value: "documentation"},
        ]);

        const shortcuts = computed(() => {
            const shortcuts = ui.shortcuts();
            shortcuts.sort((a, b) => sortableShortcut(a.keys).localeCompare(sortableShortcut(b.keys)));
            return shortcuts;
        });

        const shortcutsEnabled = computed(() => store.state.shortcuts);

        return {
            i18n,
            shortcutLabel,
            shortcuts,
            shortcutsEnabled,
            tab,
            tabs,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-modal {
    h1 {
        font-size: 1.35rem;
        font-weight: normal;
        margin-top: 0;
    }

    p {
        margin: 1.5rem 0.5rem;
    }

    table {
        border-spacing: 1rem 0.25rem;

        th {
            color: var(--light-gray);
            font-weight: normal;
        }

        + .v-message {
            margin-top: 2rem;
        }
    }

    .v-spinner {
        margin: 0 0.5rem;
    }

    ::v-deep(.v-modal-content) {
        max-width: 35vw;
    }
}
</style>
