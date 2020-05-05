<template>
    <v-content>
        <v-toolbar clip="content">
            <v-button icon="save" @click="save">{{ $t("labels.settings.save") }}</v-button>

            <v-divider vertical></v-divider>

            <v-button icon="undo" :disabled="!guarded" @click="reset">{{ $t("labels.reset") }}</v-button>
        </v-toolbar>

        <h1>{{ $t("labels.settings.display.alt") }}</h1>

        <v-form>
            <v-label>{{ $t("labels.language._") }}</v-label>
            <v-select
                :options="locales"
                :placeholder="$t('labels.language.select')"
                v-autofocus
                v-model="settings.locale"
            ></v-select>

            <v-label>{{ $t("labels.theme._") }}</v-label>
            <v-select :options="themes" :placeholder="$t('labels.theme.select')" v-model="settings.theme"></v-select>

            <v-label>{{ $t("labels.timezone._") }}</v-label>
            <v-select
                :options="timezones"
                :placeholder="$t('labels.timezone.select')"
                v-model="settings.timezoneUTC"
            ></v-select>
        </v-form>
    </v-content>
</template>

<script lang="ts">
import {Component, Mixins} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

import themes from "@/src/components/vue/app/themes";
import {resolveOption} from "@/src/helpers/select";
import {CustomMixins} from "@/src/mixins";

interface Settings {
    locale: string;
    theme: string;
    timezoneUTC: boolean;
}

const timezones: Array<SelectOption> = [
    {label: "labels.timezone.local", value: false},
    {label: "labels.timezone.utc", value: true},
];

@Component
export default class Display extends Mixins<CustomMixins>(CustomMixins) {
    public settings: Settings = {
        locale: "",
        theme: "",
        timezoneUTC: false,
    };

    public mounted(): void {
        this.reset();
    }

    public get locales(): Array<SelectOption> {
        return Object.keys(this.$i18n.messages)
            .sort()
            .map((name: string) => ({label: this.$t("name", name) as string, value: name}));
    }

    public reset(): void {
        this.settings = {
            locale: this.$store.getters.locale,
            theme: this.$store.getters.theme,
            timezoneUTC: this.$store.getters.timezoneUTC,
        };

        this.guardWatch("settings");
    }

    public save(): void {
        this.$store.commit("locale", this.settings.locale);
        this.$store.commit("theme", this.settings.theme);
        this.$store.commit("timezoneUTC", this.settings.timezoneUTC);

        this.$components.notify(
            this.$t("messages.settings.saved", this.settings.locale as string) as string,
            "success",
        );

        this.reset();
    }

    public get themes(): Array<SelectOption> {
        return Object.keys(themes).map((key: string) => ({label: themes[key].name, value: key}));
    }

    public get timezones(): Array<SelectOption> {
        return timezones.map((option: SelectOption) => resolveOption(this, option));
    }
}
</script>

<style lang="scss" scoped>
.v-content {
    .v-form {
        width: 25%;
    }

    .v-label {
        margin-top: 1rem;
    }

    .v-select {
        width: 100%;
    }
}
</style>
