<template>
    <div class="v-message" :class="{inline, [type]: type}">
        <v-icon :icon="iconValue" v-if="iconValue"></v-icon>
        <slot></slot>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class MessageComponent extends Vue {
    @Prop({default: null, type: String})
    public icon!: string | null;

    @Prop({default: false, type: Boolean})
    public inline!: boolean;

    @Prop({type: String})
    public type!: string;

    public get iconValue(): string | null {
        if (this.icon) {
            return this.icon;
        }

        switch (this.type) {
            case "error":
                return "exclamation-circle";

            case "info":
                return "info-circle";

            case "success":
                return "check-circle";

            case "warning":
                return "exclamation-triangle";
        }

        return null;
    }
}
</script>

<style lang="scss" scoped>
.v-message {
    align-items: center;
    border-radius: 0.2rem;
    display: flex;
    padding: 1rem;

    &.inline {
        display: inline-flex;
    }

    &.error {
        background-color: var(--message-error-background);
        color: var(--message-error-color);
    }

    &.info {
        background-color: var(--message-info-background);
        color: var(--message-info-color);
    }

    &.success {
        background-color: var(--message-success-background);
        color: var(--message-success-color);
    }

    &.warning {
        background-color: var(--message-warning-background);
        color: var(--message-warning-color);
    }

    .v-icon {
        font-size: 1.1rem;
        margin-right: 0.75rem;
    }
}
</style>
