/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {createApp} from "vue";

import main from "./main.vue";
import components from "./components";
import i18n from "./i18n";
import router from "./router";
import store from "./store";
import ui from "./ui";

// prettier-ignore
createApp(main)
    .use(components)
    .use(i18n)
    .use(router)
    .use(store)
    .use(ui)
    .mount(document.body);
