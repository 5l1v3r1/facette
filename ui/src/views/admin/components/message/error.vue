<template>
    <v-message type="placeholder">
        <v-icon icon="exclamation-circle"></v-icon>

        {{ $t("messages.error.load") }}

        <v-button @click="retry">{{ $t("labels.retry") }}</v-button>
    </v-message>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator";

@Component
export default class MessageErrorComponent extends Vue {
    public retry(): void {
        this.$root.$emit("set-error", false);
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
        margin-bottom: 1rem;
    }

    .v-button {
        font-size: 0.85rem;
        margin-top: 2rem;
    }
}
</style>
