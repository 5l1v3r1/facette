<template>
    <td class="v-table-cell" :class="{[`align-${align}`]: align, grow}">
        <slot></slot>
    </td>
</template>

<script lang="ts">
export default {
    props: {
        align: {
            default: "left",
            type: String,
        },
        grow: {
            default: false,
            type: Boolean,
        },
    },
};
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
