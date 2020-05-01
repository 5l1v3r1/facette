<template>
    <v-modal name="help">
        <v-form slot-scope="modal">
            <h1>{{ $t("labels.help") }}</h1>

            <v-tab-bar v-model="tab">
                <v-tab value="keyboard">{{ $t("labels.keyboardShortcuts") }}</v-tab>
                <v-tab value="documentation">{{ $t("labels.documentation") }}</v-tab>
            </v-tab-bar>

            <table v-if="tab === 'keyboard'">
                <tr :key="index" v-for="(shortcut, index) in shortcuts">
                    <th>{{ shortcut.label }}</th>
                    <td>{{ shortcut.help }}</td>
                </tr>
            </table>

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
import {Component, Vue, Watch} from "vue-property-decorator";

import {ShortcutRecord} from "@/types/components";

@Component
export default class ModalHelpComponent extends Vue {
    public shortcuts: Array<ShortcutRecord> = [];

    public tab = "keyboard";

    public mounted(): void {
        this.update();
    }

    @Watch("$route.path")
    public onRoutePath(): void {
        this.$nextTick(() => {
            this.update();
        });
    }

    private update(): void {
        this.shortcuts = this.$components.shortcuts();
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
}
</style>
