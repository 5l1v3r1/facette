<template>
    <v-content>
        <v-modal-confirm></v-modal-confirm>
        <v-modal-prompt></v-modal-prompt>

        <v-toolbar clip="content">
            <v-button
                icon="sync"
                :tooltip="$t('labels.refresh.list')"
                :shortcut="['r', $t('labels.refresh.list')]"
                @click="getObjects"
            ></v-button>

            <v-divider vertical></v-divider>

            <v-button
                icon="plus"
                :shortcut="['n', $t(`labels.${params.type}.new`)]"
                :to="{name: 'admin-providers-edit', params: {id: 'new'}}"
                :tooltip="''"
                v-if="params.type === 'providers'"
            >
                {{ $t("labels.providers.new") }}
            </v-button>

            <v-button icon="plus" v-else>
                {{ $t(`labels.${params.type}.new`) }}
                <template slot="dropdown">
                    <v-button
                        :shortcut="['n', $t(`labels.${params.type}.new`)]"
                        :to="{name: `admin-${params.type}-edit`, params: {id: 'new'}}"
                    >
                        {{ $t(`labels.${params.type}.new`) }}
                    </v-button>
                    <v-button
                        :shortcut="['shift+n', $t(`labels.templates.newFrom`)]"
                        :to="{name: `admin-${params.type}-edit`, params: {id: 'link'}}"
                    >
                        {{ $t("labels.templates.newFrom") }}
                    </v-button>
                </template>
            </v-button>

            <template v-if="params.type == 'providers'">
                <v-divider vertical></v-divider>

                <v-button
                    icon="play"
                    :disabled="
                        selection.length === 0 || (selection.enabled.length === 1 && selection.enabled[0] === true)
                    "
                    @click="toggleProviders(Object.values(selection.objects), true)"
                >
                    {{ $t("labels.providers.enable") }}
                </v-button>
                <v-button
                    icon="stop"
                    :disabled="
                        selection.length === 0 || (selection.enabled.length === 1 && selection.enabled[0] === false)
                    "
                    @click="toggleProviders(Object.values(selection.objects), false)"
                >
                    {{ $t("labels.providers.disable") }}
                </v-button>
                <v-button
                    :disabled="
                        selection.length === 0 ||
                        (selection.enabled.length === 1 && selection.enabled[0] === false) ||
                        polling
                    "
                    :icon="polling ? 'circle-notch' : 'arrow-alt-circle-down'"
                    :icon-spin="polling"
                    @click="pollProviders(Object.values(selection.objects))"
                >
                    {{ $t("labels.providers.poll") }}
                </v-button>
            </template>

            <v-divider vertical></v-divider>

            <v-button
                icon="trash"
                :disabled="selection.length === 0"
                :shortcut="['d', $t('labels.delete')]"
                :tooltip="''"
                @click="deleteObjects(Object.values(selection.objects))"
            >
                {{ $t("labels.delete") }}
            </v-button>

            <v-spacer></v-spacer>

            <v-input
                icon="filter"
                type="search"
                :placeholder="$t(`labels.${params.type}.filter`)"
                :shortcut="['f', $t(`labels.${params.type}.filter`)]"
                @clear="applyFilter"
                @focusout="applyFilter"
                @keypress.enter.native="applyFilter"
                v-model="tmpFilter"
            ></v-input>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <v-message-error @retry="getObjects" v-else-if="erred"></v-message-error>

        <template v-else>
            <h1>
                {{ $tc(`labels.${params.type}._`, 2) }}
                <span class="count" v-if="total">{{ total }}</span>
            </h1>

            <v-message class="selection" icon="clipboard-check" type="info" v-if="selection.length > 0">
                {{ $tc(`messages.${params.type}.selected`, selection.length, [selection.length]) }}
                <v-button icon="times-circle" :tooltip="$t('labels.clearSelection')" @click="clearSelection"></v-button>
            </v-message>

            <v-tab-bar v-model="options.kind" v-if="templatable">
                <v-tab value="plain">{{ $tc(`labels.${params.type}._`, 2) }}</v-tab>
                <v-tab value="template">{{ $tc("labels.templates._", 2) }}</v-tab>
            </v-tab-bar>

            <v-message type="info" v-if="objects.length === 0">
                {{ $t(`messages.${options.kind === "template" ? "templates" : params.type}.none`) }}
            </v-message>

            <template v-else>
                <v-table ref="table" selectable @selection="onSelection" v-model="objects">
                    <template slot="header">
                        <v-table-cell grow>
                            <v-flex class="sort" @click.native="toggleSort('name')">
                                {{ $t("labels.name") }}
                                <v-icon
                                    :icon="`chevron-${options.sort == '-name' ? 'down' : 'up'}`"
                                    v-if="options.sort == 'name' || options.sort == '-name'"
                                ></v-icon>
                            </v-flex>
                        </v-table-cell>

                        <v-table-cell>
                            <v-flex class="sort" @click.native="toggleSort('modified')">
                                {{ $t("labels.lastModified") }}
                                <v-icon
                                    :icon="`chevron-${options.sort == '-modified' ? 'down' : 'up'}`"
                                    v-if="options.sort == 'modified' || options.sort == '-modified'"
                                ></v-icon>
                            </v-flex>
                        </v-table-cell>

                        <v-table-cell></v-table-cell>
                    </template>

                    <template slot-scope="obj">
                        <v-table-cell grow>
                            <v-flex>
                                <router-link
                                    class="link"
                                    :to="{name: `admin-${params.type}-edit`, params: {id: obj.value.id}}"
                                >
                                    {{ obj.value.name }}
                                </router-link>

                                <v-icon
                                    class="linked"
                                    icon="link"
                                    v-tooltip="$t('labels.templates.instance')"
                                    v-if="templatable && obj.value.link"
                                ></v-icon>

                                <template v-else-if="params.type == 'providers'">
                                    <v-icon
                                        class="disabled"
                                        icon="stop-circle"
                                        v-tooltip="$t('labels.providers.disabled')"
                                        v-if="!obj.value.enabled"
                                    ></v-icon>
                                    <v-icon
                                        class="error"
                                        icon="exclamation-circle"
                                        v-tooltip="$t('messages.error._', [obj.value.error])"
                                        v-else-if="obj.value.error"
                                    ></v-icon>
                                    <v-icon
                                        class="enabled"
                                        icon="check-circle"
                                        v-tooltip="$t('labels.providers.enabled')"
                                        v-else
                                    ></v-icon>
                                </template>
                            </v-flex>
                        </v-table-cell>

                        <v-table-cell>
                            {{ formatDate(obj.value.modified, $t("date.long")) }}
                        </v-table-cell>

                        <v-table-cell>
                            <v-button class="reveal icon" dropdown-anchor="right" icon="ellipsis-v">
                                <template slot="dropdown">
                                    <template v-if="params.type === 'providers'">
                                        <v-button
                                            icon="arrow-alt-circle-right"
                                            :to="{
                                                name: 'admin-metrics-list',
                                                query: {filter: `{__provider__=${JSON.stringify(obj.value.name)}}`},
                                            }"
                                        >
                                            {{ $t("labels.goto.metrics") }}
                                        </v-button>
                                        <v-divider></v-divider>
                                    </template>

                                    <template v-else-if="templatable && options.kind === 'plain'">
                                        <v-button
                                            icon="arrow-alt-circle-right"
                                            :to="{name: `${params.type}-show`, params: {id: obj.value.name}}"
                                        >
                                            {{ $tc(`labels.goto.${params.type}`, 1) }}
                                        </v-button>
                                        <v-divider></v-divider>
                                    </template>

                                    <v-button icon="clone" @click="cloneObject(obj.value, null)">
                                        {{ $t("labels.clone") }}
                                    </v-button>

                                    <v-button icon="trash" @click="deleteObjects([obj.value], null)">
                                        {{ $t("labels.delete") }}
                                    </v-button>

                                    <template v-if="params.type == 'providers'">
                                        <v-divider></v-divider>
                                        <v-button
                                            icon="play"
                                            @click="toggleProviders([obj.value], true)"
                                            v-if="!obj.value.enabled"
                                        >
                                            {{ $t("labels.providers.enable") }}
                                        </v-button>
                                        <v-button icon="stop" @click="toggleProviders([obj.value], false)" v-else>
                                            {{ $t("labels.providers.disable") }}
                                        </v-button>
                                        <v-button
                                            icon="sync"
                                            :disabled="!obj.value.enabled"
                                            @click="pollProviders([obj.value])"
                                        >
                                            {{ $t("labels.providers.poll") }}
                                        </v-button>
                                    </template>
                                </template>
                            </v-button>
                        </v-table-cell>
                    </template>
                </v-table>

                <v-paging :handler="switchPage" :page="options.page" :page-size="limit" :total="total"></v-paging>
            </template>
        </template>
    </v-content>
