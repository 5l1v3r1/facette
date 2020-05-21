<template>
    <v-content :class="{compact: !erred && section === 'layout'}">
        <v-modal-confirm></v-modal-confirm>
        <v-modal-dashboard-item></v-modal-dashboard-item>
        <v-modal-template-variable :allow-fixed="link"></v-modal-template-variable>

        <v-toolbar clip="content">
            <v-button
                icon="save"
                :disabled="erred || saving || !validity"
                @click="save(true)"
                v-if="prevRoute.name !== 'dashboards-show' && !template && modifiers.alt"
            >
                {{ $t("labels.saveAndGo") }}
            </v-button>

            <v-button icon="save" :disabled="erred || saving || !validity" @click="save(false)" v-else>
                {{ $t(`labels.${template ? "templates" : "dashboards"}.save`) }}
            </v-button>

            <v-button icon="trash" @click="deleteDashboard()" v-if="!erred && edit && modifiers.alt">
                {{ $t("labels.delete") }}
            </v-button>

            <v-button :disabled="erred" :to="prevRoute || {name: 'admin-dashboards-list'}" v-else>
                {{ $t("labels.cancel") }}
            </v-button>

            <v-divider vertical></v-divider>

            <v-button icon="undo" :disabled="erred || loading || !guarded" @click="reset()">
                {{ $t("labels.reset") }}
            </v-button>

            <template v-if="dashboard && edit">
                <v-spacer></v-spacer>

                <v-label class="modified" v-if="dashboard.modified">
                    {{ $t("messages.lastModified", [formatDate(dashboard.modified, $t("date.long"))]) }}
                </v-label>
            </template>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <v-message-error @retry="reset(true)" v-else-if="erred"></v-message-error>

        <template v-else>
            <h1 v-if="section === 'general'">{{ $t("labels.general") }}</h1>

            <v-form class="third" @validity="onValidity" v-show="section === 'general'">
                <v-label>{{ $t("labels.name") }}</v-label>
                <v-input
                    ref="name"
                    required
                    :custom-validity="conflictCustomValidity"
                    :delay="200"
                    :pattern="namePattern"
                    v-autofocus.select
                    v-model="dashboard.name"
                ></v-input>

                <template v-if="link">
                    <v-label>{{ $tc("labels.templates._", 1) }}</v-label>
                    <v-flex class="template">
                        <v-select
                            :options="templates"
                            :placeholder="$t('labels.templates.select')"
                            v-model="dashboard.link"
                        >
                            <template slot="dropdown-placeholder" v-if="templates.length === 0">
                                <v-label>{{ $t("messages.templates.none") }}</v-label>
                            </template>
                        </v-select>

                        <v-button
                            icon="pencil-alt"
                            :to="{name: 'admin-dashboards-edit', params: {id: dashboard.link}}"
                            :style="{visibility: dashboard.link ? 'visible' : 'hidden'}"
                        >
                            {{ $t("labels.templates.edit") }}
                        </v-button>
                    </v-flex>
                </template>

                <template v-else>
                    <v-label>{{ $t("labels.title") }}</v-label>
                    <v-input
                        :delay="200"
                        :help="$t('help.dashboards.title')"
                        v-model="dashboard.options.title"
                    ></v-input>
                </template>
            </v-form>

            <template v-if="link">
                <template v-if="section === 'variables' && variables.length > 0">
                    <h1>{{ $t("labels.variables._") }}</h1>

                    <v-form-template-variables
                        :defined="dashboard.options.variables"
                        :parsed="variables"
                    ></v-form-template-variables>
                </template>
            </template>

            <template v-else>
                <template v-if="section === 'layout'">
                    <v-grid :add-handler="addItem" :layout="dashboard.layout" v-model="dashboard.items">
                        <template slot-scope="item">
                            <v-chart
                                :legend="item.value.options.legend"
                                @click.native="editItem(item.index)"
                                v-model="dashboardRefs[`chart|${item.value.options.id}`]"
                                v-if="item.value.type === 'chart'"
                            >
                            </v-chart>
                        </template>
                    </v-grid>
                </template>
            </template>
        </template>
    </v-content>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import merge from "lodash/merge";
import {Component, Mixins, Watch} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

import {ModalConfirmParams} from "@/src/components/modal/confirm.vue";
import {conflictCustomValidity} from "@/src/helpers/api";
import {beforeRoute} from "@/src/helpers/route";
import {CustomMixins} from "@/src/mixins";
import {ModalDashboardItemParams} from "@/src/views/admin/components/modal/dashboard-item.vue";

import {namePattern} from "..";

const defaultDashboard: Dashboard = {
    id: "",
    name: "",
    options: {
        title: "",
        variables: [],
    },
    layout: {
        columns: 1,
        rowHeight: 260,
        rows: 1,
    },
    items: [],
};

