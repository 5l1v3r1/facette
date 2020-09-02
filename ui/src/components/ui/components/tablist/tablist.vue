<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div class="v-tablist" role="tablist">
        <v-button
            role="tab"
            :class="{active: tab.value === value}"
            :icon="tab.icon"
            :key="index"
            @click="select(tab.value)"
            v-for="(tab, index) in tabs"
        >
            {{ tab.label }}
        </v-button>
    </div>
</template>

<script lang="ts">
import {SetupContext} from "vue";

import {Tab} from "types/ui";

export default {
    props: {
        tabs: {
            required: true,
            type: Array as () => Array<Tab>,
        },
        value: {
            required: true,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const select = (value: unknown): void => {
            ctx.emit("update:value", value);
        };

        return {
            select,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-tablist {
    align-items: center;
    border-bottom: 1px solid var(--tablist-border);
    display: flex;
    height: var(--tablist-height);
    line-height: var(--tablist-height);
    margin-bottom: 2rem;

    ::v-deep(.v-button) {
        height: var(--tablist-height);
        line-height: var(--tablist-height);
        margin: 0;

        .v-button-content {
            border-bottom: 0.25rem solid transparent;
            border-radius: 0;
            padding: 0 1rem;

            &:active,
            &:focus {
                background-color: unset;
                border-color: var(--color) !important;
                outline: none;
            }

            &:hover {
                background-color: unset;
            }
        }

        &.active .v-button-content {
            background-color: unset;
            border-color: var(--accent);
        }
    }
}
</style>
