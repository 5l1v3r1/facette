/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

export function conflictCustomValidity(
    component: Vue,
    type: string,
    id: string | null = null,
): (value: string) => Promise<string> {
    return (value: string): Promise<string> => {
        if (!value) {
            return Promise.resolve("");
        }

        return new Promise(resolve => {
            return component.$http
                .get(`/api/v1/${type}/${value}?excerpt=true`)
                .then(response => response.json())
                .then(
                    (response: APIResponse<ObjectMeta>) => {
                        resolve(response.data?.id !== id ? (component.$t(`messages.${type}.conflict`) as string) : "");
                    },
                    () => {
                        resolve("");
                    },
                );
        });
    };
}
