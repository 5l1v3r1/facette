/**
 * Copyright (c) 2020, The Facette Authors
 *
 * Licensed under the terms of the BSD 3-Clause License; a copy of the license
 * is available at: https://opensource.org/licenses/BSD-3-Clause
 */

import {Labels} from "@/lib/labels";

export const prefix = "/api/v1";

export type ObjectType = "charts" | "dashboards" | "providers";

function createURL(input: string, params?: Record<string, any>): string {
    const url = new URL(input, location.origin);

    if (params !== undefined) {
        const search = new URLSearchParams();

        Object.keys(params).forEach(key => {
            if (params[key] !== undefined) {
                search.set(key, params[key]);
            }
        });

        url.search = search.toString();
    }

    return url.toString();
}

export async function onFetch(response: Response): Promise<any> {
    let data: Promise<any>;
    if (response.headers.get("Content-Type")?.includes("application/json")) {
        data = response.json();
    } else {
        data = response.text();
    }

    if (response.status >= 400) {
        return Promise.reject(await data);
    }

    return data;
}

export class API {
    public bulk(req: Array<BulkRequest>): Promise<APIResponse<Array<BulkResult>>> {
        return fetch(`${prefix}/bulk`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(req),
        })
            .then(onFetch)
            .then((response: APIResponse<Array<BulkResult>>) => {
                if (response.data && response.data.filter(result => result.status >= 400).length > 0) {
                    return Promise.reject(response);
                }

                return response;
            });
    }

    public cloneObject(type: ObjectType, id: string, obj: Record<string, unknown>): Promise<void> {
        return fetch(createURL(`${prefix}/${type}`, {copy: id}), {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(obj),
        }).then(onFetch);
    }

    public delete(type: ObjectType, id: string): Promise<void> {
        return fetch(`${prefix}/${type}/${id}`, {method: "DELETE"}).then(onFetch);
    }

    public metrics(params?: ListParams): Promise<APIResponse<Array<Labels>>> {
        return fetch(createURL(`${prefix}/metrics`, params))
            .then(onFetch)
            .then(response => {
                if (response.data) {
                    response.data = response.data?.map((metric: string) => new Labels(metric));
                }

                return Promise.resolve(response);
            });
    }

    public object<T extends ObjectBase>(type: ObjectType, id: string, params?: ListParams): Promise<APIResponse<T>> {
        return fetch(createURL(`${prefix}/${type}/${id}`, params)).then(onFetch);
    }

    public objects<T extends ObjectBase>(type: ObjectType, params?: ListParams): Promise<APIResponse<Array<T>>> {
        return fetch(createURL(`${prefix}/${type}`, params)).then(onFetch);
    }

    public options(): Promise<APIResponse<Options>> {
        return fetch(prefix, {method: "OPTIONS"}).then(onFetch);
    }

    public query(query: SeriesQuery): Promise<APIResponse<SeriesResult>> {
        return fetch(`${prefix}/query`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(query),
        }).then(onFetch);
    }

    public resolveObject<T extends ObjectBase>(type: ObjectType, id: string): Promise<APIResponse<T>> {
        return fetch(`${prefix}/${type}/${id}/resolve`, {method: "POST"}).then(onFetch);
    }

    public saveObject<T extends ObjectBase>(type: ObjectType, obj: T): Promise<void> {
        let url = `${prefix}/${type}`;
        let method = "POST";

        if (obj.id) {
            url += `/${obj.id}`;
            method = "PUT";
        }

        return fetch(url, {
            method,
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(obj),
        })
            .then(response => {
                if (method === "POST") {
                    // Apply identifier received from back-end to the current
                    // object instance
                    obj.id = response.headers.get("Location")?.substr(url.length + 1) ?? "";
                }

                return Promise.resolve(response);
            })
            .then(onFetch);
    }

    public testObject<T extends ObjectBase>(type: ObjectType, obj: T): Promise<APIResponse<TestResult>> {
        return fetch(`${prefix}/${type}/test`, {
            method: "POST",
            headers: {"Content-Type": "application/json"},
            body: JSON.stringify(obj),
        }).then(onFetch);
    }

    public version(): Promise<APIResponse<Version>> {
        return fetch(`${prefix}/version`).then(onFetch);
    }
}

export default new API();
