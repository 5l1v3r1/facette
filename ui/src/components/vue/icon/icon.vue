<template>
    <div class="v-icon">
        <div class="fa-fw" :class="classes"></div>
        <div class="v-icon-badge" v-if="badge">{{ badge }}</div>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class IconComponent extends Vue {
    @Prop({default: null, type: [Number, String]})
    public badge!: number | string | null;

    @Prop({required: true, type: String})
    public icon!: string;

    @Prop({default: false, type: Boolean})
    public spin!: boolean;

    public get classes(): Array<string> | null {
        if (!this.icon) {
            return null;
        }

        let classes: Array<string>;
        const index = this.icon.indexOf("/");
        if (index !== -1) {
            classes = [this.icon.substr(0, index), `fa-${this.icon.substr(index + 1)}`];
        } else {
            classes = ["fa", `fa-${this.icon}`];
        }

        if (this.spin) {
            classes.push("fa-spin");
        }

        return classes;
    }
}
</script>

<style lang="scss" scoped>
$fa-font-path: "~@fortawesome/fontawesome-free/webfonts";
@import "~@fortawesome/fontawesome-free/scss/fontawesome";
@import "~@fortawesome/fontawesome-free/scss/regular";
@import "~@fortawesome/fontawesome-free/scss/solid";
@import "~@fortawesome/fontawesome-free/scss/brands";

.v-icon {
    align-items: center;
    display: flex;
    justify-content: center;
    position: relative;

    .v-icon-badge {
        background-color: var(--badge-background);
        border-radius: 1.1rem;
        color: var(--badge-color);
        left: 100%;
        line-height: 1.1rem;
        min-height: 1.1rem;
        min-width: 1.1rem;
        padding: 0 0.35rem;
        position: absolute;
        top: 100%;
        transform: scale(0.7) translate(-100%, -100%);
    }

    .fa-spin {
        animation: fa-spin 0.65s infinite linear;
    }
}
</style>
