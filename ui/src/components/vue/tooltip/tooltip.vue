<template>
    <div class="v-tooltip" :class="{[anchor]: anchor}" v-if="value">
        <template v-if="typeof value.value === 'object'">
            <span class="v-tooltip-message" v-if="value.value.message">{{ value.value.message }}</span>
            <span class="v-tooltip-shortcut" v-if="value.value.shortcut"
                >({{ shortcutLabel(value.value.shortcut) }})</span
            >
        </template>

        <span class="v-tooltip-message" v-else>{{ value.value }}</span>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue, Watch} from "vue-property-decorator";

import {TooltipState} from "@/types/components";

import {shortcutLabel} from "../shortcut";

@Component
export default class TooltipComponent extends Vue {
    @Prop({required: true, validator: (v): boolean => typeof v === "object" || v === null})
    public value!: TooltipState | null;

    public shortcutLabel: (keys: Array<string>) => string = shortcutLabel;

    public get anchor(): string | null {
        return this.value !== null ? this.value.anchor : "bottom";
    }

    @Watch("value")
    public onChange(to: TooltipState | null): void {
        if (to === null) {
            Object.assign((this.$el as HTMLElement).style, {visibility: null});
            return;
        }

        this.$nextTick(() => {
            let left = 0;
            let top = 0;

            switch (this.anchor) {
                case "bottom": {
                    left = to.rect.left + to.rect.width / 2 - this.$el.clientWidth / 2;
                    top = to.rect.top + to.rect.height;
                    break;
                }

                case "left": {
                    left = to.rect.left - this.$el.clientWidth;
                    top = to.rect.top + to.rect.height / 2 - this.$el.clientHeight / 2;
                    break;
                }

                case "right": {
                    left = to.rect.left + to.rect.width;
                    top = to.rect.top + to.rect.height / 2 - this.$el.clientHeight / 2;
                    break;
                }

                case "top": {
                    left = to.rect.left + to.rect.width / 2 - this.$el.clientWidth / 2;
                    top = to.rect.top - this.$el.clientHeight;
                    break;
                }
            }

            if (left < 0) {
                left = 0;
            }

            if (left + this.$el.clientWidth > document.body.clientWidth) {
                left = document.body.clientWidth - this.$el.clientWidth;
            }

            Object.assign((this.$el as HTMLElement).style, {
                left: `${left}px`,
                top: `${top}px`,
                visibility: "visible",
            });
        });
    }
}
</script>

<style lang="scss" scoped>
.v-tooltip {
    background-color: var(--tooltip-background);
    border-radius: 0.2rem;
    box-shadow: 0 0.2rem 0.5rem var(--tooltip-shadow);
    left: 0;
    max-width: 30vw;
    overflow-wrap: break-word;
    padding: 0.35rem 0.65rem;
    pointer-events: none;
    position: absolute;
    top: 0;
    visibility: hidden;
    white-space: pre-wrap;
    will-change: left, top, transform;
    z-index: 700;

    &.bottom {
        transform: translateY(0.5rem);
    }

    &.left {
        transform: translateX(-0.5rem);
    }

    &.right {
        transform: translateX(0.5rem);
    }

    &.top {
        transform: translateY(-0.5rem);
    }

    .v-tooltip-shortcut {
        opacity: 0.5;
    }

    .v-tooltip-message + .v-tooltip-shortcut {
        margin-left: 0.5rem;
    }
}
</style>
