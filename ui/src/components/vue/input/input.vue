<template>
    <div
        class="v-input"
        :aria-disabled="disabled"
        :aria-invalid="validity || !checkValidity()"
        :aria-readonly="readonly"
        :class="{focus: focused, [type]: true}"
        @click="focus"
        @shortcut="shortcut && focus(true)"
        v-shortcut="shortcut"
    >
        <v-icon :icon="icon" @mousedown.native.prevent v-if="icon !== null"></v-icon>

        <v-label v-if="label">{{ label }}</v-label>

        <input
            ref="input"
            :disabled="disabled"
            :pattern="pattern"
            :placeholder="placeholderLabel"
            :readonly="readonly"
            :required="required"
            :type="type"
            :value="value"
            @focus="onFocus"
            @focusout="onFocus"
            @input="update(false)"
            @keydown="onKeydown"
        />

        <v-icon
            class="validity"
            icon="exclamation-triangle"
            @mousedown.native.prevent
            v-tooltip="validity"
            v-if="validity"
        ></v-icon>

        <v-icon class="help" icon="question-circle" @mousedown.native.prevent v-tooltip="help" v-if="help"></v-icon>
    </div>
</template>

<script lang="ts">
import debounce from "lodash/debounce";
import {Component, Prop, Vue} from "vue-property-decorator";

import {shortcutLabel} from "../shortcut";

@Component
export default class InputComponent extends Vue {
    @Prop({default: false, type: Boolean})
    public autofocus!: boolean;

    @Prop({default: null, type: Function})
    public customValidity!: (value: string) => Promise<string>;

    @Prop({default: 0, type: Number})
    public delay!: number;

    @Prop({default: false, type: Boolean})
    public disabled!: boolean;

    @Prop({default: null, type: String})
    public help!: string;

    @Prop({default: null, type: String})
    public icon!: string;

    @Prop({default: null, type: String})
    public label!: string | null;

    @Prop({default: null, type: String})
    public pattern!: string | null;

    @Prop({default: null, type: String})
    public placeholder!: string | null;

    @Prop({default: false, type: Boolean})
    public readonly!: boolean;

    @Prop({default: false, type: Boolean})
    public required!: boolean;

    @Prop({default: null, type: [Array, String]})
    public shortcut!: string | Array<string> | null;

    @Prop({default: "text", type: String})
    public type!: string;

    @Prop({default: "", required: true, type: [Number, String]})
    public value!: number | string;

    public focused = false;

    public validity = "";

    private input!: HTMLInputElement;

    private pristine = true;

    private updateDebounce: (() => void) | null = null;

    public mounted(): void {
        this.input = this.$refs.input as HTMLInputElement;

        if (this.delay > 0) {
            this.updateDebounce = debounce(() => {
                this.update(true);
            }, this.delay);
        }

        if (this.autofocus) {
            this.$nextTick(() => {
                this.focus(true);
            });
        }
    }

    public checkValidity(): boolean {
        if (!this.input) {
            return true;
        }

        const valid = this.input.checkValidity();
        return valid || (this.pristine && !valid);
    }

    public clear(): void {
        if (this.input) {
            this.input.setCustomValidity("");
        }

        this.$emit("input", "");
        this.$emit("clear");
    }

    private async emit(value: string): Promise<void> {
        if (this.customValidity) {
            this.validity = await this.customValidity(value);
            this.input.setCustomValidity(this.validity);
        }

        if (this.pristine) {
            this.pristine = false;
        }

        this.$emit("input", value);
    }

    public focus(select = false): void {
        if (!this.focused) {
            const input = this.$refs.input as HTMLInputElement;
            if (input) {
                input[select ? "select" : "focus"]();
            }
        }
    }

    public onFocus(e: FocusEvent): void {
        this.focused = e.type === "focus";
        this.$emit(e.type, e);

        if (!this.focused && this.pristine) {
            this.pristine = false;
        }
    }

    public onKeydown(e: KeyboardEvent): void {
        switch (e.code) {
            case "Escape": {
                if (this.type === "search") {
                    if (this.value) {
                        this.clear();
                    } else {
                        this.input.blur();
                    }
                }
            }
        }
    }

    public get placeholderLabel(): string | null {
        let placeholder = "";

        if (this.placeholder) {
            placeholder += this.placeholder;
        }

        if (this.shortcut) {
            if (this.placeholder) {
                placeholder += " ";
            }
            placeholder += `(${shortcutLabel(this.shortcut)})`;
        }

        return placeholder ?? null;
    }

    public update(apply = false): void {
        if (apply || this.delay === 0) {
            this.emit(this.input.value);
        } else if (this.updateDebounce !== null) {
            this.updateDebounce();
        }
    }
}
</script>

<style lang="scss" scoped>
@import "mixins";

.v-input {
    @include input;

    input {
        background: none;
        border: none;
        color: inherit;
        flex-grow: 1;
        font: inherit;
        min-width: 0;
        padding: 0;

        &[type="number"] {
            -moz-appearance: textfield;

            &::-webkit-inner-spin-button,
            &::-webkit-outer-spin-button {
                -webkit-appearance: none;
            }
        }

        &:focus,
        &:invalid {
            box-shadow: none;
            outline: none;
        }
    }
}
</style>
