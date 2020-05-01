<template>
    <div
        class="v-tab"
        tabindex="0"
        :aria-disabled="disabled"
        :class="{active}"
        @click="$parent.select(value)"
        @keydown.enter.stop="$parent.select(value)"
    >
        <slot></slot>
        <div class="v-tab-badge" v-if="badge">{{ badge }}</div>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

import TabBarComponent from "./tab-bar.vue";

@Component
export default class TabComponent extends Vue {
    @Prop({default: null, type: [Number, String]})
    public badge!: number | string | null;

    @Prop({default: false, type: Boolean})
    public disabled!: boolean;

    @Prop({default: "", required: true})
    public value!: unknown;

    public get active(): boolean {
        return (this.$parent as TabBarComponent).value === this.value;
    }
}
</script>

<style lang="scss" scoped>
.v-tab {
    align-items: center;
    border-bottom: 0.25rem solid transparent;
    cursor: pointer;
    display: flex;
    padding: 0 1rem;
    transform: translateY(1px);

    &[aria-disabled] {
        color: var(--light-gray);
        pointer-events: none;

        .v-tab-badge {
            background-color: var(--light-gray);
        }
    }

    &:focus,
    &.active:focus {
        border-color: var(--color);
        outline: none;

        .v-tab-badge {
            background-color: var(--color);
        }
    }

    &.active {
        border-color: var(--tab-active-border);
        outline: none;
    }

    .v-tab-badge {
        background-color: var(--color);
        border-radius: 0.2rem;
        color: var(--background);
        font-size: 0.75rem;
        height: 0.9rem;
        line-height: 0.9rem;
        margin-left: 0.75rem;
        padding: 0 0.25rem;
    }
}
</style>
