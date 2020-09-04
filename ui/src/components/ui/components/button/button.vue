<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div
        class="v-button"
        ref="el"
        :aria-disabled="disabled || undefined"
        :class="{
            active,
            danger,
            primary,
            icon: !hasLabel && (!hasDropdown || icon.startsWith('angle-')),
            ['in-dropdown']: inDropdown,
        }"
        @autofocus="focus()"
        @shortcut="!disabled && onShortcut()"
    >
        <component
            class="v-button-content"
            ref="content"
            role="button"
            :is="tag"
            :tabindex="disabled ? -1 : 0"
            :target="target"
            :to="to"
            @click.passive="onClick"
            @focusout="onFocusout"
            @keydown.enter="onKeydownEnter"
            @mouseenter="hasDropdown && inDropdown && onMouse($event)"
            @mouseleave="hasDropdown && inDropdown && onMouse($event)"
        >
            <v-icon aria-hidden="true" :badge="iconBadge" :icon="icon" :spin="iconSpin" v-if="icon !== null"></v-icon>

            <template v-if="hasLabel">
                <div class="v-button-label">
                    <slot></slot>
                </div>

                <div class="v-button-badge" v-if="badge">{{ badge }}</div>
            </template>

            <div class="v-button-shortcut" v-if="shortcutsEnabled && inDropdown && shortcut">
                {{ shortcutLabel(shortcut) }}
            </div>

            <v-icon
                class="v-button-caret"
                :icon="inDropdown ? 'angle-right' : 'angle-down'"
                v-if="hasDropdown && !icon.startsWith('angle-')"
            ></v-icon>
        </component>

        <v-dropdown
            tabindex="-1"
            :anchor="dropdownAnchor"
            @focusout="onFocusout"
            @mouseleave="inDropdown && onMouse($event)"
            v-if="hasDropdown"
        >
            <slot name="dropdown"></slot>
        </v-dropdown>
    </div>
</template>

<script lang="ts">
import {SetupContext, computed, onMounted, ref} from "vue";

import {useUI} from "../..";
import {shortcutLabel} from "../../directives/shortcut";

