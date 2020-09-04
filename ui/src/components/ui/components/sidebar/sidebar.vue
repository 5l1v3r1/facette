<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <aside
        class="v-sidebar"
        :class="{
            active,
            [`toolbar-app-${toolbars.app}`]: toolbars.app,
            [`toolbar-sidebar-${toolbars.sidebar}`]: toolbars.sidebar,
        }"
    >
        <slot></slot>
    </aside>
</template>

<script lang="ts">
import {watch} from "vue";

import {useUI} from "../..";

export default {
    props: {
        active: {
            default: true,
            type: Boolean,
        },
    },
    setup(props: Record<string, any>): Record<string, unknown> {
        const ui = useUI();

        watch(
            () => props.active,
            to => (ui.state.sidebar.active = to),
            {immediate: true},
        );

        return {
            toolbars: ui.state.toolbars,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-sidebar {
    background-color: var(--sidebar-background);
    bottom: 0;
    color: var(--sidebar-color);
    display: flex;
    flex-direction: column;
    left: 0;
    overflow-y: auto;
    padding: 1.5rem 0;
    position: fixed;
    top: 0;
    transform: translateX(calc(-1 * var(--sidebar-width)));
    transition: transform 0.2s var(--timing-function);
    width: var(--sidebar-width);
    will-change: transform;
    z-index: 200;

    &.active {
        transform: none;
    }

    &.toolbar-app-horizontal {
        top: var(--toolbar-size);
    }

    &.toolbar-app-vertical {
        left: var(--toolbar-size);
    }

    &.toolbar-sidebar-vertical {
        padding-left: var(--toolbar-size);
    }

    ::v-deep() {
        > .v-button {
            display: block;
            margin: 0;

            > .v-button-content {
                border-radius: 0;
                padding: 0 1.5rem;
                justify-content: flex-start;

                span {
                    flex-grow: 1;
                }

                .v-button-badge {
                    min-width: 1.2rem;
                    text-align: center;
                }
            }

            &.active > .v-button-content,
            > .v-button-content.active {
                background-color: var(--sidebar-active-background);
                box-shadow: inset -0.25rem 0 0 var(--accent);
                color: var(--sidebar-active-color);
            }

            > .v-button-content:focus,
            > .v-button-content:hover {
                background-color: var(--sidebar-focus-background);
                color: var(--sidebar-focus-color);
            }
        }

        > :not(.v-toolbar):first-child,
        > .v-toolbar:first-child + * {
            margin-top: 0;
        }

        > .v-label {
            color: var(--light-gray);
            letter-spacing: 0.1rem;
            margin: 0.5rem 0;
            padding: 0 1.5rem;
            text-transform: uppercase;
        }

        > .v-spinner {
            bottom: 0;
            left: 0;
            position: absolute;
            right: 0;
            top: 0;
        }

        > .v-toolbar {
            color: inherit;
            min-height: var(--toolbar-size);
        }
    }
}
</style>
