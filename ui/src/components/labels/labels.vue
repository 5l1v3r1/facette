<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-flex class="v-labels">
        <div
            class="v-labels-entry"
            :class="{selectable}"
            :key="key"
            @click="$emit('select', key, value)"
            v-for="(value, key) in labels"
        >
            <span class="v-labels-key">{{ key }}</span>

            <span class="v-labels-value empty" v-if="value === ''">{{ i18n.t("labels.empty") }}</span>

            <span class="v-labels-value" v-else>{{ value }}</span>
        </div>
    </v-flex>
</template>

<script lang="ts">
import {useI18n} from "vue-i18n";

export default {
    props: {
        labels: {
            required: true,
            type: Object as () => Record<string, string>,
        },
        selectable: {
            default: false,
            type: Boolean,
        },
    },
    setup(): Record<string, unknown> {
        const i18n = useI18n();

        return {
            i18n,
        };
    },
};
</script>

<style lang="scss" scoped>
.v-labels {
    .v-labels-entry {
        background-color: var(--button-focus-background);
        border-radius: 0.2rem;
        line-height: 1.25rem;
        padding: 0 0.35rem;

        &.selectable {
            cursor: pointer;
        }

        .v-labels-key {
            margin-right: 0.25rem;
            opacity: 0.65;

            &::after {
                content: ":";
            }
        }

        .v-labels-value {
            color: var(--button-focus-color);

            &.empty {
                opacity: 0.5;
                text-transform: lowercase;
            }
        }

        + .v-labels-entry {
            margin-left: 0.35rem;
        }
    }
}
</style>
