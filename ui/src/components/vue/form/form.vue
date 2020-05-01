<template>
    <form class="v-form" novalidate @input="onChange" @submit.prevent>
        <slot></slot>

        <div class="v-form-bottom" v-if="bottom">
            <slot name="bottom"></slot>
        </div>
    </form>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator";

@Component
export default class FormComponent extends Vue {
    public bottom = false;

    public mounted(): void {
        this.checkSlots();
        this.$nextTick(this.onChange);
    }

    public updated(): void {
        this.checkSlots();
    }

    public onChange(): void {
        this.$emit("validity", (this.$el as HTMLFormElement).checkValidity());
    }

    private checkSlots(): void {
        this.bottom = Boolean(this.$slots?.bottom);
    }
}
</script>

<style lang="scss" scoped>
.v-form {
    .v-form-bottom {
        align-items: center;
        display: flex;
        margin-top: 1.5rem;
    }
}
</style>
