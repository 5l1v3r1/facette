<template>
    <div class="v-message" :class="{inline, [type]: type}">
        <v-icon :icon="iconValue" v-if="iconValue"></v-icon>
        <slot></slot>
    </div>
</template>

<script lang="ts">
import {computed} from "vue";

export default {
    props: {
        icon: {
            default: null,
            type: String,
        },
        inline: {
            default: false,
            type: Boolean,
        },
        type: {
            default: null,
            type: String,
        },
    },
    setup(props: Record<string, any>): Record<string, unknown> {
        const iconValue = computed((): string | null => {
            if (props.icon !== null) {
                return props.icon;
            }

            switch (props.type) {
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
        });

        return {
            iconValue,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-message {
    align-items: center;
    border-radius: 0.2rem;
    display: flex;
    margin: 1rem 0;
    padding: 0.75rem 1rem;

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
        margin-right: 1rem;
    }
}
</style>
