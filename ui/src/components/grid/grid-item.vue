<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div
        class="v-grid-item"
        :class="{highlight, placeholder}"
        :draggable="!readonly"
        :style="`
        --column-end: ${layout.x + layout.w + 1};
        --column-start: ${layout.x + 1};
        --row-end: ${layout.y + layout.h + 1};
        --row-start: ${layout.y + 1};
        `"
        @dragend="onDrag"
        @dragstart="onDrag"
    >
        <a class="v-grid-anchor" :id="`item${index}`" v-if="index !== null"></a>

        <template v-if="!readonly">
            <v-toolbar @mousedown="onMouse" @mouseup="onMouse" v-if="!placeholder">
                <v-spacer></v-spacer>

                <v-button
                    icon="times"
                    @click.stop="$emit('remove-item')"
                    v-tooltip="i18n.t('labels.items.remove')"
                ></v-button>
            </v-toolbar>

            <v-button
                class="v-grid-resize"
                icon="angle-down"
                @mousedown="$emit('resize-item')"
                v-if="!placeholder"
            ></v-button>
        </template>

        <slot></slot>
    </div>
</template>

<script lang="ts">
import {SetupContext} from "vue";
import {useI18n} from "vue-i18n";

export default {
    props: {
        highlight: {
            default: false,
            type: Boolean,
        },
        index: {
            default: null,
            type: Number,
        },
        layout: {
            required: true,
            type: Object,
        },
        placeholder: {
            default: false,
            type: Boolean,
        },
        readonly: {
            default: false,
            type: Boolean,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const i18n = useI18n();

        let locked = true;

        const onDrag = (ev: DragEvent): void => {
            switch (ev.type) {
                case "dragend":
                    ctx.emit("drag-item", null);
                    locked = true;

                    break;

                case "dragstart":
                    if (locked) {
                        ev.preventDefault();
                    } else {
                        ctx.emit("drag-item", props.index);
                    }

                    break;
            }
        };

        const onMouse = (ev: MouseEvent): void => {
            locked = ev.type !== "mousedown";
        };

        return {
            i18n,
            onDrag,
            onMouse,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-grid-item {
    background-color: #272727; // FIXME: don't hardcode here
    border: 0.15rem solid transparent;
    border-radius: 0.2rem;
    grid-column: var(--column-start) / var(--column-end);
    grid-row: var(--row-start) / var(--row-end);
    min-width: 0; // Needed to ensure proper grid resizing
    position: relative;

    &.highlight {
        border-color: var(--accent);
    }

    &:hover {
        border-color: var(--toolbar-background);
    }

    &.placeholder {
        align-items: center;
        border-radius: 0.2rem;
        display: flex;
        justify-content: center;
        pointer-events: none;
        z-index: 2;

        ::v-deep(.v-icon) {
            font-size: 1.5rem;
        }
    }

    .v-toolbar {
        cursor: move;
        visibility: hidden;
        height: 2.25rem;
        left: 0;
        line-height: 2.25rem;
        position: absolute;
        right: 0;
        top: 0;
        z-index: 1;

        .v-button {
            height: 1.5rem;
            margin: 0 0.375rem;
            min-height: 1.5rem;

            &.icon {
                min-width: 1.5rem;
                width: 1.5rem;
            }
        }
    }

    .v-grid-resize {
        bottom: 0;
        visibility: hidden;
        height: 1.25rem;
        min-width: auto;
        position: absolute;
        right: 0;
        width: 1.25rem;
        z-index: 1;

        ::v-deep() {
            .v-button-content {
                background-color: transparent;
                border-radius: 0;
                color: var(--light-gray);
                cursor: se-resize;
            }

            &.active > .v-button-content,
            &.v-button-content:active,
            &.v-button-content:focus,
            &.v-button-content:hover {
                color: inherit;
            }

            .v-icon {
                font-size: 1rem;
                transform: rotate(-45deg);
            }
        }
    }

    &:hover .v-toolbar,
    &:hover .v-grid-resize {
        visibility: visible;
    }
}
</style>
