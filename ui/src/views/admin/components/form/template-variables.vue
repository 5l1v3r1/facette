<template>
    <v-form class="half">
        <v-message type="info" v-if="variables.length === 0">
            {{ $t("messages.variables.none") }}
        </v-message>

        <v-table v-model="variables" v-else>
            <template slot="header">
                <v-table-cell>{{ $t("labels.name") }}</v-table-cell>
                <v-table-cell grow>{{ $t("labels.properties") }}</v-table-cell>
                <v-table-cell class="more"></v-table-cell>
            </template>

            <template slot-scope="variable">
                <v-table-cell>
                    {{ variable.value.name }}
                </v-table-cell>

                <v-table-cell grow>
                    <template v-if="variableDefined(variable.value)">
                        <v-labels
                            :labels="{dynamic: true, label: variable.value.label}"
                            v-if="variable.value.dynamic"
                        ></v-labels>

                        <v-labels :labels="{dynamic: false, value: variable.value.value}" v-else></v-labels>
                    </template>

                    <span class="notdefined" v-else>
                        {{ $t("messages.notDefined") }}
                    </span>
                </v-table-cell>

                <v-table-cell class="more">
                    <v-button
                        class="reveal"
                        icon="pencil-alt"
                        @click="editVariable(variable.index)"
                        v-tooltip="$t('labels.variables.edit')"
                    ></v-button>
                    <v-button
                        class="reveal"
                        icon="eraser"
                        :disabled="!variableDefined(variable.value)"
                        @click="clearVariable(variable.index)"
                        v-tooltip="$t('labels.variables.clear')"
                    ></v-button>
                </v-table-cell>
            </template>
        </v-table>
    </v-form>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import {Component, Prop, Vue} from "vue-property-decorator";

import {ModalTemplateVariableParams} from "@/src/views/admin/components/modal/template-variable.vue";

@Component
export default class TemplateVariablesComponent extends Vue {
    @Prop({required: true, type: Array})
    public defined!: Array<TemplateVariable>;

    @Prop({required: true, type: Array})
    public parsed!: Array<TemplateVariable>;

    public get available(): Array<TemplateVariable> {
        const defined = this.defined.map(variable => variable.name);
        return this.parsed.filter(variable => !defined.includes(variable.name));
    }

    public clearVariable(index: number): void {
        this.defined.splice(index, 1);
    }

    public editVariable(index: number): void {
        this.$components.modal(
            "template-variable",
            {
                available: this.available,
                edit: true,
                variable: cloneDeep(this.variables[index]),
            } as ModalTemplateVariableParams,
            (value: TemplateVariable) => {
                if (value) {
                    const definedIndex = this.defined.findIndex(variable => variable.name === value.name);
                    if (definedIndex !== -1) {
                        this.defined.splice(definedIndex, 1, value);
                    } else {
                        this.defined.push(value);
                    }
                }
            },
        );
    }

    public variableDefined(variable: TemplateVariable): boolean {
        return Boolean((!variable.dynamic && variable.value) || (variable.dynamic && variable.label));
    }

    public get variables(): Array<TemplateVariable> {
        return cloneDeep(this.defined.concat(this.available)).sort((a, b) => b.name.localeCompare(a.name));
    }
}
</script>

<style lang="scss" scoped>
.v-form {
    .notdefined {
        color: var(--light-gray);
        text-transform: lowercase;
    }
}
</style>
