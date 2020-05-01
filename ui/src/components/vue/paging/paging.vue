<template>
    <div class="v-paging" v-if="pages > 1">
        <v-button icon="chevron-left" :disabled="page === 1" @click="handler(page - 1)"></v-button>
        <v-button
            :class="{primary: item.value === page}"
            :disabled="item.value === null"
            :key="index"
            @click="handler(item.value)"
            v-for="(item, index) in items"
            >{{ item.label }}</v-button
        >
        <v-button icon="chevron-right" :disabled="page === pages" @click="handler(page + 1)"></v-button>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue, Watch} from "vue-property-decorator";

interface PagingItem {
    label: string;
    value: number | null;
}

@Component
export default class PagingComponent extends Vue {
    @Prop({default: 2, type: Number})
    public adjacency!: number;

    @Prop({required: true, type: Function})
    public handler!: (page: number) => void;

    @Prop({required: true, type: Number})
    public page!: number;

    @Prop({required: true, type: Number})
    public pageSize!: number;

    @Prop({required: true, type: Number})
    public total!: number;

    public items: Array<PagingItem> = [];

    public pages = 0;

    public addEllipsis(items: Array<PagingItem>): void {
        items.push({label: "…", value: null});
    }

    public addRange(items: Array<PagingItem>, start: number, end: number): void {
        for (let i: number = start; i <= end; i++) {
            items.push({label: String(i), value: i});
        }
    }

    @Watch("page", {immediate: true})
    public onPage(): void {
        this.update();
    }

    @Watch("total")
    public onTotal(): void {
        this.update();
    }

    private update(): void {
        const items: Array<PagingItem> = [];
        const fullAdjacency = this.adjacency * 2;

        this.pages = Math.ceil(this.total / this.pageSize);

        if (this.pages <= fullAdjacency * 2) {
            this.addRange(items, 1, this.pages);
        } else if (this.page <= fullAdjacency) {
            const end = fullAdjacency + this.adjacency + 1;

            this.addRange(items, 1, end);
            if (end !== this.pages - this.adjacency - 1) {
                this.addEllipsis(items);
            }
            this.addRange(items, this.pages - this.adjacency, this.pages);
        } else if (this.page < this.pages - fullAdjacency) {
            const start = this.page - this.adjacency;
            const end = this.page + this.adjacency;
            let startDelta = 0;
            let endDelta = 1;

            if (this.page === this.pages - fullAdjacency) {
                startDelta = 1;
            }

            if (this.page === fullAdjacency + 1) {
                endDelta = 0;
            }

            this.addRange(items, 1, this.adjacency + startDelta);
            if (endDelta !== 0) {
                this.addEllipsis(items);
            }
            this.addRange(items, start, end);
            if (startDelta === 0 && this.page !== this.pages / 2) {
                this.addEllipsis(items);
            }
            this.addRange(items, this.pages - this.adjacency + endDelta, this.pages);
        } else {
            const start = this.pages - (fullAdjacency + this.adjacency);
            const end = this.pages;

            this.addRange(items, 1, this.adjacency + 1);
            if (start !== this.adjacency + 2) {
                this.addEllipsis(items);
            }
            this.addRange(items, start, end);
        }

        this.items = items;
    }
}
</script>

<style lang="scss" scoped>
.v-paging {
    align-items: center;
    display: flex;
    justify-content: center;
    margin: 2rem 0;
}
</style>
