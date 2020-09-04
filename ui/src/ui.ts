/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {createUI} from "@/components/ui"; // FIXME: extract

import store from "./store";

export default createUI({
    theme: store.state.theme ?? undefined,
    title: "Facette",
});
