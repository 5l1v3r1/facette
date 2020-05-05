<template>
    <div class="v-radio-group">
        <slot></slot>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class RadioGroupComponent extends Vue {
    @Prop({required: true, type: String})
    public name!: string;

    @Prop({required: true, type: String})
    public value!: string;

    public focus(): void {
        const el: HTMLElement | null = this.$el.querySelector(".v-radio");
        if (el) {
            el.focus();
        }
    }

    public update(value: string): void {
        this.$emit("input", value);
    }
}
</script>
