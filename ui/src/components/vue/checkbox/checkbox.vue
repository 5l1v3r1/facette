<template>
    <div
        class="v-checkbox"
        tabindex="0"
        :aria-disabled="disabled"
        :aria-readonly="readonly"
        :class="{check: !toggle, checked: value, toggle}"
        @click="onChange"
        @keydown.enter="onChange"
    >
        <v-icon :icon="icon"></v-icon>
        <div class="v-checkbox-label" v-if="label">
            <slot></slot>
        </div>
        <input ref="input" type="checkbox" :checked="value" :readonly="readonly" @change="update" />
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class CheckboxComponent extends Vue {
    @Prop({default: false, type: Boolean})
    public autofocus!: boolean;

    @Prop({default: false, type: Boolean})
    public disabled!: boolean;

    @Prop({default: false, type: Boolean})
    public readonly!: boolean;

    @Prop({default: false, type: Boolean})
    public toggle!: boolean;

    @Prop({default: false, required: true, type: Boolean})
    public value!: boolean;

    private input!: HTMLInputElement;

    public label = false;

    public mounted(): void {
        this.checkSlots();
        this.input = this.$refs.input as HTMLInputElement;
    }

    public get icon(): string {
        if (this.toggle) {
            return this.value ? "toggle-on" : "toggle-off";
        } else {
            return this.value ? "check" : "";
        }
    }

    public onChange(): void {
        this.input.checked = !this.input.checked;
        this.update();
    }

    public updated(): void {
        this.checkSlots();
    }

    public update(): void {
        this.$emit("input", this.input.checked);
    }

    private checkSlots(): void {
        this.label = Boolean(this.$slots?.default);
    }
}
</script>

<style lang="scss" scoped>
.v-checkbox {
    align-items: center;
    cursor: pointer;
    display: inline-flex;
    height: var(--input-height);
    line-height: var(--input-height);
    vertical-align: middle;

    &[aria-disabled] {
        color: var(--light-gray);
        pointer-events: none;

        &.check.checked .v-icon {
            background-color: var(--gray);
        }
    }

    &:focus {
        outline: none;
    }

    &.check {
        &:focus .v-icon {
            border-color: var(--color);
        }

        &.checked:focus .v-icon {
            background-color: var(--checkbox-focus);
            color: var(--white);
        }

        .v-icon {
            background-color: var(--checkbox-background);
            border-radius: 0.2rem;
            border: 1.5px solid var(--checkbox-border);
            color: var(--checkbox-color);
            height: 1em;
            margin: 0.125em;
            pointer-events: none;
            width: 1em;

            ::v-deep .fa {
                font-size: 0.75rem;
            }
        }

        &.checked .v-icon {
            background-color: var(--color);
            border: none;
            color: var(--background);
        }
    }

    &.toggle .v-icon {
        font-size: 1.1rem;
    }

    .v-checkbox-label {
        flex-grow: 1;
        margin-left: 0.5rem;
        pointer-events: none;
    }

    input {
        display: none;
    }
}
</style>
