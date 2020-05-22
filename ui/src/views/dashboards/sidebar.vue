<template>
    <v-sidebar :active="sidebar">
        <v-toolbar clip="sidebar">
            <v-select
                :options="types"
                :search="false"
                @input="onType"
                v-model="type"
                v-if="$route.name === 'basket-show' || $route.name === 'dashboards-home'"
            ></v-select>

            <v-button
                exact
                icon="arrow-left"
                :shortcut="['alt+up', $t('labels.goto.parent')]"
                :to="{name: 'dashboards-show', params: {id: dashboard.parent}}"
                v-tooltip="{shortcut: 'alt+up'}"
                v-else-if="dashboard && dashboard.parent"
            >
                {{ $t("labels.goto.parent") }}
            </v-button>

            <v-button
                exact
                icon="arrow-left"
                :shortcut="['alt+up', $t('labels.goto.home')]"
                :to="{name: 'dashboards-home'}"
                v-tooltip="{shortcut: 'alt+up'}"
                v-else
            >
                {{ $t("labels.goto.home") }}
            </v-button>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <template v-else>
            <v-label>{{ title }}</v-label>

            <template v-if="params.type === 'dashboards'">
                <v-button
                    icon="folder"
                    :key="dashboard.id"
                    :to="{name: 'dashboards-show', params: {id: dashboard.name}}"
                    v-for="dashboard in dashboards"
                >
                    {{ dashboard.options && dashboard.options.title ? dashboard.options.title : dashboard.name }}
                </v-button>
            </template>

            <v-button
                class="empty"
                disabled
                icon="info-circle"
                v-if="
                    $route.name !== 'dashboards-home' &&
                    (!dashboard || !dashboard.items || dashboard.items.length === 0)
                "
            >
                {{ $t(`messages.${params.type}.${error === "notFound" ? "notFound" : "empty"}`) }}
            </v-button>

            <v-button
                :class="{unsupported: !checkType(item.type)}"
                :href="`#item${index}`"
                :icon="!checkType(item.type) ? 'exclamation-triangle' : null"
                :key="index"
                :ref="`item${index}`"
                @click="highlight($event)"
                @focus="highlight($event, index)"
                @focusout="highlight($event)"
                v-for="(item, index) in dashboard.items"
                v-else-if="dashboard"
            >
                {{ itemLabel(item) }}
            </v-button>
        </template>
    </v-sidebar>
</template>

<script lang="ts">
import {Component, Mixins, Watch} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

import {resolveOption} from "@/src/helpers/select";
import {CustomMixins} from "@/src/mixins";

import {itemTypes} from "./show.vue";

const types: Array<SelectOption> = [
    {label: "labels.home", value: "dashboards", icon: "home"},
    {label: "labels.basket._", value: "basket", icon: "shopping-basket"},
];

@Component
export default class Sidebar extends Mixins<CustomMixins>(CustomMixins) {
    public dashboard: Dashboard | null = null;

    public dashboards: Array<Dashboard> = [];

    public dashboardRefs: Record<string, unknown> = {};

    public loading = true;

    public type = "home";

    public mounted(): void {
        this.$parent.$on("dashboard-loaded", this.onDashboardLoaded);

        if (this.$route.name === "dashboards-home") {
            this.getDashboards();
        }
    }

    public updated(): void {
        this.type = this.params.type;
    }

    public beforeDestroy(): void {
        this.$parent.$off("dashboard-loaded", this.onDashboardLoaded);
    }

    public checkType(type: string): boolean {
        return itemTypes.includes(type);
    }

    public highlight(e: KeyboardEvent | MouseEvent, index?: number): void {
        if (e.ctrlKey || e.metaKey) {
            return;
        } else if (e.type === "click") {
            e.preventDefault();
            return;
        }

        history.replaceState(
            null,
            "",
            e.type === "focusout" ? window.location.pathname + window.location.search : `#item${index}`,
        );

        window.dispatchEvent(new Event("hashchange"));
    }

    public itemLabel(item: DashboardItem): string {
        switch (item.type) {
            case "chart": {
                if (item.options) {
                    const ref = this.dashboardRefs[`chart|${item.options.id}`] as Chart;
                    if (ref) {
                        return ref.options?.title ?? ref.name;
                    }

                    return item.options.id as string;
                }

                break;
            }

            default:
                return this.$t("labels.items.unsupported") as string;
        }

        return this.$t("labels.unnamed") as string;
    }

    public get title(): string {
        if (this.params.type === "basket") {
            return this.$t("labels.basket._") as string;
        } else if (this.dashboard) {
            return this.dashboard?.options?.title ?? this.dashboard.name;
        } else if (this.$route.name === "dashboards-show") {
            return this.$route.params.id;
        }

        return this.$tc("labels.dashboards._", 2) as string;
    }

    @Watch("$route.path")
    public onRoutePath(): void {
        Object.assign(this, {
            dashboard: null,
            dashboardsRefs: {},
            loading: true,
        });

        if (this.$route.name === "dashboards-home") {
            this.getDashboards();
        }
    }

    public onType(): void {
        switch (this.type) {
            case "basket":
                this.$router.push({name: "basket-show"});
                break;

            case "dashboards":
                this.$router.push({name: "dashboards-home"});
                break;
        }
    }

    public get types(): Array<SelectOption> {
        return types.map(option => resolveOption(this, option));
    }

    private getDashboards(): void {
        if (this.params.type === "basket") {
            this.loading = false;
            return;
        }

        this.$http
            .get("/api/v1/dashboards", {params: {parent: this.dashboard?.id}})
            .then(response => response.json())
            .then(
                (response: APIResponse<Array<Dashboard>>) => {
                    if (response.data) {
                        this.dashboards = response.data;
                    }
                    this.loading = false;
                },
                this.handleError(() => {
                    this.loading = false;
                }),
            );
    }

    private onDashboardLoaded(dashboard: Dashboard | null, dashboardRefs: Record<string, unknown> = {}): void {
        if (dashboard === null) {
            // Dashboard load failed, only disable loading
            this.loading = false;
            return;
        }

        Object.assign(this, {
            dashboard,
            dashboardRefs,
        });

        this.getDashboards();
    }
}
</script>

<style lang="scss" scoped>
.v-sidebar {
    .v-toolbar {
        .v-button,
        .v-select {
            flex-grow: 1;
        }

        .v-button ::v-deep .v-button-content {
            justify-content: flex-start;
        }
    }
}
</style>
