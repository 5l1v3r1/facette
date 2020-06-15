<template>
    <v-content>
        <v-modal-confirm></v-modal-confirm>
        <v-modal-provider-filter></v-modal-provider-filter>

        <v-toolbar clip="content">
            <v-button
                icon="save"
                :disabled="erred || formFailed || saving || testing || !validity"
                @click="save(true)"
                v-if="modifiers.alt"
            >
                {{ $t("labels.saveAndGo") }}
            </v-button>

            <v-button
                icon="save"
                :disabled="erred || formFailed || saving || testing || !validity"
                @click="save(false)"
                v-else
            >
                {{ $t("labels.providers.save") }}
            </v-button>

            <v-button icon="trash" @click="deleteProvider()" :disabled="testing" v-if="!erred && edit && modifiers.alt">
                {{ $t("labels.delete") }}
            </v-button>

            <v-button :disabled="erred" :to="{name: 'admin-providers-list'}" v-else>
                {{ $t("labels.cancel") }}
            </v-button>

            <v-divider vertical></v-divider>

            <v-button icon="undo" :disabled="erred || loading || testing || !guarded" @click="reset()">
                {{ $t("labels.reset") }}
            </v-button>

            <v-divider vertical></v-divider>

            <v-button icon="clipboard-check" :disabled="erred || formFailed || testing || !validity" @click="test()">
                {{ $t("labels.providers.test") }}
            </v-button>

            <template v-if="provider && edit">
                <v-spacer></v-spacer>

                <v-label class="modified" v-if="provider.modified">
                    {{ $t("messages.lastModified", [formatDate(provider.modified, $t("date.long"))]) }}
                </v-label>
            </template>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <v-message-error v-else-if="erred"></v-message-error>

        <template v-else>
            <h1 v-if="section === 'general'">{{ $t("labels.general") }}</h1>

            <v-form class="third" @validity="onValidity" v-show="section === 'general'">
                <v-label>{{ $t("labels.name") }}</v-label>
                <v-input
                    ref="name"
                    required
                    :custom-validity="conflictCustomValidity"
                    :help="$t('help.providers.name')"
                    :pattern="namePattern"
                    v-autofocus.select
                    v-model="provider.name"
                ></v-input>

                <v-label>{{ $tc("labels.connectors._", 1) }}</v-label>
                <v-select
                    class="half"
                    required
                    :options="providers"
                    :placeholder="$t('labels.connectors.select')"
                    v-model="provider.connector.type"
                ></v-select>

                <v-label>{{ $t("labels.providers.pollInterval") }}</v-label>
                <v-input
                    class="half"
                    :help="$t('help.providers.pollInterval')"
                    :placeholder="$t('labels.placeholders.example', ['1m, 30m, 2h'])"
                    v-model="provider.pollInterval"
                ></v-input>

                <v-message type="error" v-if="formFailed">
                    {{ $t("messages.providers.loadFailed", [provider.connector.type]) }}
                </v-message>

                <component
                    class="v-form"
                    :is="form"
                    :settings="provider.connector.settings"
                    v-else-if="provider.connector && form !== null"
                ></component>
            </v-form>

            <template v-if="section === 'filters'">
                <h1>{{ $tc("labels.filters._", 2) }}</h1>

                <v-form class="half">
                    <v-message type="info" v-if="!provider.filters || provider.filters.length === 0">
                        {{ $t("messages.filters.none") }}
                    </v-message>

                    <v-table draggable v-model="provider.filters" v-else>
                        <template slot="header">
                            <v-table-cell>
                                {{ $t("labels.filters.action._") }}
                            </v-table-cell>

                            <v-table-cell grow>
                                {{ $t("labels.properties") }}
                            </v-table-cell>

                            <v-table-cell></v-table-cell>
                        </template>

                        <template slot-scope="item">
                            <v-table-cell>
                                {{ item.value.action }}
                            </v-table-cell>

                            <v-table-cell grow>
                                <v-flex>
                                    <v-labels :labels="{[item.value.label]: item.value.pattern}"></v-labels>

                                    <template v-if="item.value.action === 'relabel'">
                                        <v-icon icon="arrow-right"></v-icon>
                                        <v-labels :labels="item.value.targets"></v-labels>
                                    </template>

                                    <template v-else-if="item.value.action === 'rewrite'">
                                        <v-icon icon="arrow-right"></v-icon>
                                        <v-labels :labels="{[item.value.label]: item.value.into}"></v-labels>
                                    </template>
                                </v-flex>
                            </v-table-cell>

                            <v-table-cell>
                                <v-button
                                    class="reveal"
                                    icon="pencil-alt"
                                    @click="editFilter(item.index)"
                                    v-tooltip="$t('labels.filters.edit')"
                                ></v-button>
                                <v-button
                                    class="reveal"
                                    icon="times"
                                    @click="removeFilter(item.index)"
                                    v-tooltip="$t('labels.filters.remove')"
                                ></v-button>
                            </v-table-cell>
                        </template>
                    </v-table>

                    <v-toolbar>
                        <v-button icon="plus" @click="addFilter">{{ $t("labels.filters.add") }}</v-button>
                        <v-spacer></v-spacer>
                        <v-message icon="question-circle" type="note">
                            {{ $t("messages.labels.emptyDiscarded") }}
                        </v-message>
                    </v-toolbar>
                </v-form>
            </template>
        </template>
    </v-content>
