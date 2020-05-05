<template>
    <v-modal name="provider-filter" ref="modal">
        <v-form slot-scope="modal">
            <v-label>{{ $t("labels.filters.action._") }}</v-label>
            <v-select
                :help="$t('help.filters.action')"
                :options="actions"
                :placeholder="$t('labels.filters.action.select')"
                v-autofocus
                v-model="modal.params.action"
            ></v-select>

            <v-label>{{ $tc("labels.labels", 1) }}</v-label>
            <v-input
                :help="$t('help.filters.label')"
                :placeholder="$t('labels.placeholders.example', ['__name__'])"
                v-model="modal.params.label"
            ></v-input>

            <v-label>{{ $t("labels.filters.pattern") }}</v-label>
            <v-input :help="$t('help.filters.pattern')" v-model="modal.params.pattern"></v-input>

            <template v-if="modal.params.action === 'relabel'">
                <v-label>{{ $t("labels.filters.targets._") }}</v-label>
                <v-flex class="targets" direction="column" :key="index" v-for="(target, index) in targets">
                    <v-flex>
                        <v-input
                            :delay="200"
                            :placeholder="$tc('labels.labels', 1)"
                            :ref="index === targets.length - 1 ? 'label' : ''"
                            v-model="target.key"
                        ></v-input>

                        <v-input :delay="200" v-model="target.value"></v-input>

                        <v-button icon="times" @click="removeTarget(index)"></v-button>
                    </v-flex>
                </v-flex>

                <v-button class="add" icon="plus" :disabled="'' in modal.params.targets" @click="addTarget">
                    {{ $t("labels.filters.targets.add") }}
                </v-button>
            </template>

            <template v-else-if="modal.params.action === 'rewrite'">
                <v-label>{{ $t("labels.filters.into") }}</v-label>
                <v-input :help="$t('help.filters.into')" v-model="modal.params.into"></v-input>
            </template>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.cancel") }}</v-button>
                <v-button
                    :disabled="!modal.params.action || !modal.params.label"
                    primary
                    @click="modal.close(modal.params)"
                >
                    {{ $t("labels.ok") }}
                </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import {Component, Vue, Watch} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

import ModalComponent from "@/src/components/vue/modal/modal.vue";

interface FilterParams {
    targets: Record<string, string>;
}

interface Target {
    key: string;
    value: string;
}

const Actions = ["discard", "relabel", "rewrite", "sieve"];

@Component
export default class ModalProviderFilterComponent extends Vue {
    private targets: Array<Target> = [];

    public get actions(): Array<SelectOption> {
        return Actions.map(action => ({label: action, value: action}));
    }

    public addTarget(): void {
        this.targets.push({key: "", value: ""});
    }

    public onModalShow(params: FilterParams): void {
        if (params && params.targets) {
            this.targets = Object.keys(params.targets).map(key => ({key, value: params.targets[key]}));
        } else {
            this.targets = [{key: "", value: ""}];
        }
    }

    @Watch("targets", {deep: true})
    public onTargets(to: Array<Target>): void {
        const params = (this.$refs.modal as ModalComponent).params as FilterParams;
        if (!params) {
            return;
        }

        params.targets = to.reduce((targets: Record<string, string>, target: Target) => {
            targets[target.key] = target.value;
            return targets;
        }, {});

        // Auto-focus on latest key field on modify
        if (document.activeElement && document.activeElement.nodeName !== "INPUT") {
            this.$nextTick(() => {
                const refs = this.$refs.label as Array<HTMLElement>;
                if (refs && refs[0]) {
                    refs[0].focus();
                }
            });
        }
    }

    public removeTarget(index: number): void {
        this.targets.splice(index, 1);
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    ::v-deep .v-modal-content {
        .v-button.add {
            margin-top: 0.5rem;
            width: 100%;
        }

        .targets {
            ::v-deep .v-flex.row {
                margin-top: 0.35rem;
            }

            .v-button {
                margin-left: 0.5rem;
                min-width: 2rem;
            }

            .v-input:first-child {
                margin-right: 0.5rem;
                min-width: auto;
                width: 35%;
            }
        }
    }
}
</style>
