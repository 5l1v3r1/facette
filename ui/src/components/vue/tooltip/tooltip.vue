<template>
    <div class="v-tooltip" :class="{[anchor]: anchor}" @mouseenter="onMouse" @mouseleave="onMouse" v-if="state">
        <v-markdown>{{ state.message }}</v-markdown>
        <span class="v-tooltip-shortcut" v-if="state.shortcut">({{ shortcutLabel(state.shortcut) }})</span>
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

    private locked: TooltipState | null = null;

    public shortcutLabel: (keys: Array<string>) => string = shortcutLabel;

    public get anchor(): string | null {
        return this.value !== null ? this.value.anchor : "bottom";
    }

    @Watch("value")
    public onChange(to: TooltipState | null): void {
        if (this.locked) {
            return;
        } else if (to === null) {
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

    public get state(): TooltipState | null {
        return this.locked || this.value;
    }

    public onMouse(e: MouseEvent): void {
        this.locked = e.type === "mouseenter" ? this.value : null;
    }
}
</script>

<style lang="scss" scoped>
@import "../content/mixins";

.v-tooltip {
    background-color: var(--tooltip-background);
    border-radius: 0.2rem;
    box-shadow: 0 0.2rem 0.5rem var(--tooltip-shadow);
    cursor: default;
    left: 0;
    max-width: 25vw;
    overflow-wrap: break-word;
    padding: 0.5rem 0.75rem;
    position: absolute;
    top: 0;
    visibility: hidden;
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

    :first-child {
        margin-top: 0;
    }

    :last-child {
        margin-bottom: 0;
    }

    ::v-deep {
        @include content;

        ul {
            margin: 0;
        }

        p + ul {
            transform: translateY(-0.25rem);
        }
    }
}
</style>
