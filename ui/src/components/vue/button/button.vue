<template>
    <div
        class="v-button"
        role="button"
        :aria-disabled="disabled"
        :aria-label="tooltip"
        :class="{
            active,
            danger,
            primary,
            icon: !content && (!dropdown || icon.startsWith('caret-')),
            ['in-dropdown']: inDropdown,
        }"
        @shortcut="shortcut && onShortcut()"
        v-shortcut="shortcut"
    >
        <component
            class="v-button-content"
            :exact="exact"
            :href="href"
            :is="tag"
            :tabindex="disabled ? -1 : 0"
            :target="target"
            :to="to"
            @click="onClick"
            @focus="onFocus"
            @focusout="onFocus"
            @keydown.enter="onClick"
            @mouseenter="dropdown && inDropdown && onMouse($event)"
            @mouseleave="dropdown && inDropdown && onMouse($event)"
            v-tooltip="!active && tooltip !== null ? {message: tooltip, shortcut} : null"
        >
            <v-icon :badge="iconBadge" :icon="icon" :spin="iconSpin" v-if="icon !== null"></v-icon>

            <span v-if="content">
                <slot></slot>

                <div class="v-button-badge" v-if="badge">{{ badge }}</div>
            </span>

            <div class="v-button-shortcut" v-if="shortcut && inDropdown">{{ shortcutLabel(shortcut) }}</div>

            <v-icon
                class="v-button-caret"
                :icon="inDropdown ? 'caret-right' : 'caret-down'"
                v-if="dropdown && !icon.startsWith('caret-')"
            ></v-icon>
        </component>

        <v-dropdown
            tabindex="-1"
            :anchor="dropdownAnchor"
            @focusout.native="onFocus"
            @mouseleave.native="inDropdown && onMouse($event)"
            v-if="dropdown"
        >
            <slot name="dropdown"></slot>
        </v-dropdown>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";
import {Route} from "vue-router";

import {shortcutLabel} from "../shortcut";

@Component
export default class ButtonComponent extends Vue {
    @Prop({default: false, type: Boolean})
    public autofocus!: boolean;

    @Prop({default: null, type: [Number, String]})
    public badge!: number | string | null;

    @Prop({default: false, type: Boolean})
    public danger!: boolean;

    @Prop({default: false, type: Boolean})
    public disabled!: boolean;

    @Prop({default: "bottom-left", type: String})
    public dropdownAnchor!: string;

    @Prop({default: false, type: Boolean})
    public exact!: boolean;

    @Prop({default: null, type: String})
    public href!: string;

    @Prop({default: null, type: String})
    public icon!: string;

    @Prop({default: null, type: [Number, String]})
    public iconBadge!: number | string | null;

    @Prop({default: false, type: Boolean})
    public iconSpin!: boolean;

    @Prop({default: false, type: Boolean})
    public primary!: boolean;

    @Prop({default: null, type: [Array, String]})
    public shortcut!: string | Array<string> | null;

    @Prop({default: null, type: String})
    public target!: string;

    @Prop({default: null, type: Object})
    public to!: Route | null;

    @Prop({default: null, type: String})
    public tooltip!: string;

    public active = false;

    public content = false;

    public dropdown = false;

    public inDropdown = false;

    public shortcutLabel: (keys: Array<string>) => string = shortcutLabel;

    public mounted(): void {
        this.checkSlots();
        this.inDropdown = Boolean(this.$parent.$el?.classList.contains("v-dropdown"));

        if (this.autofocus) {
            this.$nextTick(() => {
                this.focus();
            });
        }
    }

    public updated(): void {
        this.checkSlots();
    }

    public focus(): void {
        const a: HTMLAnchorElement | null = this.$el.querySelector(".v-button-content");
        if (a) {
            a.focus();
        }
    }

    public onFocus(e: FocusEvent): void {
        if (e.type === "focusout" && (e.relatedTarget === null || !this.$el.contains(e.relatedTarget as Node))) {
            this.setActive(false);
        }
        this.$emit(e.type, e);
    }

    public onClick(e: Event): void {
        if (this.dropdown) {
            if (
                this.$el.contains(e.target as Node) &&
                [this.$el, null].indexOf((e.target as Element).closest(".v-button")) !== -1
            ) {
                if (!this.active) {
                    this.setActive(true);
                }
            } else if (this.active) {
                this.setActive(false);
            }
        }

        this.$emit("click", e);
    }

