<template>
    <div class="v-radio" tabindex="0" :class="{active}" @click="update" @keypress.enter="update">
        <v-icon :icon="active ? 'dot-circle' : 'circle'"></v-icon>
        <input ref="input" type="radio" :id="id" :name="$parent.name" :value="value" />
        <label :for="id" v-if="label">
            <slot></slot>
        </label>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

import RadioGroupComponent from "./radio-group.vue";

@Component
export default class RadioComponent extends Vue {
    @Prop({required: true, type: String})
    public value!: string;

    public label = false;

    public mounted(): void {
        this.checkSlots();
    }

    public updated(): void {
        this.checkSlots();
    }

    public get active(): boolean {
        return this.value === (this.$parent as RadioGroupComponent).value;
    }

    public get id(): string {
        return `${(this.$parent as RadioGroupComponent).name}-${this.value}`;
    }

    public update(): void {
        (this.$parent as RadioGroupComponent).update(this.value);
    }

    private checkSlots(): void {
        this.label = Boolean(this.$slots && this.$slots.default);
    }
}
</script>

<style lang="scss" scoped>
.v-radio {
    align-items: center;
    display: flex;
    line-height: var(--radio-height);

    &:focus {
        outline: none;
    }

    > .v-icon {
        color: var(--radio-color);
        font-size: 1rem;
        margin-right: 0.5rem;
    }

    &:focus > .v-icon,
    &.active > .v-icon {
        color: var(--radio-focus-color);
    }

    input {
        display: none;
    }

    label {
        align-items: center;
        display: flex;
        flex-grow: 1;
    }
}
</style>
