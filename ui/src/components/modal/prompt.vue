<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-modal name="prompt" @show="onShow">
        <template v-slot="modal">
            <v-form ref="form">
                <v-label>{{ message }}</v-label>
                <v-input
                    :custom-validity="input.customValidity"
                    :delay="input.customValidity ? 350 : 0"
                    :help="input.help"
                    :required="input.required"
                    :type="input.type"
                    v-autofocus.select
                    v-model:value="input.value"
                ></v-input>

                <template v-slot:bottom>
                    <v-button @click="modal.close(false)">
                        {{ i18n.t("labels.cancel") }}
                    </v-button>

                    <v-button
                        :danger="button.danger"
                        :primary="button.primary"
                        @click="canValidate() && modal.close(input.value)"
                    >
                        {{ button.label || i18n.t("labels.ok") }}
                    </v-button>
                </template>
            </v-form>
        </template>
    </v-modal>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import merge from "lodash/merge";
import {ref} from "vue";
import {useI18n} from "vue-i18n";

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

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();

        const button = ref<ModalPromptParams["button"]>(cloneDeep(defaultParams.button));
        const form = ref<HTMLFormElement | null>(null);
        const input = ref<ModalPromptParams["input"]>(cloneDeep(defaultParams.input));
        const message = ref(defaultParams.message);

        const onShow = (params: ModalPromptParams): void => {
            button.value = merge({}, defaultParams.button, params.button);
            input.value = merge({}, defaultParams.input, params.input);
            message.value = params.message;
        };

        const canValidate = (): boolean => {
            return form.value?.$el.checkValidity();
        };

        return {
            button,
            canValidate,
            form,
            i18n,
            input,
            message,
            onShow,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-modal {
    ::v-deep(.v-modal-content) {
        max-width: 35vw;
    }
}
</style>
