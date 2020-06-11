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
import {Component, Prop, Vue, Watch} from "vue-property-decorator";

import {ToolbarState} from "@/types/components";

@Component
export default class SidebarComponent extends Vue {
    @Prop({default: true, type: Boolean})
    public active!: boolean;

    @Watch("active", {immediate: true})
    public onActive(to: boolean): void {
        this.$components.set("sidebar", {active: to});
    }

    public get toolbars(): ToolbarState {
        return this.$components.state.toolbars;
    }
}
</script>

<style lang="scss" scoped>
.v-sidebar {
    background-color: var(--sidebar-background);
    bottom: 0;
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

    > .v-button {
        display: block;
        margin: 0;

        ::v-deep > .v-button-content {
            border-radius: 0;
            padding: 0 1.5rem;
            justify-content: flex-start;

            span {
                flex-grow: 1;
            }

            &.active {
                background-color: var(--sidebar-active-background);
                box-shadow: inset -0.25rem 0 0 var(--sidebar-active-shadow);
                color: var(--sidebar-active-color);
            }
        }
    }

    > .v-label {
        color: var(--light-gray);
        letter-spacing: 0.1rem;
        margin: 0.75rem 0 0.25rem;
        padding: 0 1.5rem;
        text-transform: uppercase;

        &:first-of-type {
            margin-top: 0;
        }
    }

    > .v-spinner {
        bottom: 0;
        left: 0;
        position: absolute;
        right: 0;
        top: 0;
    }

    > .v-toolbar {
        min-height: var(--toolbar-size);
    }
}
</style>
