<template>
    <v-modal name="prompt">
        <v-form slot-scope="modal">
            <v-label>{{ modal.params.message }}</v-label>
            <v-input
                :help="modal.params.input && modal.params.input.help ? modal.params.input.help : null"
                :type="modal.params.input && modal.params.input.type ? modal.params.input.type : 'text'"
                v-autofocus.select
                v-model="modal.params.input.value"
            ></v-input>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.cancel") }}</v-button>
                <v-button
                    :danger="modal.params.button.danger"
                    :primary="modal.params.button.primary"
                    @click="modal.close(modal.params.input.value)"
                >
                    {{ modal.params.button && modal.params.button.label ? modal.params.button.label : $t("labels.ok") }}
                </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator";

@Component
export default class ModalPromptComponent extends Vue {}
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
