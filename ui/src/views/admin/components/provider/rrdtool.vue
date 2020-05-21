<template>
    <div>
        <v-label>{{ $t("labels.providers.rrdtool.path") }}</v-label>
        <v-input
            required
            :help="$t('help.providers.rrdtool.path')"
            :placeholder="$t('labels.placeholders.example', ['/var/lib/collectd/rrd'])"
            v-model="settings.path"
        ></v-input>

        <v-label>{{ $t("labels.providers.rrdtool.pattern") }}</v-label>
        <v-input
            required
            :help="$t('help.providers.rrdtool.pattern', ['https://github.com/google/re2/wiki/Syntax'])"
            :placeholder="$t('labels.placeholders.example', ['(?P<instance>[^/]+)/(?P<__name__>.+)\\.rrd'])"
            v-model="settings.pattern"
        ></v-input>

        <v-label>{{ $t("labels.providers.rrdtool.daemon") }}</v-label>
        <v-input
            required
            :help="$t('help.providers.rrdtool.daemon')"
            :placeholder="$t('labels.placeholders.example', ['unix:/var/run/rrdcached.sock'])"
            v-model="settings.daemon"
        ></v-input>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class ProviderRRDToolComponent extends Vue {
    @Prop({required: true, type: Object})
    public settings!: Record<string, unknown>;
}
</script>