const defaultDashboardLinked: Dashboard = {
    id: "",
    name: "",
    options: {
        variables: [],
    },
};

const defaultSection = "general";

@Component({
    beforeRouteLeave: beforeRoute,
    beforeRouteUpdate: beforeRoute,
})
export default class Edit extends Mixins<CustomMixins>(CustomMixins) {
    public conflictCustomValidity!: (value: string) => Promise<string>;

    public dashboard: Dashboard | null = null;

    public dashboardModel: Dashboard | null = null;

    public dashboardRefs: Record<string, unknown> = {};

    public loading = true;

    public linked: Dashboard | null = null;

    public namePattern = namePattern;

    public resolveData: Record<string, string> = {};

    public saving = false;

    public section: string = defaultSection;

    public template = false;

    public templates: Array<SelectOption> = [];

    public validity = false;

    public variables: Array<TemplateVariable> = [];

    private unwatchDashboard: (() => void) | null = null;

    public created(): void {
        this.conflictCustomValidity = conflictCustomValidity(this, "dashboards");
    }

    public mounted(): void {
        this.reset(true);
    }

    public beforeDestroy(): void {
        if (this.unwatchDashboard) {
            this.unwatchDashboard();
        }
    }

    public addItem(layout: GridItemLayout): void {
        if (this.dashboard === null) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        }

        const items = this.dashboard.items as Array<DashboardItem>;

