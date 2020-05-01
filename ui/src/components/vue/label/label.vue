<template>
    <div class="v-label">
        <slot></slot>
    </div>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator";

@Component
export default class LabelComponent extends Vue {}
</script>

<style lang="scss" scoped>
.v-label {
    line-height: var(--label-height);
    min-height: var(--label-height);
}
</style>
