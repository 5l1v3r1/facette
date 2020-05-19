<template>
    <div
        class="v-modal"
        tabindex="-1"
        @click="close(false)"
        @keydown.enter="onKeydown"
        @keydown.esc="close(false)"
        @keydown.tab="onKeydown"
        v-if="visible"
    >
        <div class="v-modal-content" @click.stop>
            <slot v-bind="{close}"></slot>
        </div>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

import {ModalHandler} from "@/types/components";

interface ComponentWithFn {
    onModalShow: (params: unknown) => void;
}

@Component
export default class ModalComponent extends Vue {
    @Prop({required: true, type: String})
    public name!: string;

    public params: unknown = null;

    public visible = false;

    private handler: ModalHandler | null = null;

    public mounted(): void {
        this.$components.$on("modal-open", this.open);
    }

    public beforeDestroy(): void {
        this.$components.$off("modal-open", this.open);
    }

    public close(value: unknown): void {
        if (this.handler) {
            this.handler(value);
        }

        Object.assign(this, {
            handler: null,
            params: null,
            visible: false,
        });
    }

    public focus(): void {
        (this.$el as HTMLElement).focus();
    }

    public onKeydown(e: KeyboardEvent): void {
        switch (e.code) {
            case "Enter": {
                if ((e.target as HTMLElement).closest(".v-form-bottom") === null) {
                    const primary: HTMLElement | null = this.$el.querySelector(".v-button.primary");
                    if (primary && !primary.hasAttribute("aria-disabled")) {
                        const content: HTMLElement | null = primary.querySelector(".v-button-content");
                        if (content) {
                            content.dispatchEvent(new Event("click"));
                        }
                    }
                }

                break;
            }

            case "Tab": {
                const focusable: NodeListOf<Element> = this.$el.querySelectorAll(
                    'input, [tabindex]:not([tabindex="-1"])',
                );
                let index: number | null = null;

                if (!e.shiftKey && e.target === focusable[focusable.length - 1]) {
                    index = 0;
                } else if (e.shiftKey && e.target === focusable[0]) {
                    index = focusable.length - 1;
                }

                if (index !== null) {
                    e.preventDefault();
                    this.$nextTick(() => {
                        (focusable[index as number] as HTMLElement).focus();
                    });
                }

                break;
            }
        }
    }

    public open(name: string, params: unknown = {}, handler: ModalHandler | null = null): void {
        if (name !== this.name) {
            return;
        }

        // TODO: hide previously existing overlays?

        Object.assign(this, {
            handler,
            params,
            visible: true,
        });

        const fn = (this.$parent as Vue & ComponentWithFn).onModalShow;
        if (fn) {
            fn(params);
        }

        this.$nextTick(() => {
            this.focus();
        });
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    align-items: center;
    background-color: var(--modal-background);
    bottom: 0;
    display: flex;
    justify-content: center;
    left: 0;
    position: fixed;
    right: 0;
    top: 0;
    z-index: 500;

    .v-modal-content {
        background-color: var(--background);
        border-radius: 0.2rem;
        box-shadow: 0 0.2rem 1rem var(--modal-content-shadow);
        min-width: 20vw;
        padding: 1.5rem;
    }

    .v-form ::v-deep .v-form-bottom {
        justify-content: flex-end;
    }
}
</style>
