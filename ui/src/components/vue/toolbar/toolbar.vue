<template>
    <nav
        class="v-toolbar"
        :class="{
            horizontal: !vertical,
            vertical,
            [`clip-${clip}`]: clip,
            [`sidebar-${sidebar.mode}`]: sidebar.active,
            [`toolbar-app-${toolbars.app}`]: toolbars.app,
            [`toolbar-content-${toolbars.content}`]: toolbars.content,
        }"
    >
        <slot></slot>
    </nav>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

import {SidebarState, ToolbarState} from "@/types/components";

@Component
export default class ToolbarComponent extends Vue {
    @Prop({default: null, type: String})
    public clip!: string;

    @Prop({default: false, type: Boolean})
    public vertical!: boolean;

    public mounted(): void {
        if (this.clip) {
            this.$components.set("toolbars", {[this.clip]: this.vertical ? "vertical" : "horizontal"});
        }
    }

    public destroyed(): void {
        if (this.clip) {
            this.$components.set("toolbars", {[this.clip]: null});
        }
    }

    public get sidebar(): SidebarState {
        return this.$components.state.sidebar;
    }

    public get toolbars(): ToolbarState {
        return this.$components.state.toolbars;
    }
}
</script>

<style lang="scss" scoped>
.v-toolbar {
    align-items: center;
    background-color: var(--toolbar-background);
    display: flex;
    line-height: var(--toolbar-size);

    &.horizontal {
        height: var(--toolbar-size);
    }

    &.vertical {
        width: var(--toolbar-size);
    }

    &.clip-app {
        background-color: var(--blue);
        color: var(--white);
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

        > .v-button {
            ::v-deep > .v-button-content {
                color: inherit;
            }

            &.active ::v-deep > .v-button-content,
            ::v-deep > .v-button-content:active,
            ::v-deep > .v-button-content:focus,
            ::v-deep > .v-button-content:hover {
                background-color: var(--dark-blue);
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
        background-color: var(--background);
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

        > .v-button {
            &.active ::v-deep > .v-button-content,
            ::v-deep > .v-button-content:active,
            ::v-deep > .v-button-content:focus,
            ::v-deep > .v-button-content:hover {
                background-color: var(--toolbar-background);
            }
        }

        > .v-select {
            &.focus,
            &:hover {
                background-color: var(--toolbar-background);
            }
        }
    }

    > .v-divider {
        height: 2rem;
    }

    > .v-button,
    > .v-input,
    > .v-label,
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

        ::v-deep > .v-button-content {
            .v-icon {
                font-size: 0.9rem;
            }
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

    > .v-select {
        border: none;

        &.focus,
        &:hover {
            background-color: var(--button-focus-background);
        }
    }
}
</style>
