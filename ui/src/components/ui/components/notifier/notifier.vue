<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div
        class="v-notifier"
        :class="{active, [notification.type]: notification.type}"
        @click="dismiss"
        v-if="notification"
    >
        <v-icon :icon="iconValue" v-if="iconValue"></v-icon>

        <div class="v-notifier-text" v-if="notification">{{ notification.text }}</div>

        <div class="v-notifier-timer" ref="timer"></div>
    </div>
</template>

<script lang="ts">
import {computed, onMounted, ref} from "vue";

import {Notification} from "types/ui";

import {useUI} from "../..";

export default {
    setup(): Record<string, unknown> {
        const ui = useUI();

        let hideTimeout: number | null = null;

        const active = ref(true);
        const notification = ref<Notification | null>(null);
        const timer = ref<HTMLElement | null>(null);

        const iconValue = computed((): string | null => {
            if (notification.value?.icon) {
                return notification.value.icon;
            }

            switch (notification.value?.type) {
                case "error":
                    return "exclamation-circle";

                case "success":
                    return "check-circle";

                case "warning":
                    return "exclamation-triangle";
            }

            return null;
        });

        const dismiss = (): void => {
            if (hideTimeout !== null) {
                clearTimeout(hideTimeout);
                hideTimeout = null;
            }
            active.value = false;
        };

        const notify = (value: Notification): void => {
            active.value = true;
            notification.value = value;

            if (hideTimeout !== null) {
                clearTimeout(hideTimeout);
                hideTimeout = null;
            }

            // Re-arm CSS animation
            if (timer.value !== null) {
                const style = timer.value.style;
                requestAnimationFrame(() => {
                    style.display = "none";
                    requestAnimationFrame(() => (style.display = "unset"));
                });
            }

            hideTimeout = setTimeout(() => (active.value = false), 6000);
        };

        onMounted(() => (ui.state.notifier = notify));

        return {
            active,
            dismiss,
            iconValue,
            notification,
            timer,
        };
    },
};
</script>

<style lang="scss" scoped>
@keyframes shrink {
    0% {
        width: 100%;
    }

    100% {
        width: 0;
    }
}

.v-notifier {
    border-radius: 0.2rem;
    box-shadow: 0 0.1rem 0.5rem var(--notifier-shadow);
    color: var(--notifier-color);
    cursor: pointer;
    display: flex;
    height: 0;
    line-height: 3rem;
    padding: 0 1rem;
    overflow: hidden;
    position: fixed;
    right: 0.65rem;
    top: 0.65rem;
    transition: height 0.25s ease-in-out;
    will-change: height;
    z-index: 600;

    &.active {
        height: 3rem;
    }

    &.error {
        background-color: var(--notifier-error-background);
    }

    &.success {
        background-color: var(--notifier-success-background);
    }

    &.warning {
        background-color: var(--notifier-warning-background);
    }

    .v-notifier-text {
        flex-grow: 1;
    }

    .v-notifier-timer {
        animation: shrink 6s linear;
        background-color: var(--notifier-highlight-background);
        bottom: 0;
        display: block;
        height: 0.2rem;
        left: 0;
        position: absolute;
        transition: width 0;
        width: 100%;
        will-change: width;
    }

    .v-icon {
        font-size: 1rem;
        height: 3rem;
        margin-right: 1rem;
    }
}
</style>
