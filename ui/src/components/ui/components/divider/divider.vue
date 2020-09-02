<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div class="v-divider" :class="{horizontal: !vertical, vertical}"></div>
</template>

<script lang="ts">
export default {
    props: {
        vertical: {
            default: false,
            type: Boolean,
        },
    },
};
</script>

<style lang="scss" scoped>
.v-divider {
    background-color: var(--divider-background);

    &.horizontal {
        height: 1px;
        width: 100%;
    }

    &.vertical {
        height: 100%;
        width: 1px;
    }
}
</style>
