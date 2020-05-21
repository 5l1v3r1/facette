<template>
    <v-modal name="prompt">
        <v-form @validity="onValidity" slot-scope="modal">
            <v-label>{{ message }}</v-label>
            <v-input
                :custom-validity="input.customValidity"
                :delay="input.customValidity ? 200 : 0"
                :help="input.help"
                :required="input.required"
                :type="input.type"
                v-autofocus.select
                v-model="input.value"
            ></v-input>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.cancel") }}</v-button>
                <v-button
                    :danger="button.danger"
                    :disabled="input.customValidity && !validity"
                    :primary="button.primary"
                    @click="modal.close(input.value)"
                >
                    {{ button.label || $t("labels.ok") }}
                </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import merge from "lodash/merge";
import {Component, Vue} from "vue-property-decorator";

export interface ModalPromptParams {
    button?: {
        label?: string | null;
        danger?: boolean;
        primary?: boolean;
    };
    input?: {
        customValidity?: ((value: string) => Promise<string>) | null;
        help?: string | null;
        required?: boolean;
        type?: string;
        value?: unknown;
    };
    message: string;
}

const defaultParams: ModalPromptParams = {
    button: {
        label: null,
        danger: false,
        primary: false,
    },
    input: {
        customValidity: null,
        help: null,
        required: false,
        type: "text",
        value: "",
    },
    message: "",
};

@Component
export default class ModalPromptComponent extends Vue {
    public button: ModalPromptParams["button"] = cloneDeep(defaultParams.button);

    public input: ModalPromptParams["input"] = cloneDeep(defaultParams.input);

    public message: ModalPromptParams["message"] = "";

    public validity = false;

    public onModalShow(params: ModalPromptParams): void {
        Object.assign(this, {
            button: merge({}, defaultParams.button, params.button),
            input: merge({}, defaultParams.input, params.input),
            message: params.message,
        });
    }

    public onValidity(to: boolean): void {
        this.validity = to;
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    ::v-deep .v-modal-content {
        max-width: 35vw;

        .v-input {
            width: 100%;
        }
    }
}
</style>
