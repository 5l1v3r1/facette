/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import Vue from "vue";
import VueI18n from "vue-i18n";

Vue.use(VueI18n);

import {messages} from "./locales";

export default new VueI18n({
    fallbackLocale: "en",
    locale: "en",
    messages,
});
