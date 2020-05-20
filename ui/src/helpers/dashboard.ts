/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

export function mapReferences(refs: Array<Reference>): Record<string, unknown> {
    return (
        refs.reduce((out: Record<string, unknown>, ref: Reference) => {
            switch (ref.type) {
                case "chart":
                    out[`chart|${(ref.value as Chart).id}`] = ref.value;
                    break;
            }

            return out;
        }, {}) ?? {}
    );
}

export async function resolveVariables(
    component: Vue,
    variables: Array<TemplateVariable>,
): Promise<Record<string, Array<string>>> {
    const req: Array<BulkRequest> = [];
    const labels: Array<string> = [];

    variables.forEach(variable => {
        if (variable.dynamic) {
            req.push({
                endpoint: `/labels/${variable.label}/values`,
                method: "GET",
                params: variable.filter ? {match: variable.filter} : undefined,
            });

            labels.push(variable.name);
        }
    });

    const data: Record<string, Array<string>> = {};

    if (req.length > 0) {
        await component.$http
            .post("/api/v1/bulk", req)
            .then(response => response.json())
            .then((response: APIResponse<Array<BulkResult>>) => {
                if (response.data && response.data.filter(result => result.status >= 400).length > 0) {
                    component.$components.notify(component.$t("messages.error.bulk") as string, "error");
                    return;
                }

                response.data?.forEach((result: BulkResult, index: number) => {
                    data[labels[index]] = result.response.data as Array<string>;
                });
            });
    }

    return data;
}
