<template>
    <div class="v-columns" :style="{columnCount: count}">
        <slot></slot>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class ColumnsComponent extends Vue {
    @Prop({required: true, type: Number})
    public count!: number;
}
</script>
