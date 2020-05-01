<template>
    <td class="v-table-cell" :class="{[`align-${align}`]: align, grow}">
        <slot></slot>
    </td>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class TableCellComponent extends Vue {
    @Prop({default: "left", type: String})
    public align!: string;

    @Prop({default: false, type: Boolean})
    public grow!: boolean;
}
</script>

<style lang="scss" scoped>
.v-table-cell {
    padding: 0 0.75rem;
    white-space: nowrap;

    &.align-center {
        text-align: center;
    }

    &.align-left {
        text-align: left;
    }

    &.align-right {
        text-align: right;
    }

    &.grow {
        width: 100%;
    }
}
</style>