</template>

<script lang="ts">
import uniq from "lodash/uniq";
import {Component, Mixins, Watch} from "vue-property-decorator";
import {Dictionary} from "vue-router/types/router";

import TableComponent from "@/src/components/vue/table/table.vue";
import {CustomMixins} from "@/src/mixins";
import {updateRouteQuery} from "@/src/router";

interface Options {
    filter: string;
    kind: string;
    page: number;
    sort: string;
}

@Component
export default class List extends Mixins<CustomMixins>(CustomMixins) {
    public limit = 20;

    public loading = true;

    public objects: Array<ObjectMeta> = [];

    public options: Options = {
        filter: "",
        kind: "plain",
        page: 1,
        sort: "name",
    };

    public polling = false;

    public selection: {
        objects: Record<string, ObjectMeta>;
        length: number;
        [key: string]: unknown;
    } = {
        objects: {},
        length: 0,
        enabled: [],
    };

    public tmpFilter = "";

    public total = 0;

    public applyFilter(): void {
        if (this.tmpFilter !== this.options.filter) {
            this.options.filter = this.tmpFilter;
        }
    }

    public clearSelection(): void {
        this.selection = {
            objects: {},
            length: 0,
            enabled: [],
        };

        const table: TableComponent = this.$refs.table as TableComponent;
        if (table) {
            table.select([]);
        }
    }

