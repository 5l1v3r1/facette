/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {AppTheme} from "types/ui";

import dark from "./dark";
import light from "./light";

const themes: Record<string, AppTheme> = {
    dark,
    light,
};

export default themes;

export function detectTheme(): string {
    return matchMedia("(prefers-color-scheme: dark)").matches ? "dark" : "light";
}
