<template>
    <div class="v-color" :class="{focus: focused}" @click="open">
        <div class="v-color-placeholder">
            <v-icon class="color" icon="eye-dropper" v-if="!value"></v-icon>

            <div class="color" :style="{backgroundColor: color}" v-else></div>

            <input
                ref="input"
                :placeholder="i18n.t('labels.default')"
                :value="value"
                @focus="onFocus"
                @focusout="onFocus"
                @input="onInput"
                @keydown="onKeydown"
            />
        </div>

        <v-dropdown @click.stop="onClick" @mousedown.prevent="onMouse" v-show="opened">
            <div
                class="v-color-saturation"
                ref="saturation"
                data-v-color-key="saturation"
                :style="{backgroundColor: hueColor}"
            >
                <div
                    class="cursor"
                    :style="{left: `${saturationCursor.x}px`, top: `${saturationCursor.y}px`}"
                    v-show="value"
                ></div>
            </div>

            <div class="v-color-hue" ref="hue" data-v-color-key="hue">
                <div class="cursor" :style="{top: `${hueCursor.y}px`}"></div>
            </div>
        </v-dropdown>
    </div>
</template>

<script lang="ts">
import {hex2Hsv, hsv2Hex} from "colorsys";
import {SetupContext, computed, nextTick, ref, watch} from "vue";
import {useI18n} from "vue-i18n";

interface Data {
    h: number;
    s: number;
    v: number;
}

const defaultData: Data = {
    h: 0,
    s: 100,
    v: 100,
};

function normalizeValue(input: string): string {
    let out = input;
    if (out.length === 3 || out.length === 4) {
        out = out[0] + out[0] + out[1] + out[1] + out[2] + out[2];
    }
    return !out.startsWith("#") ? `#${out}` : out;
}

