/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import Vue from "vue";
import VueComponents from "./components/vue"; // TODO: extract
import VueResource from "vue-resource";

Vue.use(VueComponents);
Vue.use(VueResource);

Vue.config.productionTip = false;

import "./components";

import i18n from "./i18n";
import router from "./router";
import store from "./store";
import main from "./main.vue";

new Vue({
    i18n,
    router,
    store,
    render: h => h(main),
}).$mount(document.body);
