<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <div
        class="v-select"
        ref="el"
        role="listbox"
        :aria-disabled="disabled || undefined"
        :aria-invalid="(invalid && !opened) || undefined"
        :aria-readonly="readonly || undefined"
        :class="{'focus': focused, 'open-top': openTop}"
        :tabindex="opened ? -1 : 0"
        @autofocus="focus(true)"
        @click="open"
        @focus="onFocus"
        @focusout="onFocus"
        @keydown.enter.stop="open"
    >
        <v-icon :icon="iconValue" @mousedown.prevent v-if="iconValue"></v-icon>

        <v-label v-if="label">{{ label }}</v-label>

        <div
            class="v-select-value"
            :class="{placeholder: placeholderLabel === placeholder}"
            v-show="!opened || !searchable"
        >
            <div class="label">{{ placeholderLabel }}</div>

            <v-icon class="help" icon="question-circle" @mousedown.prevent v-tooltip="help" v-if="help"></v-icon>

            <v-icon icon="angle-down"></v-icon>
        </div>

        <v-input
            ref="input"
            :delay="350"
            :placeholder="placeholder"
            @focus="onFocus"
            @focusout="onFocus"
            @keydown.stop="onKeydown"
            v-model:value="filter"
            v-if="opened && searchable"
        ></v-input>

        <select
            ref="select"
            :required="required"
            :value="value"
            @report-validity="onReportValidity"
            @reset-validity="onResetValidity"
        >
            <option :key="index" :value="option.value" v-for="(option, index) in options">
                {{ option.label }}
            </option>
        </select>

        <v-dropdown ref="dropdown" v-show="opened">
            <slot name="dropdown-placeholder" v-if="currentOptions.length === 0"></slot>

            <v-button
                role="option"
                :class="{active: index === active}"
                :icon="option.icon"
                :key="index"
                @click.stop="selectOption(index)"
                @mousedown.prevent
                @mouseenter="active = index"
                v-for="(option, index) in currentOptions"
            >
                {{ option.label }}
            </v-button>
        </v-dropdown>
    </div>
</template>

<script lang="ts">
import {ComponentPublicInstance, SetupContext, computed, nextTick, ref, watch} from "vue";

import {SelectOption} from "types/ui";

export default {
    props: {
        disabled: {
            default: false,
            type: Boolean,
        },
        help: {
            default: null,
            type: String,
        },
        icon: {
            default: null,
            type: String,
        },
        label: {
            default: null,
            type: String,
        },
        options: {
            required: true,
            type: Array as () => Array<SelectOption>,
        },
        placeholder: {
            default: null,
            type: String,
        },
        readonly: {
            default: false,
            type: Boolean,
        },
        required: {
            default: false,
            type: Boolean,
        },
        searchable: {
            default: true,
            type: Boolean,
        },
        value: {
            default: "",
            required: true,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const active = ref(0);
        const dropdown = ref<ComponentPublicInstance | null>(null);
        const el = ref<HTMLElement | null>(null);
        const filter = ref("");
        const focused = ref(false);
        const input = ref<HTMLElement | null>(null);
        const invalid = ref(false);
        const opened = ref(false);
        const openTop = ref(false);
        const select = ref<HTMLSelectElement | null>(null);
        const selected = ref(-1);

        const currentOptions = computed<Array<SelectOption>>(() =>
            filter.value
                ? props.options.filter((option: SelectOption) => option.label.toLowerCase().includes(filter.value))
                : props.options,
        );

        const iconValue = computed(
            (): string | null =>
                (props.value !== undefined && selected.value !== -1
                    ? props.options?.[selected.value].icon
                    : props.icon) ?? null,
        );

        const placeholderLabel = computed(
            (): string =>
                (props.value !== undefined && selected.value !== -1
                    ? props.options?.[selected.value].label
                    : props.placeholder) ?? "",
        );

        const focus = (next = false): void => {
            if (next) {
                nextTick(() => el.value?.focus());
            } else {
                el.value?.focus();
            }
        };

        const onFocus = (ev: FocusEvent): void => {
            // Skip focus/focusout when switching from main element to input
            if (
                (ev.type === "focusout" && (ev.relatedTarget as HTMLElement)?.closest(".v-select") === el.value) ||
                (ev.type === "focus" && ev.relatedTarget === el.value)
            ) {
                return;
            }

            focused.value = ev.type === "focus";

            if (!focused.value) {
                opened.value = false;
                onReportValidity();
            }
        };

        const onKeydown = (ev: KeyboardEvent): void => {
            switch (ev.code) {
                case "ArrowDown":
                    ev.preventDefault();
                    if (active.value < currentOptions.value.length - 1) {
                        active.value++;
                    }
                    break;

                case "ArrowUp":
                    ev.preventDefault();
                    if (active.value > 0) {
                        active.value--;
                    }
                    break;

                case "Enter":
                    selectOption(active.value);
                    focus();
                    break;

                case "Escape":
                    opened.value = false;
                    focus();
                    break;
            }
        };

        const onReportValidity = (): void => {
            invalid.value = !select.value?.validity.valid;
        };

        const onResetValidity = (): void => {
            invalid.value = false;
        };

        const open = (): void => {
            const idx = props.options?.findIndex((option: SelectOption) => option.value === props.value);
            active.value = idx !== -1 ? idx : 0;
            opened.value = true;

            nextTick(() => {
                const {top, height} = el.value?.getBoundingClientRect() as DOMRect;
                openTop.value = top + height + (dropdown.value?.$el as HTMLElement).clientHeight > window.innerHeight;
            });

            if (props.searchable) {
                filter.value = "";
                nextTick(() => input.value?.focus());
            }
        };

        const selectOption = (index: number): void => {
            const value = currentOptions.value[index]?.value;

            selected.value = index;
            opened.value = false;

            ctx.emit("update:value", value);

            nextTick(() => {
                onReportValidity();
                focus();
            });
        };

        watch(filter, () => {
            // Ensure active option index is always within current length
            if (active.value > currentOptions.value.length - 1) {
                active.value = 0;
            }
        });

        watch(
            () => props.value,
            to => (selected.value = props.options.findIndex((option: SelectOption) => option.value === to)),
            {immediate: true},
        );

        return {
            active,
            currentOptions,
            dropdown,
            el,
            filter,
            focus,
            focused,
            iconValue,
            input,
            invalid,
            onFocus,
            onKeydown,
            onReportValidity,
            onResetValidity,
            open,
            opened,
            placeholderLabel,
            select,
            selectOption,
            openTop,
        };
    },
};
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
            user-select: none;
        }

        &.placeholder .label {
            opacity: var(--placeholder-opacity);
        }
    }

    .v-input {
        background-color: transparent;
        border: none;
        color: inherit;
        flex-grow: 1;
        min-width: auto;
        padding: 0;
    }

    .v-dropdown {
        left: 0;
        max-height: 35vh;
        overflow-y: auto;
        right: 0;
        top: calc(100% + 1px);
    }

    &.open-top .v-dropdown {
        bottom: calc(100% + 1px);
        top: auto;
    }
}
</style>
