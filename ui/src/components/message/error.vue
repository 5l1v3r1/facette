<template>
    <v-message type="placeholder">
        <v-icon icon="dizzy"></v-icon>

        {{ message }}

        <v-button @click="retry">{{ $t("labels.retry") }}</v-button>
    </v-message>
</template>

<script lang="ts">
import {TranslateResult} from "vue-i18n";
import {Component, Mixins} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

@Component
export default class MessageErrorComponent extends Mixins<CustomMixins>(CustomMixins) {
    public get message(): string {
        let result: TranslateResult;

        switch (this.error) {
            case "notFound":
                result = this.$t(`messages.${this.params.type}.notFound`);
                break;

            default:
                result = this.$t(`messages.error.unhandled`);
                break;
        }

        return result as string;
    }

    public retry(): void {
        this.resetError();
        this.$emit("retry");
    }
}
</script>

<style lang="scss" scoped>
.v-message {
    color: var(--gray);
    bottom: 0;
    flex-direction: column;
    font-size: 1rem;
    justify-content: center;
    left: 0;
    position: absolute;
    right: 0;
    top: 0;

    .v-icon {
        color: var(--dark-gray);
        font-size: 4rem;
        margin: 0 0 1rem;
    }

    .v-button {
        font-size: 0.85rem;
        margin-top: 2rem;
    }
}
</style>
