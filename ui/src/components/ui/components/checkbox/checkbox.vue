<template>
    <div
        class="v-checkbox"
        role="checkbox"
        tabindex="0"
        :aria-checked="Boolean(value) || undefined"
        :aria-disabled="disabled || undefined"
        :aria-labelledby="hasLabel ? `checkbox${uid}` : undefined"
        :aria-readonly="readonly || undefined"
        :class="{[type]: true}"
        @click="toggle"
        @keydown.enter="toggle"
    >
        <v-icon :icon="icon"></v-icon>

        <label class="v-checkbox-label" :id="`checkbox${uid}`" v-if="hasLabel">
            <slot></slot>
        </label>

        <input ref="input" type="checkbox" :checked="value" :readonly="readonly" @change="update" />
    </div>
</template>

<script lang="ts">
import {SetupContext, computed, onMounted, ref} from "vue";

let uid = 0;

export default {
    props: {
        disabled: {
            default: false,
            type: Boolean,
        },
        readonly: {
            default: false,
            type: Boolean,
        },
        type: {
            default: "check",
            type: String,
        },
        value: {
            default: false,
            required: true,
            type: Boolean,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const hasLabel = ref(false);
        const input = ref<HTMLInputElement | null>(null);

        const icon = computed(() => {
            if (props.type === "toggle") {
                return props.value ? "toggle-on" : "toggle-off";
            } else {
                return props.value ? "check" : "";
            }
        });

        const toggle = (): void => {
            if (input.value) {
                input.value.checked = !input.value.checked;
                update();
            }
        };

        const update = (): void => {
            ctx.emit("update:value", input.value?.checked);
        };

        onMounted(() => {
            hasLabel.value = Boolean(ctx.slots.default);
        });

        // Increment global UID counter to avoid label tag identifiers
        // colisions
        uid++;

        return {
            hasLabel,
            icon,
            input,
            toggle,
            uid,
            update,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-checkbox {
    align-items: center;
    cursor: pointer;
    display: inline-flex;
    height: var(--input-height);
    line-height: var(--input-height);
    position: relative;
    vertical-align: middle;

    &[aria-disabled="true"] {
        color: var(--light-gray);
        pointer-events: none;

        &.check[aria-checked="true"] .v-icon {
            background-color: var(--gray);
        }
    }

    &:focus {
        outline: none;
    }

    &:focus-visible::before {
        background-color: rgba(255, 255, 255, 0.35);
        border: 0.15rem solid white;
        border-radius: 0.2rem;
        bottom: -0.25rem;
        content: "";
        left: -0.75rem;
        position: absolute;
        right: -0.75rem;
        top: -0.25rem;
        z-index: 1;
    }

    &.check {
        &:focus .v-icon {
            border-color: var(--color);
        }

        &[aria-checked="true"]:focus .v-icon {
            background-color: var(--checkbox-focus);
            color: white;
        }

        .v-icon {
            background-color: var(--checkbox-background);
            border-radius: 0.2rem;
            border: 1.5px solid var(--checkbox-border);
            color: var(--checkbox-color);
            height: 1em;
            margin: 0.125em;
            pointer-events: none;
            width: 1em;

            ::v-deep(.fa) {
                font-size: 0.75rem;
            }
        }

        &[aria-checked="true"] .v-icon {
            background-color: var(--color);
            border: none;
            color: var(--background);
        }
    }

    &.toggle .v-icon {
        font-size: 1.1rem;
    }

    .v-checkbox-label {
        flex-grow: 1;
        margin-left: 0.5rem;
        pointer-events: none;
    }

    input {
        display: none;
    }
}
</style>
