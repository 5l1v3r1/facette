<template>
    <table class="v-table">
        <template v-if="hasHeader">
            <thead class="v-table-header">
                <tr class="v-table-row">
                    <v-table-cell align="center" class="v-table-handle" v-if="draggable"></v-table-cell>

                    <v-table-cell align="center" class="v-table-select" v-if="selectable">
                        <v-checkbox v-model:value="selectedAll"></v-checkbox>
                    </v-table-cell>

                    <slot name="header"></slot>
                </tr>
            </thead>
        </template>

        <tbody class="v-table-body">
            <tr
                class="v-table-row"
                :key="index"
                :class="{selected: selectable && selectionState[index]}"
                v-for="(item, index) in value"
            >
                <v-table-cell align="center" class="v-table-handle" v-if="draggable">
                    <v-icon icon="grip-lines-vertical"></v-icon>
                </v-table-cell>

                <v-table-cell align="center" class="v-table-select" v-if="selectable">
                    <v-checkbox
                        @click.prevent="select(index, $event.shiftKey)"
                        v-model:value="selectionState[index]"
                    ></v-checkbox>
                </v-table-cell>

                <slot v-bind="{index: index, value: item}"></slot>
            </tr>
        </tbody>
    </table>
</template>

<script lang="ts">
import {xxHash32} from "js-xxhash";
import {SetupContext, computed, onMounted, ref, unref} from "vue";

// TODO: reimplement draggable

function hash(input: unknown): string {
    return xxHash32(JSON.stringify(input), 0).toString(16);
}

export default {
    props: {
        draggable: {
            default: false,
            type: Boolean,
        },
        selectable: {
            default: false,
            type: Boolean,
        },
        selection: {
            default: (): Array<unknown> => [],
            type: Array as () => Array<unknown>,
        },
        value: {
            required: true,
            type: Array as () => Array<unknown>,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        let selectionLast: number | null = null;

        const hasHeader = ref(false);

        const selectedAll = computed({
            get: () => props.selection.length === props.value.length,
            set: to => selectAll(to),
        });

        const selectionState = computed(() => {
            const hashes = props.selection.map((value: unknown) => hash(value));

            return props.value.reduce((out: Record<number, boolean>, value: unknown, index: number) => {
                out[index] = hashes.includes(hash(value));
                return out;
            }, {});
        });

        const selectAll = (state: boolean): void => {
            ctx.emit("update:selection", state ? props.value : []);
        };

        const select = (index: number, multiple: boolean): void => {
            const state = unref(selectionState);

            if (multiple && selectionLast !== null) {
                const [start, end] = [selectionLast, index].sort((a, b) => a - b);
                for (let i: number = start; i <= end; i++) {
                    state[i] = state[selectionLast];
                }
            }

            ctx.emit(
                "update:selection",
                props.value.filter((value: never, index: number) => state[index]),
            );

            selectionLast = index;
        };

        onMounted(() => {
            hasHeader.value = Boolean(ctx.slots.header);
        });

        return {
            hasHeader,
            select,
            selectAll,
            selectedAll,
            selectionState,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-table {
    border-collapse: collapse;
    border-spacing: 0;
    max-width: 100%;
    width: 100%;

    .v-table-row {
        height: var(--table-row-height);

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

        ::v-deep(.v-table-cell) {
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
                background-color: var(--accent);
                border-bottom: 1px solid var(--light-accent);
                border-top-color: var(--light-accent);
                color: var(--table-row-selected-color);

                ::v-deep() {
                    .v-button {
                        > .v-button-content {
                            color: var(--table-row-selected-color);
                        }

                        .active > .v-button-content,
                        > .v-button-content:active,
                        > .v-button-content:focus,
                        > .v-button-content:hover {
                            background-color: var(--dark-accent);
                        }
                    }

                    .v-checkbox {
                        &:focus .v-icon {
                            background-color: var(--dark-accent);
                            color: var(--table-row-selected-color);
                        }

                        .v-icon {
                            background-color: var(--table-row-selected-color);
                            color: var(--accent);
                        }
                    }
                }
            }
        }
    }
}
</style>
