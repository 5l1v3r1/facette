<template>
    <div class="v-paging" v-if="pages > 1">
        <v-button icon="chevron-left" :disabled="page === 1" @click="switchPage(page - 1)"></v-button>

        <v-button
            :class="{primary: item.value === page}"
            :disabled="item.value === null"
            :key="index"
            @click="switchPage(item.value)"
            v-for="(item, index) in items"
        >
            {{ item.label }}
        </v-button>

        <v-button icon="chevron-right" :disabled="page === pages" @click="switchPage(page + 1)"></v-button>
    </div>
</template>

<script lang="ts">
import {SetupContext, computed} from "vue";

interface PagingItem {
    label: string;
    value: number | null;
}

export default {
    props: {
        adjacency: {
            default: 2,
            type: Number,
        },
        page: {
            required: true,
            type: Number,
        },
        pageSize: {
            required: true,
            type: Number,
        },
        total: {
            required: true,
            type: Number,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const addEllipsis = (items: Array<PagingItem>): void => {
            items.push({label: "…", value: null});
        };

        const addRange = (items: Array<PagingItem>, start: number, end: number): void => {
            for (let i: number = start; i <= end; i++) {
                items.push({label: String(i), value: i});
            }
        };

        const items = computed<Array<PagingItem>>(() => {
            const items: Array<PagingItem> = [];
            const fullAdjacency = props.adjacency * 2;

            if (pages.value <= fullAdjacency * 2) {
                addRange(items, 1, pages.value);
            } else if (props.page <= fullAdjacency) {
                const end = fullAdjacency + props.adjacency + 1;

                addRange(items, 1, end);
                if (end !== pages.value - props.adjacency - 1) {
                    addEllipsis(items);
                }
                addRange(items, pages.value - props.adjacency, pages.value);
            } else if (props.page < pages.value - fullAdjacency) {
                const start = props.page - props.adjacency;
                const end = props.page + props.adjacency;
                const startDelta = props.page === pages.value - fullAdjacency ? 1 : 0;
                const endDelta = props.page === fullAdjacency + 1 ? 0 : 1;

                addRange(items, 1, props.adjacency + startDelta);
                if (endDelta !== 0) {
                    addEllipsis(items);
                }
                addRange(items, start, end);
                if (startDelta === 0 && props.page !== pages.value / 2) {
                    addEllipsis(items);
                }
                addRange(items, pages.value - props.adjacency + endDelta, pages.value);
            } else {
                const start = pages.value - (fullAdjacency + props.adjacency);
                const end = pages.value;

                addRange(items, 1, props.adjacency + 1);
                if (start !== props.adjacency + 2) {
                    addEllipsis(items);
                }
                addRange(items, start, end);
            }

            return items;
        });

        const pages = computed(() => Math.ceil(props.total / props.pageSize));

        const switchPage = (page: number): void => {
            ctx.emit("update:page", page);
        };

        return {
            items,
            pages,
            switchPage,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-paging {
    align-items: center;
    display: flex;
    justify-content: center;
    margin: 3rem 0;
}
</style>