    public onMouse(e: MouseEvent): void {
        if (this.dropdown && this.inDropdown) {
            switch (e.type) {
                case "mouseenter": {
                    this.active = true;
                    break;
                }

                case "mouseleave": {
                    this.active = true;
                    const dropdown: Element | null = (e.relatedTarget as HTMLElement).closest(".v-dropdown");
                    this.active = dropdown !== null && (dropdown.parentNode as Node).isSameNode(this.$el);
                    break;
                }
            }
        }
    }

    public onShortcut(): void {
        const a: HTMLAnchorElement | null = this.$el.querySelector(".v-button-content");
        if (a) {
            a.click();
        }
    }

    public get tag(): string {
        return this.to ? "router-link" : "a";
    }

    private checkSlots(): void {
        Object.assign(this, {
            content: Boolean(this.$slots?.default),
            dropdown: Boolean(this.$slots?.dropdown),
        });
    }

    private setActive(state: boolean): void {
        this.active = state;
        document[state ? "addEventListener" : "removeEventListener"]("click", this.onClick);
    }
}
</script>

<style lang="scss" scoped>
.v-button {
    display: inline-block;
    height: var(--button-height);
    line-height: var(--button-height);
    position: relative;
    vertical-align: middle;

    &[aria-disabled] {
        pointer-events: none;

        > .v-button-content > * {
            opacity: 0.35;
        }
    }

    + .v-button {
        margin-left: 0.5rem;
    }

    > .v-button-content {
        border-radius: 0.2rem;
        align-items: center;
        color: var(--button-color);
        cursor: pointer;
        display: flex;
        height: 100%;
        padding: 0 0.65rem;
        justify-content: center;
        text-decoration: none;
        user-select: none;
        white-space: nowrap;

        span {
            align-items: center;
            display: flex;
        }

        .v-icon {
            height: 1.25rem;

            &.v-button-caret {
                width: auto;
            }

            + span {
                margin-left: 0.5rem;
            }
        }

        .v-button-badge {
            background-color: var(--button-color);
            border-radius: 0.2rem;
            color: var(--background);
            font-size: 0.75rem;
            height: 0.9rem;
            line-height: 0.9rem;
            margin-left: 0.75rem;
            padding: 0 0.25rem;
        }

        .v-button-shortcut {
            margin-left: 2rem;
            opacity: 0.5;
        }

        .v-button-caret {
            margin-left: 0.5rem;
        }

        > * {
            pointer-events: none;
        }
    }

    &.active > .v-button-content,
    > .v-button-content:active,
    > .v-button-content:focus,
    > .v-button-content:hover {
        background-color: var(--button-focus-background);
        color: var(--button-focus-color);
        outline: none;

        .v-button-badge {
            background-color: var(--button-focus-color);
        }
    }

    &.danger > .v-button-content,
    &.primary > .v-button-content {
        color: var(--button-focus-color);
    }

    &.danger {
        .v-button-content {
            background-color: var(--button-danger-background);
        }

        &.active > .v-button-content,
        > .v-button-content:active,
        > .v-button-content:focus,
        > .v-button-content:hover {
            background-color: var(--button-danger-focus-background);
        }
    }

    &.primary {
        > .v-button-content {
            background-color: var(--button-primary-background);
        }

        &.active > .v-button-content,
        > .v-button-content:active,
        > .v-button-content:focus,
        > .v-button-content:hover {
            background-color: var(--button-primary-focus-background);
        }
    }

    &.icon {
        min-width: var(--button-height);
        width: var(--button-height);

        > .v-button-content {
            height: 100%;
            padding: 0;
            width: 100%;

            .v-button-caret {
                display: none;
            }
        }
    }

    > .v-dropdown {
        display: none;
        min-width: 100%;

        &:focus {
            outline: none;
        }
    }

    &.active {
        cursor: pointer;

        > .v-button-content {
            pointer-events: none;
        }

        &.in-dropdown > .v-button-content {
            pointer-events: all;
        }

        > .v-dropdown {
            display: unset;
            pointer-events: all;
        }
    }
}
</style>
