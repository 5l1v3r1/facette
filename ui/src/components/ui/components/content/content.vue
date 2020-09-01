<template>
    <main
        class="v-content"
        :class="{
            [`sidebar-${sidebar.mode}`]: sidebar.active,
            [`toolbar-app-${toolbars.app}`]: toolbars.app,
            [`toolbar-content-${toolbars.content}`]: toolbars.content,
        }"
    >
        <slot></slot>
    </main>
</template>

<script lang="ts">
import {useUI} from "../..";

export default {
    setup(): Record<string, unknown> {
        const ui = useUI();

        return {
            sidebar: ui.state.sidebar,
            toolbars: ui.state.toolbars,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "./mixins";

.v-content {
    min-height: 100vh;
    padding: 2.25rem;
    position: relative;
    transition: margin-left 0.2s var(--timing-function);
    will-change: margin-left;

    &.sidebar-float {
        pointer-events: none;
    }

    &.sidebar-static {
        margin-left: var(--sidebar-width);
    }

    &.toolbar-app-horizontal,
    &.toolbar-content-horizontal {
        margin-top: var(--toolbar-size);
        min-height: calc(100vh - var(--toolbar-size));
    }

    &.toolbar-app-horizontal.toolbar-content-horizontal {
        margin-top: calc(var(--toolbar-size) * 2);
        min-height: calc(100vh - var(--toolbar-size) * 2);
    }

    &.toolbar-content-vertical {
        margin-left: var(--toolbar-size);

        &.sidebar-static {
            margin-left: calc(var(--sidebar-width) + var(--toolbar-size));
        }
    }

    ::v-deep() {
        > .v-spinner {
            bottom: 0;
            left: 0;
            position: absolute;
            right: 0;
            top: 0;
        }

        @include content;
    }
}
</style>
