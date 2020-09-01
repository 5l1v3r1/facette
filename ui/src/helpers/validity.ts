/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import i18n from "@/i18n";
import api, {ObjectType} from "@/lib/api";

export function objectNameValidity(type: ObjectType, id: string | null = null): (value: string) => Promise<string> {
    return (value: string): Promise<string> => {
        if (!value) {
            return Promise.resolve("");
        }

        return new Promise(resolve => {
            api.object(type, value, {excerpt: true}).then(
                response => {
                    resolve(response.data?.id !== id ? i18n.global.t(`messages.${type}.conflict`) : "");
                },
                () => {
                    resolve("");
                },
            );
        });
    };
}
