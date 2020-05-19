<template>
    <v-modal name="confirm">
        <v-form slot-scope="modal">
            <p>{{ message }}</p>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.cancel") }}</v-button>
                <v-button :danger="button.danger" :primary="button.primary" @click="modal.close(true)" v-autofocus>
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

@Component
export default class ModalConfirmComponent extends Vue {
    public button: ModalConfirmParams["button"] = cloneDeep(defaultParams.button);

    public message: ModalConfirmParams["message"] = "";

    public onModalShow(params: ModalConfirmParams): void {
        Object.assign(this, {
            button: merge({}, defaultParams.button, params.button),
            message: params.message,
        });
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    ::v-deep .v-modal-content {
        max-width: 35vw;
    }
}
</style>
