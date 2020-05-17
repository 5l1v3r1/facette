<template>
    <v-content>
        <v-toolbar clip="content">
            <v-spacer></v-spacer>
            <v-button href="https://facette.io" icon="globe" target="_blank">{{ $t("labels.visit.website") }}</v-button>
        </v-toolbar>

        <v-spinner v-if="loading"></v-spinner>

        <v-message-error @retry="getVersion" v-else-if="erred"></v-message-error>

        <template v-else>
            <h1>{{ $t("labels.info._") }}</h1>

            <dl>
                <template v-if="version">
                    <dt>{{ $t("labels.info.version") }}</dt>
                    <dd>{{ version.version }}</dd>

                    <dt>{{ $t("labels.info.branch") }}</dt>
                    <dd>{{ version.branch }}</dd>

                    <dt>{{ $t("labels.info.revision") }}</dt>
                    <dd>{{ version.revision }}</dd>

                    <dt>{{ $t("labels.info.compiler") }}</dt>
                    <dd>{{ version.compiler }}</dd>

                    <dt>{{ $t("labels.info.buildDate") }}</dt>
                    <dd>{{ version.buildDate }}</dd>
                </template>

                <dt>{{ $t("labels.info.connectors") }}</dt>
                <template v-if="connectors">
                    <dd :key="name" v-for="name in connectors">{{ name }}</dd>
                </template>
                <dd v-else>{{ $t("messages.notAvailable") }}</dd>
            </dl>
        </template>
    </v-content>
</template>

<script lang="ts">
import {Component, Mixins} from "vue-property-decorator";

import {CustomMixins} from "@/src/mixins";

interface Version {
    version?: string;
    branch?: string;
    revision?: string;
    compiler?: string;
    buildDate?: string;
}

@Component
export default class Info extends Mixins<CustomMixins>(CustomMixins) {
    public loading = true;

    public version: Version | null = null;

    public mounted(): void {
        this.getVersion();
    }

    public get connectors(): Array<string> {
        return this.$store.getters.connectors;
    }

    private getVersion(): void {
        this.$http
            .get("/api/v1/version")
            .then(response => response.json())
            .then(
                (response: APIResponse<Version>) => {
                    this.version = response.data ?? null;
                    this.loading = false;
                },
                this.handleError(() => {
                    this.loading = false;
                }, true),
            );
    }
}
</script>

<style lang="scss" scoped>
.v-content {
    dl {
        dt {
            color: var(--light-gray);
            margin: 0.65rem 0;
        }

        dd {
            margin: 0 0 0 1.5rem;
        }
    }
}
</style>