    public cloneObject(obj: ObjectMeta, name?: string): void {
        if (name) {
            this.$http
                .post(`/api/v1/${this.params.type}?inherit=${obj.id}`, {
                    name,
                })
                .then(() => {
                    this.getObjects();
                }, this.handleError());
            return;
        }

        this.$components.modal(
            "prompt",
            {
                button: {
                    label: this.$t("labels.clone"),
                    primary: true,
                },
                input: {
                    value: `${obj.name}-clone`,
                },
                message: this.$t(`labels.${this.params.type}.name`),
            },
            (value: string) => {
                if (value) {
                    this.cloneObject(obj, value);
                }
            },
        );
    }

    public deleteObjects(objects: Array<ObjectMeta>, apply = false): void {
        if (apply) {
            this.$http
                .post(
                    "/api/v1/bulk",
                    objects.map(obj => ({
                        endpoint: `/${this.params.type}/${obj.id}`,
                        method: "DELETE",
                    })),
                )
                .then(response => response.json())
                .then((response: APIResponse<Array<BulkResult>>) => {
                    if (response.data && response.data.filter(result => result.status >= 400).length > 0) {
                        this.$components.notify(this.$t("messages.error.bulk") as string, "error");
                    } else {
                        this.$components.notify(
                            this.$tc(`messages.${this.params.type}.deleted`, objects.length),
                            "success",
                        );
                    }

                    // Remove any selection reference to deleted objects
                    objects.forEach(obj => delete this.selection.objects[obj.id]);

                    this.getObjects();
                }, this.handleError());

            return;
        }

        this.$components.modal(
            "confirm",
            {
                button: {
                    icon: "trash",
                    label: this.$tc(`labels.${this.params.type}.delete`, objects.length),
                    danger: true,
                },
                message: this.$tc(
                    `messages.${this.params.type}.delete`,
                    objects.length,
                    objects.length > 1 ? {count: objects.length} : objects[0],
                ),
            },
            (value: boolean) => {
                if (value) {
                    this.deleteObjects(objects, true);
                }
            },
        );
    }