        this.$components.modal(
            "dashboard-item",
            {
                item: {layout: cloneDeep(layout)},
                step: 1,
            } as ModalDashboardItemParams,
            (value: DashboardItem) => {
                if (value) {
                    items.push(value);
                }
            },
        );
    }

    public deleteDashboard(apply = false): void {
        if (this.dashboard === null) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        }

        if (apply) {
            this.$http.delete(`/api/v1/dashboards/${this.dashboard.id}`).then(() => {
                this.unguard();
                this.$components.notify(this.$tc("messages.dashboards.deleted", 1) as string, "success");
                this.$router.push(
                    this.prevRoute?.name === "dashboards-show"
                        ? {name: "dashboards-home"}
                        : {name: "admin-dashboards-list", query: this.template ? {kind: "template"} : {}},
                );
            }, this.handleError());

            return;
        }

        this.$components.modal(
            "confirm",
            {
                button: {
                    icon: "trash",
                    label: this.$tc(`labels.${this.params.type}.delete`, 1),
                    danger: true,
                },
                message: this.$tc(`messages.${this.params.type}.delete`, 1, this.dashboard),
            } as ModalConfirmParams,
            (value: boolean) => {
                if (value) {
                    this.deleteDashboard(true);
                }
            },
        );
    }

    public get dynamicVariables(): Array<TemplateVariable> {
        return this.dashboard?.options?.variables?.filter(variable => variable.dynamic) ?? [];
    }

    public get edit(): boolean {
        return this.params.id !== "link" && this.params.id !== "new";
    }

    public editItem(index: number): void {
        if (this.dashboard === null || !this.dashboard.items) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        }

        const items = this.dashboard.items as Array<DashboardItem>;

        this.$components.modal(
            "dashboard-item",
            {
                item: cloneDeep(items[index]),
                step: 2,
            } as ModalDashboardItemParams,
            (value: DashboardItem) => {
                if (value) {
                    items.splice(index, 1, value);
                }
            },
        );
    }

    public get link(): boolean {
        return this.params.id === "link" || Boolean(this.dashboard?.link);
    }

    @Watch("params.id")
    public onParamsID(to: string, from: string): void {
        if (to !== from) {
            this.section = defaultSection;
            this.reset(true);
        }
    }

    @Watch("$route.hash", {immediate: true})
    public onRouteHash(to: string): void {
        this.section = to ? to.substr(1) : defaultSection;
    }

    public onValidity(to: boolean): void {
        this.validity = to;
    }

    public removeItem(index: number): void {
        if (this.dashboard === null || !this.dashboard.items) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        }

        this.dashboard.items.splice(index, 1);
    }

    public async reset(apply = false): Promise<void> {
        if (!apply) {
            this.$components.modal(
                "confirm",
                {
                    button: {
                        label: this.$t("labels.dashboards.reset"),
                        danger: true,
                    },
                    message: this.$t("messages.confirmLeave"),
                } as ModalConfirmParams,
                (value: boolean) => {
                    if (value) {
                        this.reset(true);
                    }
                },
            );

            return;
        }

        if (this.unwatchDashboard) {
            this.unwatchDashboard();
        }

        this.validity = false;
        this.dashboardRefs = {};

        if (!this.edit) {
            let dashboard: Dashboard;

            if (this.link) {
                await this.getTemplates();
                dashboard = cloneDeep(defaultDashboardLinked);
                this.unwatchDashboard = this.$watch("dashboard", this.onDashboardLinked, {deep: true});
            } else {
                dashboard = cloneDeep(defaultDashboard);
                this.unwatchDashboard = this.$watch("dashboard", this.onDashboard, {deep: true});
            }

            this.dashboard = dashboard;
            this.loading = false;
            this.guardWatch("dashboard");

            return;
        }

        this.loading = true;

        this.$http
            .get(`/api/v1/dashboards/${this.params.id}`)
            .then(response => response.json())
            .then(
                async (response: APIResponse<Dashboard>) => {
                    if (response?.data?.link) {
                        await this.getTemplates();
                        this.unwatchDashboard = this.$watch("dashboard", this.onDashboardLinked, {deep: true});
                        this.dashboard = merge({}, defaultDashboardLinked, response.data);
                    } else {
                        this.unwatchDashboard = this.$watch("dashboard", this.onDashboard, {deep: true});
                        this.dashboard = merge({}, defaultDashboard, response.data);
                    }

                    this.guardWatch("dashboard");
                },
                this.handleError(() => {
                    this.loading = false;
                }, true),
            );
    }

    public save(go: boolean): void {
        if (this.dashboard === null) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        } else if (this.saving) {
            return;
        }
        this.saving = true;

        const dashboard: Dashboard = cloneDeep(this.dashboard);
        dashboard.template = this.template;

        this.$http({
            url: `/api/v1/dashboards${this.edit ? `/${this.params.id}` : ""}`,
            method: this.edit ? "PUT" : "POST",
            body: dashboard,
        }).then(
            () => {
                this.unguard();
                this.$components.notify(this.$t("messages.dashboards.saved") as string, "success");
                this.$router.push(
                    go || this.prevRoute?.name === "dashboards-show"
                        ? {name: "dashboards-show", params: {id: dashboard.name}}
                        : {name: "admin-dashboards-list", query: this.template ? {kind: "template"} : {}},
                );
                this.saving = false;
            },
            this.handleError(() => {
                this.saving = false;
            }),
        );
    }

    private emitUpdate(): void {
        this.$parent.$emit(
            "dashboard-updated",
            this.link,
            this.dashboard?.items?.length ?? null,
            this.variables.length ?? null,
        );
    }

    private getTemplates(): PromiseLike<void> {
        return this.$http
            .get("/api/v1/dashboards", {params: {kind: "template"}})
            .then(response => response.json())
            .then((response: APIResponse<Array<Dashboard>>) => {
                if (response.data) {
                    this.templates = response.data.map((dashboard: Dashboard) => ({
                        label: dashboard.name,
                        value: dashboard.id,
                    }));
                }
            });
    }

    private onDashboard(to: Dashboard): void {
        if (!to.items) {
            this.emitUpdate();
            return;
        }

        const keys: Array<string> = [];
        const types: Array<DashboardItemType> = [];

        const req: Array<BulkRequest> = to.items.reduce((req: Array<BulkRequest>, item: DashboardItem) => {
            switch (item.type) {
                case "chart": {
                    const id = item.options?.id as string | undefined;
                    const key = `chart|${id}`;

                    if (
                        id !== undefined &&
                        !keys.includes(key) &&
                        (!this.dashboardRefs[key] || (this.dashboardRefs[key] as Chart).template)
                    ) {
                        req.push({
                            endpoint: `/charts/${id}/resolve`,
                            method: "POST",
                        });

                        keys.push(key);
                    }
                }
            }

            types.push(item.type);

            return req;
        }, []);

        if (req.length === 0) {
            this.emitUpdate();
            return;
        }

        this.$http
            .post("/api/v1/bulk", req)
            .then(response => response.json())
            .then(
                (response: APIResponse<Array<BulkResult>>) => {
                    if (response.data && response.data.filter(result => result.status >= 400).length > 0) {
                        this.$components.notify(this.$t("messages.error.bulk") as string, "error");
                        return;
                    }

                    response.data?.forEach((result: BulkResult, index: number) => {
                        switch (types[index]) {
                            case "chart": {
                                const chart = result.response.data as Chart;
                                this.dashboardRefs[`chart|${chart.id}`] = chart;
                                break;
                            }
                        }
                    });

                    this.loading = false;
                    this.emitUpdate();
                },
                () => {
                    this.loading = false;
                    this.emitUpdate();
                },
            );
    }

    private onDashboardLinked(to: Dashboard, from: Dashboard | undefined): void {
        // noop
    }
}
</script>

<style lang="scss" scoped>
@import "../mixins";

.v-content {
    ::v-deep {
        @include content;

        canvas {
            cursor: pointer;
        }
    }
}
</style>
