<template>
    <div class="v-color" :class="{focus: focused}" @click="open">
        <div class="v-color-placeholder">
            <v-icon class="color" icon="eye-dropper" v-if="!value"></v-icon>

            <div class="color" :style="{backgroundColor: value}" v-else></div>

            <input
                ref="input"
                :placeholder="$t('labels.auto')"
                :value="value"
                @focus="onFocus"
                @focusout="onFocus"
                @input="onInput"
                @keydown="onKeydown"
            />
        </div>

        <v-dropdown @mousedown.native.prevent v-show="opened">
            <v-chrome disable-alpha disable-fields :value="value" @input="onPickerInput"></v-chrome>
        </v-dropdown>
    </div>
</template>

<script lang="ts">
import {Chrome} from "vue-color";
import {Component, Prop, Vue} from "vue-property-decorator";

interface Color {
    hex?: string;
}

@Component({
    components: {
        "v-chrome": Chrome,
    },
})
export default class ColorComponent extends Vue {
    @Prop({required: true, type: String})
    public value!: string;

    public focused = false;

    public opened = false;

    public onFocus(e: Event): void {
        this.focused = e.type === "focus";

        if (!this.focused && this.opened) {
            this.opened = false;
        }
    }

    public onInput(e: KeyboardEvent): void {
        this.$emit("input", (e.target as HTMLInputElement).value);
    }

    public onKeydown(e: KeyboardEvent): void {
        switch (e.code) {
            case "Enter":
                if (!this.opened) {
                    this.open();
                    e.stopPropagation();
                }
                break;

            case "Escape":
                if (this.opened) {
                    this.opened = false;
                    e.stopPropagation();
                }
                break;
        }
    }

    public onPickerInput(to: Color): void {
        if (to.hex) {
            this.$emit("input", to.hex.toLowerCase());
        }
    }

    public open(): void {
        this.opened = true;
    }

    // public reset(e: KeyboardEvent): void {
    //     if (this.opened) {
    //         this.$emit("input", "");
    //         this.opened = false;
    //         e.preventDefault();
    //         e.stopPropagation();
    //     }
    // }
}
</script>

<style lang="scss" scoped>
@import "~@/src/components/vue/input/mixins";

.v-color {
    @include input;

    min-width: auto;

    .v-color-placeholder {
        align-items: center;
        display: flex;
        flex-grow: 1;
        height: calc(var(--input-height) - 1rem);
        justify-content: center;

        .v-icon {
            color: var(--light-gray);
            border: 1px dotted var(--dark-gray);
            font-size: 0.65rem;
        }

        .color {
            border-radius: 0.1rem;
            height: 100%;
            width: 1.5rem;
            margin-right: 0.35rem;
        }

        input {
            background-color: transparent;
            border: none;
            color: inherit;
            flex-grow: 1;
            font: inherit;

            &:focus {
                outline: none;
            }
        }
    }

    .v-dropdown {
        margin-top: 1px;
        padding: 0.35rem;
    }

    .vc-chrome {
        background-color: transparent;
        box-shadow: none;

        ::v-deep {
            .vc-chrome-saturation-wrap {
                border-radius: 0;
            }

            .vc-chrome-body {
                background-color: inherit;
            }

            .vc-chrome-color-wrap {
                display: none;
            }
        }
    }
}
</style>
