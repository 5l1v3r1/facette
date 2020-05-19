<template>
    <v-content>
        <v-toolbar clip="content"></v-toolbar>

        <h1>{{ $t("labels.home") }}</h1>
    </v-content>
</template>

<script lang="ts">
import {Component, Mixins} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

@Component
export default class Home extends Mixins<CustomMixins>(CustomMixins) {}
</script>
