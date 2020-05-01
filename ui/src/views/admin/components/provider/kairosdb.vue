<template>
    <div>
        <v-label>{{ $t("labels.url") }}</v-label>
        <v-input
            required
            type="url"
            :help="$t('help.providers.url', ['KairosDB'])"
            :placeholder="$t('labels.placeholders.example', ['http://localhost:8080/'])"
            v-model="settings.url"
        ></v-input>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";
@Component
export default class ProviderKairosDBComponent extends Vue {
    @Prop({required: true, type: Object})
    public settings!: Record<string, unknown>;
}
</script>
