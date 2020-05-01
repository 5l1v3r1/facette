<template>
    <div class="v-divider" :class="{horizontal: !vertical, vertical}"></div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class DividerComponent extends Vue {
    @Prop({default: false, type: Boolean})
    public vertical!: boolean;
}
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
