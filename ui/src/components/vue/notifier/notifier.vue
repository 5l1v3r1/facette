<template>
    <div class="v-notifier" :class="{active: !hideTimeout && current && current.text, [type]: type}" @click="dequeue">
        <v-icon :icon="current.icon" v-if="current && current.icon"></v-icon>

        <div class="v-notifier-text" v-if="current">{{ current.text }}</div>

        <div class="v-notifier-timer"></div>
    </div>
</template>

<script lang="ts">
import {Cancelable} from "lodash";
import debounce from "lodash/debounce";
import {Component, Vue} from "vue-property-decorator";

import {Notification} from "@/types/components";

@Component
export default class NotifierComponent extends Vue {
    public current: Notification | null = null;

    public hideTimeout: number | null = null;

    public type: string | null = null;

    private debounce!: (() => unknown) & Cancelable;

    private queue: Array<Notification> = [];

    private timer!: HTMLElement;

    public mounted(): void {
        this.$components.$on("notify", this.enqueue);

        this.timer = this.$el.querySelector(".v-notifier-timer") as HTMLElement;
        this.debounce = debounce(() => this.dequeue(), 6000);
    }

    public beforeDestroy(): void {
        this.$components.$off("notify", this.enqueue);
    }

    public dequeue(): void {
        const current: Notification | null = this.queue.pop() ?? null;
        if (!current) {
            this.hideTimeout = setTimeout(() => {
                this.current = null;
            }, 350);
            return;
        }

        this.current = current;
        this.hideTimeout = null;

        if (this.current.type) {
            this.type = this.current.type;
        }

        // Re-arm CSS animation
        requestAnimationFrame(() => {
            this.timer.style.display = "none";
            requestAnimationFrame(() => {
                this.timer.style.display = "unset";
            });
        });

        if (this.debounce) {
            this.debounce();
        }
    }

    private enqueue(text: string, type = "info", icon?: string): void {
        if (!icon) {
            switch (type) {
                case "error": {
                    icon = "times-circle";
                    break;
                }

                case "success": {
                    icon = "check-circle";
                    break;
                }

                case "warning": {
                    icon = "exclamation-circle";
                    break;
                }
            }
        }

        // Cancel hide timeout if any
        if (this.hideTimeout) {
            clearTimeout(this.hideTimeout);
        }

        // Re-enqueue current message
        if (this.current) {
            this.queue.push(this.current);
        }

        // Remove preexisting message duplicate if found
        const mesg: Notification = {text, type, icon};

        const index = this.queue
            .map(notification => Object.values(notification).join("\x1e"))
            .indexOf(Object.values(mesg).join("\x1e"));

        if (index !== -1) {
            this.queue.splice(index, 1);
        }
        this.queue.push(mesg);

        this.dequeue();
    }
}
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
    box-shadow: 0 0.2rem 0.5rem var(--notifier-shadow);
    color: var(--notifier-color);
    cursor: pointer;
    display: flex;
    height: 0;
    line-height: 3rem;
    padding: 0 1.5rem;
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

    &.info {
        background-color: var(--notifier-info-background);
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
        margin-right: 0.75rem;
    }
}
</style>