export default {
    props: {
        badge: {
            default: null,
            type: [Number, String],
        },
        danger: {
            default: false,
            type: Boolean,
        },
        disabled: {
            default: false,
            type: Boolean,
        },
        dropdownAnchor: {
            default: "bottom-left",
            type: String,
        },
        href: {
            default: null,
            type: String,
        },
        icon: {
            default: null,
            type: String,
        },
        iconBadge: {
            default: null,
            type: [Number, String],
        },
        iconSpin: {
            default: false,
            type: Boolean,
        },
        primary: {
            default: false,
            type: Boolean,
        },
        target: {
            default: null,
            type: String,
        },
        to: {
            default: null,
            type: Object,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const ui = useUI();

        const active = ref(false);
        const content = ref<HTMLAnchorElement | null>(null);
        const el = ref<HTMLElement | null>(null);
        const hasDropdown = ref(false);
        const hasLabel = ref(false);
        const inDropdown = ref(false);

        const shortcut = computed(() => el.value?.dataset.vShortcut);
        const shortcutsEnabled = computed(() => ui.state.shortcuts.enabled);
        const tag = computed(() => (props.to ? "router-link" : "a"));

        const focus = (): void => {
            content.value?.focus();
        };

        const onClick = (ev: MouseEvent): void => {
            if (inDropdown.value) {
                // Ensure dropdown parent button is focused out (i.e. close
                // dropdown).
                (ev.target as HTMLElement)
                    .closest(".v-dropdown")
                    ?.parentNode?.querySelector(".v-button-content")
                    ?.dispatchEvent(new Event("focusout"));
            } else if (hasDropdown.value) {
                onDropdown(ev);
            }
        };

        const onDropdown = (ev: Event): void => {
            // Stop propagation needs to be handle in function instead of as a
            // modifier to avoid stopping on non-dropdown buttons.
            ev.stopPropagation();

            if (el.value?.contains(ev.target as Node) && el.value === (ev.target as Element).closest(".v-button")) {
                active.value = !active.value;
            } else if (active.value) {
                active.value = false;
            }
        };

        const onFocusout = (ev: FocusEvent): void => {
            if (ev.relatedTarget === null || !el.value?.contains(ev.relatedTarget as Node)) {
                active.value = false;
            }
        };

        const onKeydownEnter = (ev: KeyboardEvent): void => {
            if (hasDropdown.value) {
                onDropdown(ev);
            } else {
                el.value?.dispatchEvent(new Event("click"));
            }
        };

        const onMouse = (ev: MouseEvent): void => {
            if (hasDropdown.value && inDropdown.value) {
                switch (ev.type) {
                    case "mouseenter": {
                        active.value = true;
                        break;
                    }

                    case "mouseleave": {
                        const dropdown = (ev.relatedTarget as HTMLElement).closest(".v-dropdown");
                        active.value = dropdown !== null && (dropdown.parentNode as Node).isSameNode(el.value);
                        break;
                    }
                }
            }
        };

        const onShortcut = (): void => {
            if (props.href || props.to) {
                el.value?.querySelector(".v-button-content")?.dispatchEvent(new Event("click"));
            } else {
                ctx.emit("click");
            }
        };

        onMounted(() => {
            if (props.href && content.value) {
                content.value.href = props.href;
            }

            hasDropdown.value = Boolean(ctx.slots.dropdown);
            hasLabel.value = Boolean(ctx.slots.default);
            inDropdown.value = el.value?.closest(".v-dropdown") !== null;
        });

        return {
            active,
            content,
            el,
            focus,
            hasLabel,
            hasDropdown,
            inDropdown,
            onClick,
            onDropdown,
            onFocusout,
            onKeydownEnter,
            onMouse,
            onShortcut,
            shortcut,
            shortcutLabel,
            shortcutsEnabled,
            tag,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-button {
    display: inline-block;
    height: var(--button-height);
    line-height: var(--button-height);
    position: relative;
    vertical-align: middle;

    &[aria-disabled="true"] {
        pointer-events: none;

        > .v-button-content > * {
            opacity: 0.425;
        }
    }

    + .v-button {
        margin-left: 0.5rem;
    }

    > .v-button-content {
        border-radius: 0.2rem;
        align-items: center;
        cursor: pointer;
        display: flex;
        height: 100%;
        padding: 0 0.65rem;
        justify-content: center;
        text-decoration: none;
        user-select: none;
        white-space: nowrap;

        .v-icon:first-child {
            font-size: 0.9rem;
            height: 1.25rem;
            margin-right: 0.5rem;

            &.v-button-caret {
                width: auto;
            }

            + .v-button-caret {
                margin-left: 0;
            }
        }

        .v-button-label {
            align-items: center;
            display: flex;
            flex-grow: 1;
        }

        .v-button-badge {
            background-color: var(--color);
            border-radius: 0.2rem;
            color: var(--background);
            font-size: 0.75rem;
            height: 1rem;
            line-height: 1rem;
            margin-left: 0.75rem;
            opacity: 0.85;
            padding: 0 0.25rem;
        }

        .v-button-shortcut {
            font-size: 0.8rem;
            margin-left: 2rem;
            opacity: 0.425;
        }

        .v-button-caret {
            margin-left: 0.5rem;
            opacity: 0.75;
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

        .v-icon {
            --background: var(--button-focus-background);
        }
    }

    &.danger > .v-button-content,
    &.primary > .v-button-content {
        color: white;
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
            background-color: var(--accent);
        }

        &.active > .v-button-content,
        > .v-button-content:active,
        > .v-button-content:focus,
        > .v-button-content:hover {
            background-color: var(--light-accent);
        }
    }

    &.icon {
        min-width: var(--button-height);
        width: var(--button-height);

        > .v-button-content {
            height: 100%;
            padding: 0;
            width: 100%;

            .v-icon {
                margin-right: 0;
            }

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
