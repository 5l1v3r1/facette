<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-modal name="confirm" @show="onShow">
        <template v-slot="modal">
            <v-form>
                <v-markdown :content="message"></v-markdown>

                <template v-slot:bottom>
                    <v-button @click="modal.close(false)">
                        {{ i18n.t("labels.cancel") }}
                    </v-button>

                    <v-button :danger="button.danger" :primary="button.primary" @click="modal.close(true)" v-autofocus>
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

export interface ModalConfirmParams {
    button?: {
        label?: string | null;
        danger?: boolean;
        primary?: boolean;
    };
    message: string;
}

const defaultParams: ModalConfirmParams = {
    button: {
        label: null,
        danger: false,
        primary: false,
    },
    message: "",
};

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();

        const button = ref<ModalConfirmParams["button"]>(cloneDeep(defaultParams.button));
        const message = ref(defaultParams.message);

        const onShow = (params: ModalConfirmParams): void => {
            button.value = merge({}, defaultParams.button, params.button);
            message.value = params.message;
        };

        return {
            button,
            i18n,
            message,
            onShow,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-modal {
    ::v-deep(.v-modal-content) {
        max-width: 25vw;

        .v-markdown {
            :first-child {
                margin-top: 0;
            }

            :last-child {
                margin-bottom: 0;
            }
        }
    }
}
</style>
