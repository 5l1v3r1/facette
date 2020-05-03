<template>
    <div>
        <v-label>{{ $t("labels.url") }}</v-label>
        <v-input
            required
            type="url"
            :help="$t('help.providers.url', ['Prometheus'])"
            :placeholder="$t('labels.placeholders.example', ['http://localhost:9090/'])"
            v-model="settings.url"
        ></v-input>

        <v-checkbox toggle v-model="settings.skipVerify" v-if="secured">{{ $t("labels.tls.skipVerify") }}</v-checkbox>

        <v-label>{{ $tc("labels.filters._", 1) }}</v-label>
        <v-input
            :help="$t('help.providers.prometheus.filter')"
            :placeholder="$t('labels.placeholders.example', ['{job=&quot;prometheus&quot;}'])"
            v-model="settings.filter"
        ></v-input>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class ProviderPrometheusComponent extends Vue {
    @Prop({required: true, type: Object})
    public settings!: Record<string, unknown>;

    public get secured(): boolean {
        return Boolean(((this.settings.url ?? "") as string).startsWith("https://"));
    }
}
</script>