</template>

<script lang="ts">
import cloneDeep from "lodash/cloneDeep";
import merge from "lodash/merge";
import {VueConstructor} from "vue";
import {Component, Mixins, Watch} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

import {ModalConfirmParams} from "@/src/components/modal/confirm.vue";
import {conflictCustomValidity} from "@/src/helpers/api";
import {beforeRoute} from "@/src/helpers/route";
import {CustomMixins} from "@/src/mixins";
import {ModalProviderFilterParams} from "@/src/views/admin/components/modal/provider-filter.vue";

import {namePattern} from "..";

const defaultProvider: Provider = {
    id: "",
    name: "",
    connector: {
        type: "",
        settings: {},
    },
    filters: [],
};

const defaultSection = "general";

const providers: Array<SelectOption> = [
    {label: "KairosDB", value: "kairosdb"},
    {label: "Prometheus", value: "prometheus"},
    {label: "RRDTool", value: "rrdtool"},
];

@Component({
    beforeRouteLeave: beforeRoute,
    beforeRouteUpdate: beforeRoute,
})
export default class Edit extends Mixins<CustomMixins>(CustomMixins) {
    public conflictCustomValidity!: (value: string) => Promise<string>;

    public form: VueConstructor | null = null;

    public formFailed = false;

    public formValidity = false;

    public loading = true;

    public namePattern = namePattern;

    public provider: Provider | null = null;

    public providers: Array<SelectOption> = providers;

    public saving = false;

    public section = defaultSection;

    public testing = false;

    public created(): void {
        this.conflictCustomValidity = conflictCustomValidity(this, "providers", this.params.id);
    }

    public mounted(): void {
        this.reset(true);
    }

    public addFilter(): void {
        this.$components.modal(
            "provider-filter",
            {
                rule: {},
            } as ModalProviderFilterParams,
            (value: FilterRule) => {
                if (value && this.provider) {
                    if (!this.provider.filters) {
                        this.provider.filters = [];
                    }
                    this.provider.filters.push(value);
                    this.emitUpdate();
                }
            },
        );
    }

    public deleteProvider(apply = false): void {
        if (this.provider === null) {
            return;
        }

        if (apply) {
            this.$http
                .delete(`/api/v1/providers/${this.provider.id}`)
                .then(response => response.json())
                .then(() => {
                    this.unguard();
                    this.$components.notify(this.$tc("messages.providers.deleted", 1) as string, "success");
                    this.$router.push({name: "admin-providers-list"});
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
                message: this.$tc(`messages.${this.params.type}.delete`, 1, this.provider),
            } as ModalConfirmParams,
            (value: boolean) => {
                if (value) {
                    this.deleteProvider(true);
                }
            },
        );
    }

    public get edit(): boolean {
        return this.params.id !== "new";
    }

    public editFilter(index: number): void {
        if (this.provider === null || !this.provider.filters) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        }

        const filters: Array<FilterRule> = this.provider.filters;

        this.$components.modal(
            "provider-filter",
            {
                rule: cloneDeep(filters[index]),
            } as ModalProviderFilterParams,
            (value: FilterRule) => {
                if (value) {
                    filters.splice(index, 1, value);
                }
            },
        );
    }

