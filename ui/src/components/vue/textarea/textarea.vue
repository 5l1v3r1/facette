<template>
    <div
        class="v-textarea"
        :aria-disabled="disabled"
        :aria-readonly="readonly"
        :class="{focus: focused, help}"
        @click="focus"
    >
        <v-label v-if="label">{{ label }}</v-label>

        <textarea
            ref="textarea"
            :disabled="disabled"
            :placeholder="placeholderLabel"
            :readonly="readonly"
            :required="required"
            :value="value"
            @focus="onFocus"
            @focusout="onFocus"
            @input="update(false)"
        ></textarea>

        <v-icon class="help" icon="question-circle" @mousedown.native.prevent v-tooltip="help" v-if="help"></v-icon>
    </div>
</template>

<script lang="ts">
import debounce from "lodash/debounce";
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class TextAreaComponent extends Vue {
    @Prop({default: 0, type: Number})
    public delay!: number;

    @Prop({default: false, type: Boolean})
    public disabled!: boolean;

    @Prop({default: null, type: String})
    public help!: string;

    @Prop({default: null, type: String})
    public label!: string | null;

    @Prop({default: null, type: String})
    public placeholder!: string | null;

    @Prop({default: false, type: Boolean})
    public readonly!: boolean;

    @Prop({default: false, type: Boolean})
    public required!: boolean;

    @Prop({default: "", required: true, type: [Number, String]})
    public value!: number | string;

    public focused = false;

    private form: HTMLFormElement | null = null;

    private textarea!: HTMLTextAreaElement;

    private pristine = true;

    private updateDebounce: (() => void) | null = null;

    public mounted(): void {
        this.textarea = this.$refs.textarea as HTMLTextAreaElement;
        this.form = this.textarea.closest("form");

        if (this.delay > 0) {
            this.updateDebounce = debounce(() => {
                this.update(true);
            }, this.delay);
        }
    }

    private async emit(value: string): Promise<void> {
        if (this.pristine) {
            this.pristine = false;
        }

        this.$emit("input", value);

        if (this.form !== null) {
            this.$components.$emit("form-input", this.form);
        }
    }

    public focus(select = false): void {
        if (!this.focused) {
            const textarea = this.$refs.textarea as HTMLTextAreaElement;
            if (textarea) {
                textarea[select ? "select" : "focus"]();
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

    public get placeholderLabel(): string | null {
        let placeholder = "";
        if (this.placeholder) {
            placeholder += this.placeholder;
        }

        return placeholder ?? null;
    }

    public update(apply = false): void {
        if (apply || this.delay === 0) {
            this.emit(this.textarea.value);
        } else if (this.updateDebounce !== null) {
            this.updateDebounce();
        }
    }
}
</script>

<style lang="scss" scoped>
@import "../input/mixins";

.v-textarea {
    @include input;

    align-items: flex-start;
    flex-direction: column;
    height: auto;
    line-height: 1.5;
    padding: 0.35rem;

    textarea {
        background: none;
        border: none;
        color: inherit;
        flex-grow: 1;
        font: inherit;
        margin: 0;
        min-height: 3.5rem;
        padding: 0 0.25rem;
        resize: vertical;
        width: 100%;

        &:focus,
        &:invalid {
            box-shadow: none;
            outline: none;
        }
    }

    &.help textarea {
        padding-right: 1.5rem;
    }

    .v-label {
        line-height: 1.75rem;
        min-height: auto;
        padding: 0 0.25rem;
    }

    .v-icon {
        position: absolute;
        right: 0.25rem;
        top: 0.5rem;
    }
}
</style>
