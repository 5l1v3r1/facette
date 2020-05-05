<template>
    <v-modal name="template-variable">
        <v-form slot-scope="modal">
            <v-tab-bar v-model="variable.dynamic">
                <v-tab :value="false">{{ $t("labels.variables.fixed") }}</v-tab>
                <v-tab :value="true">{{ $t("labels.variables.dynamic") }}</v-tab>
            </v-tab-bar>

            <v-label>{{ $t("labels.name") }}</v-label>
            <v-input disabled v-model="variable.name" v-if="edit"></v-input>
            <v-select
                :options="variableOptions"
                :placeholder="$t('labels.variables.select')"
                v-autofocus.select
                v-model="variable.name"
                v-else
            ></v-select>

            <template v-if="variable.dynamic">
                <v-label>{{ $tc("labels.labels", 1) }}</v-label>
                <v-input
                    :placeholder="$t('labels.placeholders.example', ['instance'])"
                    v-autofocus.select="edit"
                    v-model="variable.label"
                ></v-input>

                <v-label>{{ $tc("labels.filters._", 1) }}</v-label>
                <v-input
                    :placeholder="$t('labels.placeholders.example', ['{__provider__=&quot;prometheus&quot;}'])"
                    v-model="variable.filter"
                ></v-input>
            </template>

            <template v-else>
                <v-label>{{ $t("labels.value") }}</v-label>
                <v-input
                    v-autofocus.select="edit"
                    :placeholder="$t('labels.placeholders.example', ['host.example.net'])"
                    v-model="variable.value"
                ></v-input>
            </template>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.close") }}</v-button>
                <v-button
                    class="primary"
                    :disabled="!variable.name || (variable.dynamic && !variable.label)"
                    @click="modal.close(cleanedVariable)"
                >
                    {{ $t("labels.ok") }}
                </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import merge from "lodash/merge";
import {Component, Vue} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

interface TemplateParams {
    available: Array<TemplateVariable>;
    edit: boolean;
    variable: TemplateVariable;
}

const defaultVariable: TemplateVariable = {
    name: "",
    value: "",
    label: "",
    filter: "",
    dynamic: false,
};

@Component
export default class ModalTemplateVariableComponent extends Vue {
    public edit = false;

    public variable: TemplateVariable = cloneDeep(defaultVariable);

    private available: Array<TemplateVariable> = [];

    public get cleanedVariable(): TemplateVariable {
        return Object.assign({}, this.variable, this.variable.dynamic ? {value: ""} : {label: "", filter: ""});
    }

    public onModalShow(params: TemplateParams): void {
        Object.assign(this, {
            available: params.available,
            edit: params.edit,
            variable: merge({}, defaultVariable, params.variable),
        });
    }

    public get variableOptions(): Array<SelectOption> {
        return this.available.map(variable => ({label: variable.name, value: variable.name}));
    }
}
</script>