    @Watch("objects")
    public onObjects(to: Array<ObjectMeta>): void {
        if (this.selection.length > 0) {
            this.$nextTick(() => {
                const table = this.$refs.table as TableComponent;
                if (table) {
                    table.select(to.filter(obj => this.selection.objects[obj.id]));
                }
            });
        }
    }

    @Watch("options", {deep: true})
    public onOptions(to: Options): void {
        const q: Dictionary<string> = {};

        if (to.filter !== "") {
            q.filter = to.filter;
        }
        if (to.kind !== "plain") {
            q.kind = to.kind;
        }
        if (to.page !== 1) {
            q.page = to.page.toString();
        }
        if (to.sort !== "name") {
            q.sort = to.sort;
        }

        updateRouteQuery(this.$route, q);

        this.getObjects();
    }

    @Watch("$route.query", {immediate: true})
    public onRouteQuery(to: Dictionary<string>, from: Dictionary<string> | undefined): void {
        if (from === undefined) {
            const options: Options = {
                filter: "",
                kind: "plain",
                page: 1,
                sort: "name",
            };

            if (to.filter) {
                options.filter = to.filter;
                this.tmpFilter = to.filter;
            }
            if (to.kind) {
                options.kind = to.kind;
            }
            if (to.page) {
                options.page = parseInt(to.page, 10) ?? 1;
            }
            if (to.sort) {
                options.sort = to.sort;
            }

            Object.assign(this.options, options);
        }

        // Get objects if no query parameter (i.e. initial get)
        if (Object.keys(to).length === 0) {
            this.getObjects();
        }
    }

    public onSelection(to: Array<ObjectMeta>): void {
        const ids = this.objects.map(obj => obj.id);

        const objects: Record<string, ObjectMeta> = Object.values(this.selection.objects)
            .filter(obj => !ids.includes(obj.id))
            .concat(to)
            .reduce((objects: Record<string, ObjectMeta>, obj: ObjectMeta) => {
                objects[obj.id] = obj;
                return objects;
            }, {});

        this.selection = {
            objects,
            length: Object.keys(objects).length,
            enabled:
                this.params.type === "providers"
                    ? uniq(Object.values(objects as Record<string, Provider>).map(obj => obj.enabled))
                    : [],
        };
    }

    public pollProviders(objects: Array<ObjectMeta>): void {
        this.polling = true;

        this.$http
            .post(
                "/api/v1/bulk",
                objects.map(obj => ({
                    endpoint: `/${this.params.type}/${obj.id}/poll`,
                    method: "POST",
                })),
            )
            .then(response => response.json())
            .then(
                (response: APIResponse<Array<BulkResult>>) => {
                    if (response.data && response.data.filter(result => result.status >= 400).length > 0) {
                        this.$components.notify(this.$t("messages.error.bulk") as string, "error");
                    }

                    this.polling = false;
                    this.getObjects();
                },
                this.handleError(() => {
                    this.polling = false;
                }),
            );
    }

    public switchPage(page: number): void {
        this.options.page = page;
    }

    public get templatable(): boolean {
        return this.params.type === "charts" || this.params.type === "dashboards";
    }

