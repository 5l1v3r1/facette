<template>
    <div class="v-spinner">
        <svg :width="size" :height="size" xmlns="http://www.w3.org/2000/svg">
            <g fill="none" fill-rule="evenodd" :stroke-width="strokeWidth">
                <circle stroke="var(--spinner-background)" :cx="half" :cy="half" :r="innerHalf" />
                <circle
                    stroke="var(--spinner-color)"
                    stroke-linecap="round"
                    stroke-dasharray="32.9867228627,131.9468914508"
                    :cx="half"
                    :cy="half"
                    :r="innerHalf"
                />
            </g>
        </svg>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class SpinnerComponent extends Vue {
    @Prop({default: 48, type: Number})
    public size!: number;

    @Prop({default: 3, type: Number})
    public strokeWidth!: number;

    public get half(): number {
        return this.size / 2;
    }

    public get innerHalf(): number {
        return (this.size - this.strokeWidth) / 2;
    }
}
</script>

<style lang="scss" scoped>
@keyframes rotate {
    from {
        transform: rotate(0deg);
    }
    to {
        transform: rotate(359deg);
    }
}

.v-spinner {
    align-items: center;
    display: inline-flex;
    justify-content: center;

    svg {
        animation: rotate 0.65s infinite linear;
    }
}
</style>
