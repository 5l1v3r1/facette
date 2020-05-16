<template>
    <div
        class="v-grid-item"
        :class="{draggable: !parent.readonly, dragged, placeholder}"
        :draggable="!parent.readonly"
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

        <template v-if="!parent.readonly">
            <v-toolbar @mousedown.native="onMouse" @mouseup.native="onMouse" v-if="!placeholder">
                <v-spacer></v-spacer>
                <v-button
                    icon="times"
                    @click="parent.removeItem(index)"
                    v-tooltip="$t('labels.items.remove')"
                ></v-button>
            </v-toolbar>

            <v-button
                class="v-grid-resize"
                icon="angle-down"
                @mousedown.native="parent.setResizing(index)"
                v-if="!placeholder || parent.state.resizing.item !== null"
            ></v-button>
        </template>

        <slot></slot>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

import GridComponent from "./grid.vue";

@Component
export default class GridItemComponent extends Vue {
    @Prop({default: null, type: Number})
    public index!: number;

    @Prop({required: true, type: Object})
    public layout!: GridItemLayout;

    @Prop({default: false, type: Boolean})
    public placeholder!: boolean;

    public dragged = false;

    private locked = true;

    public onDrag(e: DragEvent): void {
        switch (e.type) {
            case "dragend":
                this.parent.setDragging(null);
                this.dragged = false;
                this.locked = true;

                break;

            case "dragstart":
                if (this.locked) {
                    e.preventDefault();
                } else {
                    (e.dataTransfer as DataTransfer).setData("text/plain", this.index.toString());
                    this.parent.setDragging(this.index);
                    this.dragged = true;
                }

                break;
        }
    }

    public onMouse(e: MouseEvent): void {
        this.locked = e.type !== "mousedown";
    }

    private get parent(): GridComponent {
        return this.$parent as GridComponent;
    }
}
</script>

<style lang="scss" scoped>
.v-grid-item {
    background-color: #272727; // FIXME: don't hardcode here
    border: 0.15rem solid transparent;
    border-radius: 0.2rem;
    grid-column-end: var(--column-end);
    grid-column-start: var(--column-start);
    grid-row-end: var(--row-end);
    grid-row-start: var(--row-start);
    min-width: 0; // This value is need to ensure proper grid resizing
    position: relative;

    &.draggable:hover {
        border-color: var(--toolbar-background);
    }

    &.highlight {
        border-color: var(--blue);
    }

    &.placeholder {
        align-items: center;
        border-radius: 0.2rem;
        display: flex;
        justify-content: center;
        pointer-events: none;
        z-index: 1;

        &.add,
        &.drop,
        &.resize {
            background-color: rgba(2, 136, 209, 0.35);
            border: 0.15rem solid rgba(2, 136, 209, 0.65);
            color: var(--blue);

            .v-label {
                background-color: var(--blue);
                border-radius: 0.2rem;
                color: var(--white);
                line-height: 1.5rem;
                min-height: auto;
                padding: 0 0.5rem;
            }

            .v-grid-resize {
                visibility: visible;

                ::v-deep .v-icon {
                    color: var(--blue);
                }
            }
        }

        &.remove {
            background-color: rgba(211, 47, 47, 0.35);
            border: 0.15rem solid rgba(211, 47, 47, 0.65);
            color: var(--red);
        }

        .v-icon {
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

        ::v-deep {
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
