<template>
    <div
        class="v-select"
        role="listbox"
        :aria-invalid="invalid"
        :class="{focus: focused}"
        :tabindex="opened ? -1 : 0"
        @click="open"
        @focus="onFocus"
        @focusout="onFocus"
        @keydown.enter="open"
    >
        <v-icon :icon="iconValue" @mousedown.native.prevent v-if="iconValue"></v-icon>

        <v-label v-if="label">{{ label }}</v-label>

        <div
            class="v-select-value"
            :class="{placeholder: placeholderLabel === placeholder}"
            v-show="!search || !opened"
        >
            <div class="label">{{ placeholderLabel }}</div>

            <v-icon class="help" icon="question-circle" @mousedown.native.prevent v-tooltip="help" v-if="help"></v-icon>

            <v-icon icon="chevron-down"></v-icon>
        </div>

        <v-input
            ref="input"
            :delay="200"
            :placeholder="placeholder"
            @focus="onFocus"
            @focusout="onFocus"
            @keydown.native.stop="onKeydown"
            v-model="filter"
            v-if="search && opened"
        ></v-input>

        <select :required="required" :value="value">
            <option :key="index" :value="option.value" v-for="(option, index) in options">{{ option.label }}</option>
        </select>

        <v-dropdown v-show="opened">
            <slot name="dropdown-placeholder" v-if="currentOptions.length === 0"></slot>

            <v-button
                role="option"
                :class="{active: index === active}"
                :icon="option.icon"
                :key="index"
                @click.stop="select(index)"
                @mousedown.native.prevent
                @mouseenter.native="setActive(index)"
                v-for="(option, index) in currentOptions"
                >{{ option.label }}</v-button
            >
        </v-dropdown>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue, Watch} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

@Component
export default class SelectComponent extends Vue {
    @Prop({default: false, type: Boolean})
    public autofocus!: boolean;

    @Prop({default: null, type: String})
    public help!: string;

    @Prop({default: null, type: String})
    public icon!: string;

    @Prop({default: null, type: String})
    public label!: string | null;

    @Prop({required: true, type: Array})
    public options!: Array<SelectOption>;

    @Prop({default: null, type: String})
    public placeholder!: string | null;

    @Prop({default: false, type: Boolean})
    public required!: boolean;

    @Prop({default: true, type: Boolean})
    public search!: boolean;

    @Prop({required: true})
    public value!: unknown;

    public active = 0;

    public filter = "";

    public focused = false;

    public invalid = false;

    public opened = false;

    private filteredOptions: Array<SelectOption> = [];

    private selected = -1;

    public mounted(): void {
        if (this.autofocus) {
            this.$nextTick(() => {
                this.focus();
            });
        }
    }

    public get currentOptions(): Array<SelectOption> {
        return (this.filter ? this.filteredOptions : this.options) ?? [];
    }

    public focus(): void {
        (this.$el as HTMLElement).focus();
    }

    public get iconValue(): string | null {
        return (
            (this.value !== null && this.selected !== -1 && this.options && this.options[this.selected]
                ? this.options[this.selected].icon
                : this.icon) ?? null
        );
    }

    @Watch("value", {immediate: true})
    public onChange(to: unknown): void {
        for (let i = 0; i < this.options.length; i++) {
            if (this.options[i].value === to) {
                this.selected = i;
                return;
            }
        }

        this.selected = -1;
    }

    @Watch("filter")
    public onFilter(to: string): void {
        if (to) {
            to = to.toLowerCase();
        }

        this.filteredOptions =
            to && this.options
                ? this.options.filter((option: SelectOption) => option.label.toLowerCase().indexOf(to) !== -1)
                : [];

        if (this.active > this.filteredOptions.length - 1) {
            this.active = 0;
        }
    }

    public onFocus(e: FocusEvent): void {
        // Skip focus/focusout when switching from main element to input
        if (
            (e.type === "focusout" &&
                e.relatedTarget &&
                (e.relatedTarget as Element).closest(".v-select") === this.$el) ||
            (e.type === "focus" && e.relatedTarget === this.$el)
        ) {
            return;
        }

        this.focused = e.type === "focus";

        if (!this.focused) {
            this.opened = false;
            if (this.required) {
                this.invalid = !this.value;
            }
        }
    }

    public onKeydown(e: KeyboardEvent): void {
        switch (e.code) {
            case "ArrowDown": {
                e.preventDefault();
                const options: Array<SelectOption> = this.currentOptions;
                if (this.active < Object.keys(options).length - 1) {
                    this.active++;
                }
                break;
            }

            case "ArrowUp": {
                e.preventDefault();
                if (this.active > 0) {
                    this.active--;
                }
                break;
            }

            case "Enter": {
                this.select(this.active);
                this.focus();
                break;
            }

            case "Escape": {
                this.focus();
                break;
            }
        }
    }

    public open(): void {
        const active = this.options
            ? this.options.findIndex((option: SelectOption) => option.value === this.value)
            : -1;

        this.active = active !== -1 ? active : 0;
        this.opened = true;

        if (this.search) {
            this.filter = "";
            this.$nextTick(() => (this.$refs.input as HTMLInputElement).focus());
        }
    }

    public get placeholderLabel(): string {
        return (
            (this.value !== null && this.selected !== -1 && this.options && this.options[this.selected]
                ? this.options[this.selected].label
                : this.placeholder) ?? ""
        );
    }

    public select(index: number): void {
        const value = this.currentOptions[index] ? this.currentOptions[index].value : undefined;

        this.selected = index;
        this.opened = false;
        this.$emit("input", value);
        if (this.required) {
            this.invalid = !value;
        }
        this.$nextTick(() => {
            this.focus();
        });
    }

    public setActive(index: number): void {
        this.active = index;
    }
}
</script>

<style lang="scss" scoped>
@import "../input/mixins";

.v-select {
    @include input;

    &:focus {
        outline: none;
    }

    select {
        display: none;
    }

    .v-select-value {
        align-items: center;
        cursor: pointer;
        display: flex;
        flex-grow: 1;

        .label {
            flex-grow: 1;
        }

        &.placeholder .label {
            opacity: var(--placeholder-opacity);
        }
    }

    &[aria-invalid] {
        .v-select-value.placeholder .label {
            opacity: 0.5;
        }
    }

    .v-input {
        background-color: transparent;
        border: none;
        color: inherit;
        flex-grow: 1;
        min-width: auto;
        padding: 0;

        ::v-deep input {
            width: 0;
        }
    }

    .v-dropdown {
        left: 0;
        max-height: 35vh;
        overflow-y: auto;
        right: 0;
        top: calc(100% + 1px);
    }
}
</style>
