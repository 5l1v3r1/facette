<template>
    <div class="v-dropdown" :class="{[anchor]: anchor}">
        <slot></slot>
    </div>
</template>

<script lang="ts">
export default {
    props: {
        anchor: {
            default: "bottom-left",
            type: String,
        },
    },
};
</script>

<style lang="scss" scoped>
.v-dropdown {
    background-color: var(--dropdown-background);
    border-radius: 0.2rem;
    box-shadow: 0 0.1rem 0.5rem var(--dropdown-shadow);
    color: var(--dropdown-color);
    cursor: default;
    display: flex;
    flex-direction: column;
    line-height: 2rem;
    padding: 0.35rem 0;
    position: absolute;
    z-index: 400;

    &.bottom-left,
    &.bottom-right {
        top: 100%;
    }

    &.top-left,
    &.top-right {
        bottom: 100%;
    }

    &.bottom-left,
    &.top-left {
        left: 0;
    }

    &.bottom-right,
    &.top-right {
        right: 0;
    }

    &.left,
    &.right {
        top: -0.35rem;
    }

    &.left {
        left: 100%;
    }

    &.right {
        right: 100%;
    }

    ::v-deep() {
        > .v-checkbox,
        > .v-label {
            padding: 0 1rem;
        }

        > .v-button,
        > .v-columns > .v-button {
            display: block;
            margin: 0;

            > .v-button-content {
                border-radius: 0;
                justify-content: flex-start;
                padding: 0 1rem;

                span {
                    flex-grow: 1;
                }
            }

            > .v-button-content:active,
            > .v-button-content:focus,
            > .v-button-content:hover {
                background-color: var(--accent);
                color: white;
            }
        }

        > .v-button.active > .v-button-content {
            background-color: var(--accent);
            color: white;
        }

        .v-checkbox {
            width: 100%;
        }

        .v-columns {
            column-count: 2;
        }

        > .v-divider {
            margin: 0.35rem 0;
        }

        > .v-label {
            color: var(--light-gray);
        }
    }
}
</style>
