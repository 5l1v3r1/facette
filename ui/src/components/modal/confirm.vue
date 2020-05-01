<template>
    <v-modal name="confirm">
        <v-form slot-scope="modal">
            <p>{{ modal.params.message }}</p>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.cancel") }}</v-button>
                <v-button
                    autofocus
                    :danger="modal.params.button.danger"
                    :primary="modal.params.button.primary"
                    @click="modal.close(true)"
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
export default class ModalConfirmComponent extends Vue {}
</script>

<style lang="scss" scoped>
.v-modal {
    ::v-deep .v-modal-content {
        max-width: 35vw;
    }
}
</style>
