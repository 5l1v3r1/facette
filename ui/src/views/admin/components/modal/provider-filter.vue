<template>
    <v-modal name="provider-filter">
        <v-form slot-scope="modal">
            <v-label>{{ $t("labels.filters.action._") }}</v-label>
            <v-select
                :help="$t('help.filters.action')"
                :options="actions"
                :placeholder="$t('labels.filters.action.select')"
                v-autofocus
                v-model="rule.action"
            ></v-select>

            <v-label>{{ $tc("labels.labels", 1) }}</v-label>
            <v-input
                :help="$t('help.filters.label')"
                :placeholder="$t('labels.placeholders.example', ['__name__'])"
                v-model="rule.label"
            ></v-input>

            <v-label>{{ $t("labels.filters.pattern") }}</v-label>
            <v-input :help="$t('help.filters.pattern')" v-model="rule.pattern"></v-input>

            <template v-if="rule.action === 'relabel'">
                <v-label>{{ $t("labels.filters.targets._") }}</v-label>
                <v-flex class="targets" direction="column">
                    <v-flex :key="index" v-for="(target, index) in targets">
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

                <v-button class="add" icon="plus" :disabled="'' in rule.targets" @click="addTarget">
                    {{ $t("labels.filters.targets.add") }}
                </v-button>
            </template>

            <template v-else-if="rule.action === 'rewrite'">
                <v-label>{{ $t("labels.filters.into") }}</v-label>
                <v-input :help="$t('help.filters.into')" v-model="rule.into"></v-input>
            </template>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.cancel") }}</v-button>
                <v-button :disabled="!rule.action || !rule.label" primary @click="modal.close(rule)">
                    {{ $t("labels.ok") }}
                </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import merge from "lodash/merge";
import {Component, Vue, Watch} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

export interface ModalProviderFilterParams {
    rule: FilterRule;
}

interface Target {
    key: string;
    value: string;
}

const defaultRule: FilterRule = {
    action: "discard",
    label: "",
    pattern: "",
    into: "",
    targets: {},
};

const Actions = ["discard", "relabel", "rewrite", "sieve"];

@Component
export default class ModalProviderFilterComponent extends Vue {
    public rule: FilterRule = cloneDeep(defaultRule);

    public targets: Array<Target> = [];

    public get actions(): Array<SelectOption> {
        return Actions.map(action => ({label: action, value: action}));
    }

    public addTarget(): void {
        this.targets.push({key: "", value: ""});
    }

    public onModalShow(params: ModalProviderFilterParams): void {
        this.rule = merge({}, defaultRule, params.rule);

        if (this.rule.targets) {
            const targets = this.rule.targets;
            this.targets = Object.keys(this.rule.targets).map(key => ({key, value: targets[key]}));
        } else {
            this.targets = [{key: "", value: ""}];
        }
    }

    @Watch("targets", {deep: true})
    public onTargets(to: Array<Target>): void {
        this.rule.targets = to.reduce((targets: Record<string, string>, target: Target) => {
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
