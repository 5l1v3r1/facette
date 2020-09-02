<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

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
