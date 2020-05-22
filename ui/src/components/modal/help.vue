<template>
    <v-modal name="help">
        <v-form slot-scope="modal">
            <h1>{{ $t("labels.help") }}</h1>

            <v-tab-bar v-model="tab">
                <v-tab value="keyboard">{{ $t("labels.keyboard.shortcuts._") }}</v-tab>
                <v-tab value="documentation">{{ $t("labels.documentation") }}</v-tab>
            </v-tab-bar>

            <template v-if="tab === 'keyboard'">
                <v-spinner :size="24" v-if="loading"></v-spinner>

                <table v-else>
                    <tr :key="index" v-for="(shortcut, index) in shortcuts">
                        <th>{{ shortcut.label }}</th>
                        <td>{{ shortcut.help }}</td>
                    </tr>
                </table>

                <v-message type="info" v-if="!shortcutsEnabled">
                    {{ $t("messages.keyboard.shortcutsDisabled") }}
                </v-message>
            </template>

            <template v-else-if="tab === 'documentation'">
                <p>{{ $t("messages.documentation") }}</p>

                <v-button href="https://docs.facette.io" icon="external-link-alt" target="_blank">
                    {{ $t("labels.visit.documentation") }}
                </v-button>
            </template>

            <template slot="bottom">
                <v-button primary @click="modal.close()">{{ $t("labels.ok") }}</v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator";

import {ShortcutRecord} from "@/types/components";

function sortableShortcut(input: string): string {
    return input.split(" ").reverse().join("");
}

@Component
export default class ModalHelpComponent extends Vue {
    public loading = false;

    public shortcuts: Array<ShortcutRecord> = [];

    public tab = "keyboard";

    public onModalShow(): void {
        this.loading = true;
        this.update();
    }

    public get shortcutsEnabled(): boolean {
        return this.$components.state.shortcuts;
    }

    private update(): void {
        const shortcuts = this.$components.shortcuts();
        shortcuts.sort((a, b) => sortableShortcut(a.label).localeCompare(sortableShortcut(b.label)));

        Object.assign(this, {
            loading: false,
            shortcuts,
        });
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    ::v-deep .v-modal-content {
        max-width: 35vw;
    }

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
    }

    .v-spinner {
        margin: 0 0.5rem;
    }
}
</style>
