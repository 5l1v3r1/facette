<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <component
        class="v-toolbar"
        :class="{
            horizontal: !vertical,
            vertical,
            [`clip-${clip}`]: clip,
            [`sidebar-${sidebar.mode}`]: sidebar.active,
            [`toolbar-app-${toolbars.app}`]: toolbars.app,
            [`toolbar-content-${toolbars.content}`]: toolbars.content,
        }"
        :is="tag"
    >
        <slot></slot>
    </component>
</template>

<script lang="ts">
import {computed, onMounted, onUnmounted} from "vue";

import {useUI} from "../..";

export default {
    props: {
        clip: {
            default: null,
            type: String,
        },
        vertical: {
            default: false,
            type: Boolean,
        },
    },
    setup(props: Record<string, any>): Record<string, unknown> {
        const ui = useUI();

        const tag = computed(() => (props.clip === "app" ? "header" : "div"));

        onMounted(() => {
            if (props.clip !== null) {
                ui.state.toolbars[props.clip] = props.vertical ? "vertical" : "horizontal";
            }
        });

        onUnmounted(() => {
            if (props.clip !== null) {
                ui.state.toolbars[props.clip] = null;
            }
        });

        return {
            sidebar: ui.state.sidebar,
            tag,
            toolbars: ui.state.toolbars,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-toolbar {
    align-items: center;
    background-color: var(--toolbar-background);
    color: var(--toolbar-color);
    display: flex;
    line-height: var(--toolbar-size);

    &.horizontal {
        height: var(--toolbar-size);
    }

    &.vertical {
        width: var(--toolbar-size);
    }

    &.clip-app {
        background-color: var(--accent);
        color: white;
        left: 0;
        position: fixed;
        top: 0;
        z-index: 300;

        &.horizontal {
            right: 0;
        }

        &.vertical {
            bottom: 0;
            width: var(--toolbar-size);
        }

        ::v-deep() {
            > .v-button {
                > .v-button-content {
                    color: inherit;

                    > .v-icon .v-icon-badge {
                        --background: var(--accent);

                        background-color: white;
                    }
                }

                &.active > .v-button-content,
                > .v-button-content:active,
                > .v-button-content:focus,
                > .v-button-content:hover {
                    background-color: var(--dark-accent);

                    > .v-icon .v-icon-badge {
                        --background: var(--dark-accent);
                    }
                }
            }

            > .v-divider {
                background-color: white;
                opacity: 0.35;
            }
        }
    }

    &.clip-content {
        top: 0;
        left: 0;
        position: fixed;
        transition: left 0.2s var(--timing-function);
        will-change: left;
        z-index: 100;

        &.sidebar-static {
            left: var(--sidebar-width);
        }

        &.horizontal {
            right: 0;

            &.toolbar-app-horizontal {
                top: var(--toolbar-size);
            }
        }

        &.vertical {
            bottom: 0;
        }
    }

    &.clip-sidebar {
        background-color: var(--sidebar-toolbar-background);
        transform: translateY(-1.5rem);
        z-index: 100;

        &.horizontal {
            border-right: 1px solid var(--sidebar-background);
            position: sticky;
            top: 0;

            &.toolbar-content-horizontal,
            &.toolbar-content-vertical {
                border-right: none;
            }
        }

        &.vertical {
            bottom: 0;
            left: 0;
            position: fixed;
            top: 0;

            &.toolbar-app-horizontal {
                top: var(--toolbar-size);
            }
        }

        ::v-deep() {
            > .v-button {
                &.active > .v-button-content,
                > .v-button-content:active,
                > .v-button-content:focus,
                > .v-button-content:hover {
                    background-color: var(--sidebar-toolbar-focus-background);
                    color: var(--sidebar-toolbar-focus-color);
                }
            }

            > .v-select {
                &.focus,
                &:hover {
                    background-color: var(--sidebar-toolbar-focus-background);
                    color: var(--sidebar-toolbar-focus-color);
                }
            }
        }
    }

    ::v-deep() {
        > .v-divider {
            height: 2rem;
        }

        > .v-button,
        > .v-input,
        > .v-label,
        > .v-message,
        > .v-select {
            height: 2.2rem;
            line-height: 2.2rem;
            margin: 0 0.4rem;
        }

        > .v-button {
            &.icon {
                width: 2.2rem;
            }

            + .v-button {
                margin-left: 0;
            }
        }

        > .v-input {
            border: none;
            background-color: var(--button-focus-background);

            ::placeholder {
                opacity: 0;
            }

            &.focus ::placeholder {
                opacity: var(--placeholder-opacity);
            }
        }

        > .v-message {
            padding: 0 0.65rem;

            .v-icon {
                font-size: 0.9rem;
                margin-right: 0.65rem;
            }
        }

        > .v-select {
            border: none;

            &.focus,
            &:hover {
                background-color: var(--button-focus-background);
            }
        }

        .v-icon {
            --background: var(--toolbar-background);
        }
    }
}
</style>
