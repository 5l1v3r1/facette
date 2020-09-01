<template>
    <div class="v-flex" :class="{[direction]: true, reverse}">
        <slot></slot>
    </div>
</template>

<script lang="ts">
export default {
    props: {
        direction: {
            default: "row",
            type: String,
        },
        reverse: {
            default: false,
            type: Boolean,
        },
    },
};
</script>

<style lang="scss" scoped>
.v-flex {
    display: flex;

    &.column {
        flex-direction: column;

        &.reverse {
            flex-direction: column-reverse;
        }
    }

    &.row {
        flex-direction: row;

        &.reverse {
            flex-direction: row-reverse;
        }
    }
}
</style>
