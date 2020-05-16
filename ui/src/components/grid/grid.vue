<template>
    <div
        class="v-grid"
        :aria-readonly="readonly"
        :class="{
            dragging: state.dragging !== null,
            editable: !readonly,
            hovering: state.hovering,
            resizing: state.resizing.item !== null || state.resizing.row !== null,
        }"
        :style="`
            --columns: ${layout.columns || 1};
            --row-height: ${layout.rowHeight || 260}px;
            --rows: ${layout.rows || 1};
        `"
    >
        <template v-if="!readonly">
            <div class="v-grid-handle columns">
                <template v-if="layout.columns > 1">
                    <v-button
                        icon="minus"
                        :key="index"
                        :style="`grid-column-start: ${index}`"
                        @click="shrinkLayout('columns', index - 1)"
                        @focusin.native="toggleRemovePlaceholder('columns', index - 1, true)"
                        @focusout.native="toggleRemovePlaceholder('columns', index - 1, false)"
                        @mouseenter.native="toggleRemovePlaceholder('columns', index - 1, true)"
                        @mouseleave.native="toggleRemovePlaceholder('columns', index - 1, false)"
                        v-for="index in layout.columns"
                        v-show="removable.x.includes(index - 1)"
                    ></v-button>
                </template>
            </div>

            <div class="v-grid-handle rows">
                <v-button
                    class="v-grid-height"
                    icon="grip-lines"
                    :style="`top: ${state.resizing.row || layout.rowHeight}px`"
                    @dragstart.native.prevent
                    @mousedown.native="setRowResizing(true)"
                    @mouseup.native="setRowResizing(false)"
                ></v-button>

                <template v-if="layout.rows > 1">
                    <v-button
                        icon="minus"
                        :key="index"
                        :style="`grid-row-start: ${index}`"
                        @click="shrinkLayout('rows', index - 1)"
                        @focusin.native="toggleRemovePlaceholder('rows', index - 1, true)"
                        @focusout.native="toggleRemovePlaceholder('rows', index - 1, false)"
                        @mouseenter.native="toggleRemovePlaceholder('rows', index - 1, true)"
                        @mouseleave.native="toggleRemovePlaceholder('rows', index - 1, false)"
                        v-for="index in layout.rows"
                        v-show="removable.y.includes(index - 1)"
                    ></v-button>
                </template>
            </div>

            <div class="v-grid-grow columns">
                <v-button icon="plus" @click="growLayout('rows')"></v-button>
            </div>

            <div class="v-grid-grow rows">
                <v-button icon="plus" @click="growLayout('columns')"></v-button>
            </div>
        </template>

        <div class="v-grid-container" ref="container" tabindex="-1">
            <v-grid-item
                placeholder
                :class="{
                    add: state.hovering,
                    drop: state.dragging !== null,
                    remove: state.removing,
                    resize: state.resizing.item !== null,
                }"
                :layout="placeholder"
                v-if="placeholder"
            >
                <v-icon icon="trash" v-if="state.removing"></v-icon>
                <v-icon icon="plus" v-else-if="state.hovering"></v-icon>
                <v-icon icon="hand-point-right" v-else-if="state.dragging !== null"></v-icon>

                <v-label v-else>{{ placeholder.w }} x {{ placeholder.h }}</v-label>
            </v-grid-item>

            <v-grid-item
                :class="{highlight: index === highlightIndex}"
                :index="index"
                :key="index"
                :layout="item.layout"
                v-for="(item, index) in value"
            >
                <slot v-bind="{index, value: item}"></slot>
            </v-grid-item>
        </div>
    </div>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import debounce from "lodash/debounce";
import ResizeObserver from "resize-observer-polyfill";
import {Component, Mixins, Prop, Watch} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

interface ItemWithLayout {
    layout: GridItemLayout;
}

interface Position {
    x: number;
    y: number;
}

const minRowHeight = 160;

@Component
export default class GridComponent extends Mixins<CustomMixins>(CustomMixins) {
    @Prop({default: null, type: Function})
    public addHandler!: ((layout: GridItemLayout) => void) | undefined;

    @Prop({default: null, type: Number})
    public highlightIndex!: number;

    @Prop({required: true, type: Object})
    public layout!: GridLayout;

    @Prop({default: false, type: Boolean})
    public readonly!: boolean;

    @Prop({required: true, type: Array})
    public value!: Array<unknown>;

    public columnWidth = 0;

    public placeholder: GridItemLayout | null = null;

    public removable: {
        x: Array<number>;
        y: Array<number>;
    } = {
        x: [],
        y: [],
    };

    public state: {
        dragging: number | null;
        hovering: boolean;
        removing: boolean;
        resizing: {
            item: number | null;
            row: number | null;
        };
    } = {
        dragging: null,
        hovering: false,
        removing: false,
        resizing: {
            item: null,
            row: null,
        },
    };

    private domRect: DOMRect | null = null;

    private gridGap = 0;

    private matrix: Array<Array<number | null>> | null = null;

    private position: {
        x: number;
        y: number;
    } | null = null;

    private resize: ResizeObserver | null = null;

    private unwatchChange: (() => void) | null = null;

    public mounted(): void {
        if (!this.readonly) {
            const el = this.$el as HTMLElement;
            el.addEventListener("keydown", this.onKeydown);
            el.addEventListener("mousemove", this.onMouse);
            el.addEventListener("mouseleave", this.onMouse);

            if (this.addHandler) {
                el.addEventListener("click", this.onMouse);
            }

            this.unwatchChange = this.$watch(() => [this.layout, this.value], this.onChange, {
                deep: true,
                immediate: true,
            });

            const container = this.$refs.container as HTMLElement;

            // Observer container resizing for DOMRect update (useful for
            // position computation).
            this.resize = new ResizeObserver(
                debounce(() => {
                    this.domRect = container.getBoundingClientRect();
                }, 200),
            );

            this.resize.observe(container);

            // Get computed grid gap for future coordinate computations
            this.gridGap = parseFloat(getComputedStyle(container).rowGap);

            this.$nextTick(() => {
                this.onLayout(this.layout);
            });
        }
    }

    public beforeDestroy(): void {
        if (!this.readonly) {
            const el = this.$el as HTMLElement;
            el.removeEventListener("keydown", this.onKeydown);
            el.removeEventListener("mousemove", this.onMouse);
            el.removeEventListener("mouseleave", this.onMouse);

            if (this.addHandler) {
                el.removeEventListener("click", this.onMouse);
            }

            if (this.unwatchChange) {
                this.unwatchChange();
            }

            if (this.resize) {
                this.resize.disconnect();
            }
        }
    }

    public growLayout(key: "columns" | "rows"): void {
        this.layout[key]++;
    }

    public removeItem(index: number): void {
        const value = cloneDeep(this.value);
        value.splice(index, 1);

        this.$emit("input", value);
    }

    public scrollTo(id: string): void {
        const el: HTMLElement | null = document.getElementById(id);
        if (el !== null) {
            el.scrollIntoView(true);
        }
    }

    public setDragging(index: number | null): void {
        this.state.dragging = index;

        if (this.state.dragging !== null) {
            this.$el.addEventListener("dragover", this.onDragover);
            this.$el.addEventListener("drop", this.onDrop);
            this.placeholder = Object.assign({}, (this.value[this.state.dragging] as ItemWithLayout).layout);
        } else {
            this.$el.removeEventListener("dragover", this.onDragover);
            this.$el.removeEventListener("drop", this.onDrop);
            this.placeholder = null;
        }
    }

    public setResizing(index: number | null): void {
        this.state.resizing.item = index;

        if (this.state.resizing.item !== null) {
            this.$el.addEventListener("mouseup", this.onMouse);
            this.placeholder = Object.assign({}, (this.value[this.state.resizing.item] as ItemWithLayout).layout);
        } else {
            this.$el.removeEventListener("mouseup", this.onMouse);
            this.placeholder = null;
        }
    }

    public setRowResizing(state: boolean): void {
        if (state) {
            this.$el.addEventListener("mouseup", this.onMouse);
            this.state.resizing.row = this.layout.rowHeight;
        } else {
            this.$el.removeEventListener("mouseup", this.onMouse);
            this.state.resizing.row = null;
        }
    }

    public shrinkLayout(key: "columns" | "rows", index: number): void {
        const value = cloneDeep(this.value);

        value.forEach(value => {
            const v = value as ItemWithLayout;

            if (key === "columns" && v.layout.x > index) {
                v.layout.x--;
            } else if (key === "rows" && v.layout.y > index) {
                v.layout.y--;
            }
        });

        this.$emit("input", value);

        if (key === "columns") {
            this.layout.columns--;
        } else if (key === "rows") {
            this.layout.rows--;
        }

        this.placeholder = null;
    }

    public toggleRemovePlaceholder(key: "columns" | "rows", index: number, state: boolean): void {
        if (!state) {
            this.placeholder = null;
            this.state.removing = false;
            return;
        }

        if (key === "columns") {
            this.placeholder = {x: index, y: 0, w: 1, h: this.layout.rows};
        } else {
            this.placeholder = {x: 0, y: index, w: this.layout.columns, h: 1};
        }

        this.state.removing = true;
    }

    private getCollisions(layout: GridItemLayout): Array<GridItemLayout> {
        return this.value
            .reduce((layouts: Array<GridItemLayout>, item: unknown) => {
                const cur = (item as ItemWithLayout).layout;

                if (
                    !(
                        (
                            cur === layout || // is same node
                            cur.x + cur.w <= layout.x || // is at left of node
                            cur.x >= layout.x + layout.w || // is at right of node
                            cur.y + cur.h <= layout.y || // is above node
                            cur.y >= layout.y + layout.h
                        ) // is below node
                    )
                ) {
                    layouts.push(cur);
                }

                return layouts;
            }, [])
            .sort((a: GridItemLayout, b: GridItemLayout) => b.y - a.y); // Sort by Y to start with farthest collision
    }

    private getPosition(e: MouseEvent): Position | null {
        if (this.domRect === null || (e.target as HTMLElement).closest(".v-grid-container") === null) {
            return null;
        }

        const x = Math.floor((e.pageX - this.domRect.left) / (this.columnWidth + this.gridGap));
        const y = Math.floor((e.pageY - this.domRect.top) / (this.layout.rowHeight + this.gridGap));

        return {
            x: x < 0 ? 0 : Math.min(x, this.layout.columns - 1),
            y: y < 0 ? 0 : Math.min(y, this.layout.rows - 1),
        };
    }

    private move(layout: GridItemLayout, x: number, y: number): void {
        Object.assign(layout, {x, y});

        const yDelta: number = layout.y + layout.h - this.layout.rows;
        if (yDelta > 0) {
            this.layout.rows += yDelta;
        }

        this.getCollisions(layout).forEach(l => this.move(l, l.x, layout.y + layout.h));
    }

    private onChange(): void {
        const matrix: Array<Array<number>> = Array.from({length: this.layout.rows}, () =>
            Array(this.layout.columns).fill(null),
        );

        this.value.forEach((item: unknown, index: number) => {
            const layout: GridItemLayout = (item as ItemWithLayout).layout;
            for (let y = layout.y; y < layout.y + layout.h; y++) {
                for (let x = layout.x; x < layout.x + layout.w; x++) {
                    if (matrix[y][x] === null) {
                        matrix[y][x] = index;
                    }
                }
            }
        });

        const removable: {x: Array<number>; y: Array<number>} = {x: [], y: []};

        for (let i = 0; i < this.layout.columns; i++) {
            if (matrix.filter(v => v[i] === null).length === this.layout.rows) {
                removable.x.push(i);
            }
        }

        for (let i = 0; i < this.layout.rows; i++) {
            if (matrix[i].filter(v => v === null).length === this.layout.columns) {
                removable.y.push(i);
            }
        }

        Object.assign(this, {matrix, removable});
    }

    private onDragover(e: Event): void {
        e.preventDefault();

        const de: DragEvent = e as DragEvent;
        const index = parseInt((de.dataTransfer as DataTransfer).getData("text/plain"), 10);
        const layout = (this.value[index] as ItemWithLayout).layout;
        const position = this.getPosition(de);

        if (position === null) {
            return;
        }

        const xDelta = this.layout.columns - (position.x + layout.w);
        if (xDelta < 0) {
            position.x += xDelta;
        }

        const yDelta = this.layout.rows - (position.y + layout.h);
        if (yDelta < 0) {
            position.y += yDelta;
        }

        if (this.placeholder === null || position.x !== this.placeholder.x || position.y !== this.placeholder.y) {
            Object.assign(this.placeholder, {x: position.x, y: position.y});
        }
    }

    private onDrop(e: Event): void {
        e.preventDefault();

        if (this.placeholder === null) {
            return;
        }

        const de: DragEvent = e as DragEvent;
        const index: number = parseInt((de.dataTransfer as DataTransfer).getData("text/plain"), 10);
        const layout = (this.value[index] as ItemWithLayout).layout;

        // Only update layout if position has changed
        if (this.placeholder.x !== layout.x || this.placeholder.y !== layout.y) {
            this.move(layout, this.placeholder.x, this.placeholder.y);
            this.reorder();
        }
    }

    private onKeydown(e: KeyboardEvent): void {
        if (e.code === "Escape") {
            if (this.state.resizing.item !== null) {
                this.setResizing(null);
            } else if (this.state.resizing.row !== null) {
                this.setRowResizing(false);
            }
        }
    }

    @Watch("layout", {deep: true})
    public onLayout(to: GridLayout): void {
        this.columnWidth = this.$el ? (this.$el.clientWidth - this.gridGap * (to.columns + 1)) / to.columns : 0;
    }

    private onMouse(e: Event): void {
        // Don't handle any mouse event while dragging a grid item (handled via
        // onDragover).
        if (this.state.dragging) {
            return;
        }

        const resizingItem: GridItemLayout | null =
            this.state.resizing.item !== null ? (this.value[this.state.resizing.item] as ItemWithLayout).layout : null;

        switch (e.type) {
            case "click":
                if (this.placeholder !== null && this.state.hovering && this.addHandler) {
                    this.addHandler(this.placeholder);
                }

                break;

            case "mouseleave":
                // Cursor has moved outside of the grid area, reset position and
                // cancel hovering.
                this.position = null;
                if (this.state.hovering) {
                    this.placeholder = null;
                    this.state.hovering = false;
                }

                break;

            case "mousemove": {
                const me = e as MouseEvent;

                // Check whether or not we are resizing grid rows height
                if (this.state.resizing.row !== null) {
                    if (this.domRect !== null) {
                        const height = me.pageY - this.domRect.top;
                        if (height >= minRowHeight) {
                            this.state.resizing.row = height;
                        }
                    }

                    break;
                }

                const position = this.getPosition(me);
                if (position === null) {
                    if (this.state.hovering) {
                        this.state.hovering = false;
                        this.placeholder = null;
                    }

                    this.position = null;

                    break;
                }

                // Check whether or not we only need to handle item resizing
                if (resizingItem !== null) {
                    const w = position.x - resizingItem.x + 1;
                    const h = position.y - resizingItem.y + 1;

                    if (
                        this.placeholder === null ||
                        (w !== this.placeholder.w && w >= 1) ||
                        (h !== this.placeholder.h && h >= 1)
                    ) {
                        this.placeholder = {x: resizingItem.x, y: resizingItem.y, w, h};
                    }

                    break;
                }

                // Stop if position hasn't change
                if (this.position !== null && position.x === this.position.x && position.y === this.position.y) {
                    break;
                }

                this.position = position;

                // Only trigger hovering if cursor is moving on top of an empty
                // grid cell.
                if (this.matrix?.[this.position.y]?.[this.position.x] === null) {
                    if (!this.state.hovering) {
                        this.state.hovering = true;
                    }
                    this.updatePlaceholder();
                } else {
                    this.state.hovering = false;
                    this.placeholder = null;
                }

                // Ensure removing flag is reset (useful when accessing handles
                // using keyboard)
                this.state.removing = false;

                break;
            }

            case "mouseup":
                if (this.state.resizing.row !== null) {
                    this.layout.rowHeight = this.state.resizing.row;
                    this.setRowResizing(false);
                    break;
                }

                if (resizingItem === null || this.placeholder === null) {
                    break;
                }

                if (resizingItem.w !== this.placeholder.w || resizingItem.h !== this.placeholder.h) {
                    Object.assign(resizingItem, {w: this.placeholder.w, h: this.placeholder.h});
                    this.getCollisions(resizingItem).forEach(l => this.move(l, l.x, resizingItem.y + resizingItem.h));
                    this.reorder();
                }

                this.setResizing(null);

                break;
        }
    }

    private reorder(): void {
        if (this.matrix === null) {
            return;
        }

        const value: Array<unknown> = [];
        const found: Array<number> = [];

        for (let y = 0; y < this.layout.rows; y++) {
            for (let x = 0; x < this.layout.columns; x++) {
                const index = this.matrix?.[y]?.[x] ?? null;
                if (index !== null && !found.includes(index)) {
                    value.push(this.value[index]);
                    found.push(index);
                }
            }
        }

        this.$emit("input", value);
    }

    @Watch("modifiers", {deep: true})
    private updatePlaceholder(): void {
        if (this.position === null || !this.state.hovering) {
            return;
        }

        const [x, y] = Object.values(this.position);

        if (this.modifiers.alt) {
            if (this.modifiers.shift) {
                const boundaries = [y, y];
                while (boundaries[0] > 0 && this.matrix?.[boundaries[0] - 1]?.[x] === null) {
                    boundaries[0]--;
                }
                while (boundaries[1] < this.layout.columns && this.matrix?.[boundaries[1] + 1]?.[x] === null) {
                    boundaries[1]++;
                }

                this.placeholder = {x, y: boundaries[0], w: 1, h: boundaries[1] - boundaries[0] + 1};
            } else {
                const boundaries = [x, x];
                while (boundaries[0] > 0 && this.matrix?.[y]?.[boundaries[0] - 1] === null) {
                    boundaries[0]--;
                }
                while (boundaries[1] < this.layout.columns && this.matrix?.[y]?.[boundaries[1] + 1] === null) {
                    boundaries[1]++;
                }

                this.placeholder = {x: boundaries[0], y, w: boundaries[1] - boundaries[0] + 1, h: 1};
            }
        } else {
            this.placeholder = {x, y, w: 1, h: 1};
        }
    }
}
</script>

<style lang="scss" scoped>
.v-grid {
    position: relative;

    &.editable {
        margin: -1rem;
        padding: 1rem;
    }

    &.dragging,
    &.resizing {
        .v-grid-item {
            pointer-events: none;
        }
    }

    .v-grid-handle,
    .v-grid-grow {
        align-items: center;
        background-color: var(--toolbar-background);
        line-height: 1rem;
        position: absolute;

        &.columns {
            height: 1rem;
            left: 1rem;
            right: 1rem;
        }

        &.rows {
            bottom: 1rem;
            flex-direction: column;
            top: 1rem;
            width: 1rem;
        }

        .v-button {
            font-size: 0.65rem;
            height: inherit;
            margin: 0;
            min-width: auto;
            width: inherit;

            ::v-deep .v-button-content {
                border-radius: 0;
                color: var(--light-gray);
            }
        }
    }

    &:hover .v-grid-handle,
    &:hover .v-grid-grow {
        opacity: 1;
    }

    .v-grid-handle {
        display: grid;

        &.columns {
            column-gap: 0.5rem;
            grid-template-columns: repeat(var(--columns), 1fr);
            top: -0.5rem;
        }

        &.rows {
            grid-template-rows: repeat(var(--rows), var(--row-height));
            left: -0.5rem;
            row-gap: 0.5rem;

            .v-button {
                height: 100%;
            }

            .v-grid-height {
                height: 0.5rem;
                position: absolute;

                ::v-deep .v-button-content {
                    cursor: ns-resize;
                }
            }
        }
    }

    &.resizing {
        cursor: ns-resizex;

        .v-grid-height {
            pointer-events: none;
        }
    }

    .v-grid-grow {
        display: flex;

        &.columns {
            bottom: -0.5rem;
        }

        &.rows {
            right: -0.5rem;
        }

        .v-button {
            flex-grow: 1;
        }
    }

    .v-grid-container {
        column-gap: 0.5rem;
        display: grid;
        grid-template-columns: repeat(var(--columns), 1fr);
        grid-template-rows: repeat(var(--rows), var(--row-height));
        row-gap: 0.5rem;

        &:focus {
            outline: none;
        }
    }

    &.editable .v-grid-container {
        box-shadow: 0 0 0 0.5rem transparent;
    }

    &.hovering .v-grid-container {
        cursor: pointer;
    }
}
</style>
