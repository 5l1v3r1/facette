<template>
    <div class="v-flex" :class="{reverse, [direction]: true}">
        <slot></slot>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class FlexComponent extends Vue {
    @Prop({default: "row", type: String})
    public direction!: string;

    @Prop({default: false, type: Boolean})
    public reverse!: boolean;
}
</script>

<style lang="scss" scoped>
.v-flex {
    align-items: center;
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
