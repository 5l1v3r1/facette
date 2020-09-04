<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div class="v-icon">
        <div class="fa-fw" :class="classes"></div>
        <div class="v-icon-badge" :class="{small}" v-if="badge">{{ badge }}</div>
    </div>
</template>

<script lang="ts">
import {computed} from "vue";

export default {
    props: {
        badge: {
            default: null,
            type: [Number, String],
        },
        icon: {
            required: true,
            type: String,
        },
        spin: {
            default: false,
            type: Boolean,
        },
    },
    setup(props: Record<string, any>): Record<string, unknown> {
        const classes = computed((): Array<string> | null => {
            if (!props.icon) {
                return null;
            }

            const idx = props.icon.indexOf("/");

            let classes: Array<string>;
            if (idx !== -1) {
                classes = [props.icon.substr(0, idx), `fa-${props.icon.substr(idx + 1)}`];
            } else {
                classes = ["fa", `fa-${props.icon}`];
            }

            if (props.spin) {
                classes.push("fa-spin");
            }

            return classes;
        });

        const small = computed(() => props.badge === " ");

        return {
            classes,
            small,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "~@fortawesome/fontawesome-free/scss/fontawesome";

.v-icon {
    align-items: center;
    display: flex;
    justify-content: center;
    position: relative;

    .fa {
        pointer-events: none;
    }

    .v-icon-badge {
        background-color: var(--accent);
        box-shadow: 0 0 0 0.2rem var(--background);
        border-radius: 1.1rem;
        color: var(--background);
        left: 100%;
        line-height: 1.1rem;
        min-height: 1.1rem;
        min-width: 1.1rem;
        padding: 0 0.35rem;
        position: absolute;
        top: 100%;
        transform: scale(0.7) translate(-100%, -100%);

        &.small {
            min-height: 0.75rem;
            min-width: 0.75rem;
            padding: 0;
        }
    }

    .fa-spin {
        animation: fa-spin 0.65s infinite linear;
    }
}
</style>
