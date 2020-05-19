<template>
    <table class="v-table">
        <template v-if="header">
            <thead class="v-table-header">
                <tr class="v-table-row">
                    <v-table-cell align="center" class="v-table-handle" v-if="draggable"></v-table-cell>
                    <v-table-cell align="center" class="v-table-select" v-if="selectable">
                        <v-checkbox @click.native="selectAll(selection.all)" v-model="selection.all"></v-checkbox>
                    </v-table-cell>
                    <slot name="header"></slot>
                </tr>
            </thead>
        </template>

        <draggable class="v-table-body" handle="v-table-handle" tag="tbody" :disabled="!draggable" :list="value">
            <tr
                class="v-table-row"
                :key="index"
                :class="{selected: selectable && selection.items[index]}"
                v-for="(item, index) in value"
            >
                <v-table-cell align="center" class="v-table-handle" v-if="draggable">
                    <v-icon icon="grip-lines-vertical"></v-icon>
                </v-table-cell>
                <v-table-cell align="center" class="v-table-select" v-if="selectable">
                    <v-checkbox
                        @click.native="selectRange($event, index)"
                        v-model="selection.items[index]"
                    ></v-checkbox>
                </v-table-cell>
                <slot v-bind="{index: index, value: item}"></slot>
            </tr>
        </draggable>
    </table>
</template>

<script lang="ts">
import {Component, Prop, Vue, Watch} from "vue-property-decorator";
import draggable from "vuedraggable";

type TableSelection = Record<number, boolean>;

@Component({
    components: {
        draggable,
    },
})
export default class TableComponent extends Vue {
    @Prop({default: false, type: Boolean})
    public draggable!: boolean;

    @Prop({default: false, type: Boolean})
    public selectable!: boolean;

    @Prop({required: true, type: Array})
    public value!: Array<unknown>;

    public header = false;

    public selection: {
        all: boolean;
        items: TableSelection;
        last: number | null;
    } = {
        all: false,
        items: {},
        last: null,
    };

    public mounted(): void {
        this.checkSlots();
    }

    public updated(): void {
        this.checkSlots();
    }

    @Watch("value", {immediate: true})
    public onChange(to: Array<unknown>): void {
        this.selection.items = to.reduce((items: Record<number, boolean>, _: unknown, index: number) => {
            items[index] = false;
            return items;
        }, {});
    }

    @Watch("selection.items", {deep: true})
    public onSelectionItems(): void {
        if (!this.value) {
            this.selection.all = false;
            return;
        }

        const selection: Array<unknown> = this.value.filter((_: unknown, index: number) => this.selection.items[index]);
        this.selection.all = selection.length > 0 && selection.length === this.value.length;
        this.$emit("selection", selection);
    }

    public select(selection: Array<unknown>): void {
        this.selection.items = selection.reduce((items: TableSelection, item: unknown) => {
            const index = this.value.indexOf(item);
            if (index !== -1) {
                items[index] = true;
            }
            return items;
        }, {});
    }

    public selectAll(state: boolean): void {
        this.selection.items = this.value.reduce((selection: TableSelection, _: unknown, index: number) => {
            selection[index] = state;
            return selection;
        }, {});
    }

    public selectRange(e: KeyboardEvent, index: number): void {
        e.preventDefault();

        if (e.shiftKey && this.selection.last !== null) {
            const items: TableSelection = Object.assign({}, this.selection.items);
            const [start, end]: Array<number> = [this.selection.last, index].sort((a, b) => a - b);

            for (let i: number = start; i <= end; i++) {
                items[i] = this.selection.items[this.selection.last];
            }

            this.selection.items = items;
        }

        this.selection.last = index;
    }

    private checkSlots(): void {
        this.header = Boolean(this.$slots?.header);
    }
}
</script>

<style lang="scss" scoped>
.v-table {
    border-collapse: collapse;
    border-spacing: 0;
    max-width: 100%;
    width: 100%;

    .v-table-row {
        height: 2.5rem;

        .v-table-handle {
            color: var(--table-handle-color);
            cursor: move;
            padding-right: 0;
            width: 1.5rem;
        }

        .v-table-select {
            padding-right: 0;
            user-select: none;
            width: 2.5rem;
        }
    }

    .v-table-header {
        .v-table-row {
            border-top: none;
        }

        .v-table-cell {
            color: var(--table-header-color);
            text-transform: uppercase;
        }
    }

    .v-table-body {
        .v-table-row {
            border-top: 1px solid var(--table-row-border);

            &:last-child {
                border-bottom: 1px solid var(--table-row-border);
            }

            &.selected {
                background-color: var(--table-row-selected-background);
                border-bottom: 1px solid var(--table-row-selected-border);
                border-top-color: var(--table-row-selected-border);
                color: var(--table-row-selected-color);

                ::v-deep {
                    .v-button {
                        > .v-button-content {
                            color: var(--table-row-selected-color);
                        }

                        .active > .v-button-content,
                        > .v-button-content:active,
                        > .v-button-content:focus,
                        > .v-button-content:hover {
                            background-color: var(--table-row-selected-form);
                        }
                    }

                    .v-checkbox {
                        &:focus .v-icon {
                            background-color: var(--table-row-selected-form);
                            color: var(--table-row-selected-color);
                        }

                        .v-icon {
                            background-color: var(--table-row-selected-color);
                            color: var(--table-row-selected-background);
                        }
                    }
                }
            }
        }
    }
}
</style>
