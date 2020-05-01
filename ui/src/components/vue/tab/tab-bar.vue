<template>
    <div class="v-tab-bar">
        <slot></slot>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class TabBarComponent extends Vue {
    @Prop({default: "", required: true})
    public value!: unknown;

    public select(value: string): void {
        this.$emit("input", value);
    }
}
</script>

<style lang="scss" scoped>
.v-tab-bar {
    align-items: center;
    border-bottom: 1px solid var(--tab-bar-border);
    display: flex;
    line-height: var(--tab-bar-height);
    margin-bottom: 1rem;
}
</style>
