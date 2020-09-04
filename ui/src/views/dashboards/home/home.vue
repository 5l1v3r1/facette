<!--
 Copyright (c) 2020, The Facette Authors

 Licensed under the terms of the BSD 3-Clause License; a copy of the license
 is available at: https://opensource.org/licenses/BSD-3-Clause
-->

<template>
    <v-content>
        <v-toolbar clip="content">
            <template v-if="basket.length > 0">
                <v-spacer></v-spacer>

                <v-button dropdown-anchor="bottom-right" icon="shopping-basket" :icon-badge="basket.length">
                    <template v-slot:dropdown>
                        <v-button icon="eye" :to="{name: 'basket-show'}">
                            {{ i18n.t("labels.basket.preview") }}
                        </v-button>

                        <v-divider></v-divider>

                        <v-button icon="broom" @click="clearBasket">
                            {{ i18n.t("labels.basket.clear") }}
                        </v-button>
                    </template>
                </v-button>
            </template>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>
    </v-content>
</template>

<script lang="ts">
import {computed, nextTick, onMounted} from "vue";
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

        const basket = computed(() => {
            return store.state.basket;
        });

        const clearBasket = (): void => {
            store.commit("basket", []);
        };

        onMounted(() => {
            ui.title();

            // TODO: implement home and get rid of this hack
            store.commit("loading", true);
            nextTick(() => store.commit("loading", false));
        });

        return {
            basket,
            clearBasket,
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