    public toggleProviders(objects: Array<ObjectMeta>, state: boolean, apply = false): void {
        if (apply) {
            this.$http
                .post(
                    "/api/v1/bulk",
                    objects.map(obj => ({
                        endpoint: `/providers/${obj.id}`,
                        method: "PATCH",
                        data: {
                            enabled: state,
                        },
                    })),
                )
                .then(response => response.json())
                .then((response: APIResponse<Array<BulkResult>>) => {
                    if (response.data && response.data.filter(result => result.status >= 400).length > 0) {
                        this.$components.notify(this.$t("messages.error.bulk") as string, "error");
                    } else {
                        this.$components.notify(
                            this.$tc(`messages.${this.params.type}.${state ? "enabled" : "disabled"}`, objects.length),
                            "success",
                        );
                    }

                    this.getObjects();
                }, this.handleError());

            return;
        }

        this.$components.modal(
            "confirm",
            {
                button: {
                    label: this.$tc(`labels.providers.${state ? "enable" : "disable"}`, objects.length),
                    danger: !state,
                    primary: state,
                },
                message: this.$tc(
                    `messages.providers.${state ? "enable" : "disable"}`,
                    objects.length,
                    objects.length > 1 ? {count: objects.length} : objects[0],
                ),
            },
            (value: boolean) => {
                if (value) {
                    this.toggleProviders(objects, state, true);
                }
            },
        );
    }

    public toggleSort(key: string): void {
        const desc: boolean = this.options.sort.startsWith("-");
        const sort: string = desc ? this.options.sort.substr(1) : this.options.sort;
        this.options.sort = key === sort && !desc ? `-${key}` : key;
    }

    private getObjects(): void {
        const params: ListParams = {
            limit: this.limit,
            offset: (this.options.page - 1) * this.limit ?? undefined,
        };

        if (this.templatable && this.options.kind) {
            params.kind = this.options.kind;
        }

        if (this.options.sort) {
            params.sort = this.options.sort;
        }

        if (this.options.filter) {
            let parts: Array<string> = this.options.filter.split(" ");

            if (this.params.type === "providers") {
                parts = parts.filter((part: string) => !part.startsWith("enabled:"));

                const enabled: string | undefined = parts.filter((part: string) => part.startsWith("enabled:"))[0];
                if (enabled) {
                    params.enabled = `${enabled.substr(8)}`;
                }
            }

            if (parts.length > 0) {
                params.name = `~(?:${parts.join("|")})`;
            }
        }

        this.loading = true;

        this.$http
            .get(`/api/v1/${this.params.type}`, {params})
            .then(response => response.json())
            .then(
                (response: APIResponse<Array<ObjectMeta>>) => {
                    const pagesCount: number | undefined = response.total
                        ? Math.ceil(response.total / this.limit)
                        : undefined;

                    // Switch back to first/last page if current empty
                    if (response.data?.length === 0 && this.options.page > 1) {
                        this.options.page = pagesCount !== undefined ? pagesCount : 1;
                        return;
                    }

                    this.objects = response.data ?? [];
                    this.total = response.total ?? 0;
                    this.loading = false;
                },
                this.handleError(() => {
                    this.loading = false;
                }),
            );
    }
}
</script>

<style lang="scss" scoped>
@import "./mixins";

.v-content {
    @include content;

    .v-toolbar .v-input {
        width: 20rem;
    }

    .v-table {
        .v-flex {
            &.row .v-icon {
                margin-left: 0.35rem;

                &.disabled,
                &.linked {
                    color: var(--light-gray);
                }

                &.enabled {
                    color: var(--green);
                }

                &.error {
                    color: var(--red);
                }
            }

            &.sort {
                cursor: pointer;
                display: inline-flex;
                height: 2rem;
            }
        }

        .v-table-row {
            &.selected {
                .v-flex.row .v-icon {
                    color: var(--toolbar-row-selected-color);
                    opacity: 0.65;
                }
            }
        }
    }
}
</style>
