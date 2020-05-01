/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {xxHash32} from "js-xxhash";

export function hash(input: unknown): string {
    return xxHash32(JSON.stringify(input), 0).toString(16);
}
