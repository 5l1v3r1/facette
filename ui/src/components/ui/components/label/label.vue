<template>
    <div class="v-label">
        <v-icon aria-hidden="true" :badge="iconBadge" :icon="icon" :spin="iconSpin" v-if="icon !== null"></v-icon>

        <slot></slot>
    </div>
</template>

<script lang="ts">
export default {
    props: {
        icon: {
            default: null,
            type: String,
        },
        iconBadge: {
            default: null,
            type: [Number, String],
        },
        iconSpin: {
            default: false,
            type: Boolean,
        },
    },
};
</script>

<style lang="scss" scoped>
.v-label {
    align-items: center;
    display: flex;
    line-height: var(--label-height);
    min-height: var(--label-height);

    .v-icon {
        margin-right: 0.5rem;
    }
}
</style>
