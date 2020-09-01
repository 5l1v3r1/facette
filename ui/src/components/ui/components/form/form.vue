<template>
    <form class="v-form" novalidate ref="el" @submit.prevent>
        <slot></slot>

        <div class="v-form-bottom" v-if="hasBottom">
            <slot name="bottom"></slot>
        </div>
    </form>
</template>

<script lang="ts">
import {SetupContext, onMounted, ref} from "vue";

export default {
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const el = ref<HTMLFormElement | null>(null);
        const hasBottom = ref(false);

        const checkValidity = (): boolean => {
            const valid = el.value?.checkValidity() ?? false;

            if (!valid) {
                el.value
                    ?.querySelectorAll(":invalid")
                    .forEach(input => input.dispatchEvent(new CustomEvent("report-validity")));
            }

            return valid;
        };

        const reset = (): void => {
            el.value
                ?.querySelectorAll(":invalid")
                .forEach(input => input.dispatchEvent(new CustomEvent("reset-validity")));

            el.value?.reset();
        };

        onMounted((): void => {
            hasBottom.value = Boolean(ctx.slots.bottom);
        });

        return {
            checkValidity,
            el,
            hasBottom,
            reset,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "./mixins";

.v-form {
    .v-form-bottom {
        align-items: center;
        display: flex;
        margin-top: 2rem;
    }

    ::v-deep() {
        @include form;

        > :first-child {
            margin-top: 0;
        }
    }
}
</style>
