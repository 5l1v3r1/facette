<template>
    <div class="v-spacer"></div>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator";

@Component
export default class SpacerComponent extends Vue {}
</script>

<style lang="scss" scoped>
.v-spacer {
    flex-grow: 1;
}
</style>