export default {
    props: {
        value: {
            required: true,
            type: String,
        },
    },
    setup(props: Record<string, any>, ctx: SetupContext): Record<string, unknown> {
        const i18n = useI18n();

        let dragging: "hue" | "saturation" | null = null;

        const data = ref<Data>(Object.assign({}, defaultData));
        const focused = ref(false);
        const hue = ref<HTMLElement | null>(null);
        const input = ref<HTMLInputElement | null>(null);
        const opened = ref(false);
        const rect = ref<Record<"hue" | "saturation", DOMRect> | null>(null);
        const saturation = ref<HTMLElement | null>(null);

        const color = computed(() => hsv2Hex(data.value));

        const hueColor = computed(() => hsv2Hex(data.value.h, 100, 100));

        const hueCursor = computed(() => {
            if (!rect.value) {
                return {y: 0};
            }

            return {
                y: rect.value.hue.height - (data.value.h * rect.value.hue.height) / 360,
            };
        });

        const saturationCursor = computed(() => {
            if (!rect.value) {
                return {x: 0, y: 0};
            }

            return {
                x: (data.value.s * rect.value.saturation.width) / 100,
                y: rect.value.saturation.height - (data.value.v * rect.value.saturation.height) / 100,
            };
        });

        const onFocus = (ev: FocusEvent): void => {
            focused.value = ev.type === "focus";
            if (!focused.value) {
                opened.value = false;
                ctx.emit("update:value", props.value ? normalizeValue(props.value) : "");
            }
        };

        const onClick = (ev: MouseEvent): void => {
            const key = (ev.target as HTMLElement).dataset.vColorKey;
            if (key !== undefined) {
                select(ev, key);
            }
        };

        const onKeydown = (ev: KeyboardEvent): void => {
            switch (ev.code) {
                case "Enter":
                    if (!opened.value) {
                        open();
                        ev.stopPropagation();
                    }
                    break;

                case "Escape":
                    if (opened.value) {
                        opened.value = false;
                        ev.stopPropagation();
                    }
                    break;
            }
        };

        const onInput = (): void => {
            ctx.emit("update:value", input.value?.value);
        };

        const onMouse = (ev: MouseEvent): void => {
            switch (ev.type) {
                case "mousedown": {
                    const key = (ev.target as HTMLElement).dataset.vColorKey;
                    if (key !== "hue" && key !== "saturation") {
                        return;
                    }

                    document.addEventListener("mousemove", onMouse);
                    document.addEventListener("mouseup", onMouse);
                    dragging = key;

                    break;
                }

                case "mouseup": {
                    document.removeEventListener("mousemove", onMouse);
                    document.removeEventListener("mouseup", onMouse);
                    dragging = null;

                    break;
                }

                case "mousemove": {
                    if (dragging !== null) {
                        select(ev, dragging);
                    }

                    break;
                }
            }
        };

        const open = (): void => {
            opened.value = true;

            if (rect.value === null) {
                nextTick(() => {
                    rect.value = {
                        hue: hue.value?.getBoundingClientRect() as DOMRect,
                        saturation: saturation.value?.getBoundingClientRect() as DOMRect,
                    };
                });
            }
        };

        const select = (ev: MouseEvent, key: string) => {
            if (key === "hue") {
                const {y, height} = rect.value?.hue as DOMRect;

                data.value.h = Math.round(
                    360 - ((ev.pageY < y ? 0 : ev.pageY > y + height ? height : ev.pageY - y) * 360) / height,
                );
            } else if (key === "saturation") {
                const {x, y, width, height} = rect.value?.saturation as DOMRect;

                Object.assign(data.value, {
                    s: ((ev.pageX < x ? 0 : ev.pageX > x + width ? width : ev.pageX - x) * 100) / width,
                    v: 100 - ((ev.pageY < y ? 0 : ev.pageY > y + height ? height : ev.pageY - y) * 100) / height,
                });
            }

            ctx.emit("update:value", hsv2Hex(data.value));
        };

        watch(
            () => props.value,
            (to: string) => {
                if (to) {
                    try {
                        data.value = hex2Hsv(normalizeValue(to));
                    } catch (err) {}
                }
            },
            {deep: true, immediate: true},
        );

        return {
            color,
            data,
            focused,
            hue,
            hueColor,
            hueCursor,
            i18n,
            input,
            onClick,
            onFocus,
            onKeydown,
            onInput,
            onMouse,
            open,
            opened,
            saturation,
            saturationCursor,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../input/mixins";

.v-color {
    @include input;

    .v-color-placeholder {
        align-items: center;
        display: flex;
        flex-grow: 1;
        height: calc(var(--input-height) - 1rem);
        justify-content: center;

        .v-icon {
            color: var(--light-gray);
            border: 1px solid var(--dark-gray);
            font-size: 0.65rem;
        }

        .color {
            border-radius: 0.2rem;
            height: 100%;
            width: 1.5rem;
            margin-right: 0.75rem;
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
        flex-direction: row;
        height: 10rem;
        margin-top: 1px;
        padding: 0.65rem;
        width: 16rem;

        .cursor {
            border: 0.1rem solid white;
            filter: drop-shadow(0 0 0.1rem var(--dropdown-shadow));
            pointer-events: none;
            position: absolute;
            z-index: 1;
        }

        .v-color-saturation {
            border-radius: 0.15rem;
            flex-grow: 1;
            position: relative;

            &::before,
            &::after {
                border-radius: inherit;
                content: "";
                position: absolute;
                width: 100%;
                height: 100%;
            }

            &::before {
                background-image: linear-gradient(to right, white, transparent);
            }

            &::after {
                background-image: linear-gradient(to top, black, transparent);
            }

            .cursor {
                border: 0.1rem solid white;
                border-radius: 100%;
                height: 0.65rem;
                transform: translate(-0.325rem, -0.325rem);
                width: 0.65rem;
            }
        }

        .v-color-hue {
            background-image: linear-gradient(to top, red, yellow, green, cyan, blue, magenta, red);
            border-radius: 0.15rem;
            margin-left: 0.65rem;
            position: relative;
            width: 0.75rem;

            .cursor {
                border-radius: 0.1rem;
                height: 0.35rem;
                transform: translate(-0.1rem, -0.125rem);
                width: 0.95rem;
            }
        }
    }
}
</style>
