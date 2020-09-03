<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-content>
        <v-toolbar clip="content"></v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <template v-else>
            <h1>{{ i18n.t("labels.home") }}</h1>
        </template>
    </v-content>
</template>

<script lang="ts">
import {nextTick, onMounted} from "vue";
import {useI18n} from "vue-i18n";
import {useStore} from "vuex";

import common from "@/common";
import {useUI} from "@/components/ui";
import {State} from "@/store";

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();
        const store = useStore<State>();
        const ui = useUI();

        const {loading} = common;

        onMounted(() => {
            ui.title();

            // TODO: implement home and get rid of this hack
            store.commit("loading", true);
            nextTick(() => store.commit("loading", false));
        });

        return {
            i18n,
            loading,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../../mixins";

.v-content {
    @include content;
}
</style>