    @Watch("provider.connector.type", {immediate: true})
    public async onConnectorType(to: string, from?: string): Promise<void> {
        if (!to) {
            return;
        }

        // Reset settings on type change (check for from to avoid emptying upon
        // load)
        if (this.provider?.connector && from) {
            this.provider.connector.settings = {};
        }

        try {
            Object.assign(this, {
                form: (await import(`../components/provider/${to}.vue`)).default,
                formFailed: false,
            });
        } catch {
            this.formFailed = true;
        }
    }

    @Watch("$route.hash", {immediate: true})
    public onRouteHash(to: string): void {
        this.section = to ? to.substr(1) : defaultSection;
    }

    public onValidity(to: boolean): void {
        this.formValidity = to;
    }

    public removeFilter(index: number): void {
        if (this.provider?.filters) {
            this.provider.filters.splice(index, 1);
            this.emitUpdate();
        }
    }

    public reset(apply = false): void {
        if (!apply) {
            this.$components.modal(
                "confirm",
                {
                    button: {
                        label: this.$t("labels.providers.reset"),
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

        this.formValidity = false;

        if (!this.edit) {
            this.provider = cloneDeep(defaultProvider);
            this.loading = false;
            this.emitUpdate();
            this.guardWatch("provider");

            return;
        }

        this.loading = true;

        this.$http
            .get(`/api/v1/providers/${this.params.id}`)
            .then(response => response.json())
            .then(
                (response: APIResponse<Provider>) => {
                    this.provider = merge({}, defaultProvider, response.data);
                    this.loading = false;
                    this.emitUpdate();
                    this.guardWatch("provider");
                },
                this.handleError(() => {
                    this.emitUpdate();
                    this.loading = false;
                }, true),
            );
    }

    public save(go: boolean): void {
        if (!this.guarded) {
            this.redirect(go);
            return;
        }

        if (this.provider === null) {
            this.$components.notify(this.$t("messages.error.unhandled") as string, "error");
            return;
        } else if (this.saving) {
            return;
        }
        this.saving = true;

        const provider: Provider = cloneDeep(this.provider);
        if (!this.edit) {
            provider.enabled = true;
        }

        this.$http({
            url: `/api/v1/providers${this.edit ? `/${this.params.id}` : ""}`,
            method: this.edit ? "PUT" : "POST",
            body: provider,
        }).then(
            () => {
                this.unguard();
                this.$components.notify(this.$t("messages.providers.saved") as string, "success");
                this.redirect(go);
                this.saving = false;
            },
            this.handleError(() => {
                this.saving = false;
            }),
        );
    }

    public test(): void {
        this.testing = true;

        this.$http
            .post(`/api/v1/providers/test`, this.provider)
            .then(response => response.json())
            .then(
                (response: APIResponse<unknown>) => {
                    if (response.error) {
                        this.$components.notify(
                            this.$t("messages.providers.test.error", [response.error]) as string,
                            "error",
                        );
                    } else {
                        this.$components.notify(this.$t("messages.providers.test.success") as string, "success");
                    }

                    this.testing = false;
                },
                this.handleError(() => {
                    this.testing = false;
                }),
            );
    }

    public get validity(): boolean {
        const validity: Record<string, boolean> = {
            general: this.formValidity,
        };

        this.$parent.$emit("provider-validity", validity);

        return !Object.values(validity).includes(false);
    }

    private emitUpdate(): void {
        this.$parent.$emit("provider-updated", this.provider?.filters ? this.provider.filters.length : null);
    }

    private redirect(go: boolean): void {
        this.$router.push(
            go && this.provider?.name
                ? {
                      name: "admin-metrics-list",
                      query: {filter: `{__provider__=${JSON.stringify(this.provider.name)}}`},
                  }
                : {name: "admin-providers-list"},
        );
    }
}
</script>

<style lang="scss" scoped>
@import "../mixins";

.v-content {
    ::v-deep {
        @include content;
    }

    .empty {
        color: var(--light-gray);
        text-transform: lowercase;
    }

    .v-table .v-icon {
        color: var(--light-gray);
        margin: 0 0.75rem;
    }

    .v-toolbar .v-label {
        color: var(--light-gray);
    }
}
</style>
