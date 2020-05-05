<template>
    <v-flex class="v-labels">
        <v-flex
            class="v-labels-item"
            :class="{handler: handler}"
            :key="key"
            @click.native="handler && handler(key, value)"
            v-for="(value, key) in filteredLabels"
        >
            <div class="key">{{ key }}</div>
            <div class="empty value" v-if="value === ''">{{ $t("labels.empty") }}</div>
            <div class="value" v-else>{{ value }}</div>
        </v-flex>
    </v-flex>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

import {NameLabel} from "@/src/helpers/labels";

export type LabelHandler = (key: string, value: string) => void;

@Component
export default class LabelsComponent extends Vue {
    @Prop({default: null, type: Function})
    public handler!: LabelHandler;

    @Prop({required: true, type: Object})
    public labels!: Labels;

    public get filteredLabels(): Labels {
        const ls = Object.assign({}, this.labels);
        delete ls[NameLabel];
        delete ls["_hash"];

        return ls;
    }
}
</script>

<style lang="scss" scoped>
.v-labels {
    font-size: 0.8rem;

    .v-labels-item {
        background-color: var(--button-focus-background);
        border-radius: 0.2rem;
        line-height: 1.25rem;
        padding: 0 0.35rem;

        &.handler {
            cursor: pointer;
        }

        .key {
            margin-right: 0.25rem;
            opacity: 0.65;

            &::after {
                content: ":";
            }
        }

        .value {
            color: var(--button-focus-color);

            &.empty {
                opacity: 0.5;
                text-transform: lowercase;
            }
        }

        + .v-labels-item {
            margin-left: 0.35rem;
        }
    }
}
</style>
