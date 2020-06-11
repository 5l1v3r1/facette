<template>
    <v-content>
        <v-modal-confirm></v-modal-confirm>

        <v-toolbar clip="content">
            <v-button icon="save" @click="save">{{ $t("labels.settings.save") }}</v-button>

            <v-divider vertical></v-divider>

            <v-button icon="undo" :disabled="!guarded" @click="reset()">{{ $t("labels.reset") }}</v-button>
        </v-toolbar>

        <template v-if="section === 'display'">
            <h1>{{ $t("labels.settings.display") }}</h1>

            <v-form>
                <v-label>{{ $t("labels.language._") }}</v-label>
                <v-select
                    :options="locales"
                    :placeholder="$t('labels.language.select')"
                    v-autofocus
                    v-model="settings.locale"
                ></v-select>

                <v-label>{{ $t("labels.theme._") }}</v-label>
                <v-select
                    :options="themes"
                    :placeholder="$t('labels.theme.select')"
                    v-model="settings.theme"
                ></v-select>

                <v-label>{{ $t("labels.timezone._") }}</v-label>
                <v-select
                    :options="timezones"
                    :placeholder="$t('labels.timezone.select')"
                    v-model="settings.timezoneUTC"
                ></v-select>
            </v-form>
        </template>

        <template v-else-if="section === 'keyboard'">
            <h1>{{ $t("labels.keyboard.shortcuts._") }}</h1>

            <v-form>
                <v-markdown @click.native="onHelpClick">
                    {{ $t("help.keyboard.shortcuts", ["#keyboard"]) }}
                </v-markdown>

                <v-checkbox toggle v-model="settings.shortcuts">
                    {{ $t("labels.keyboard.shortcuts.enable") }}
                </v-checkbox>
            </v-form>
        </template>

        <v-message type="warning" v-if="guarded">{{ $t("messages.settings.reload") }}</v-message>
    </v-content>
</template>

<script lang="ts">
import {Component, Mixins, Watch} from "vue-property-decorator";

import {SelectOption} from "@/types/components";

import {ModalConfirmParams} from "@/src/components/modal/confirm.vue";
import themes from "@/src/components/vue/app/themes";
import {beforeRoute} from "@/src/helpers/route";
import {resolveOption} from "@/src/helpers/select";
import {CustomMixins} from "@/src/mixins";

interface Settings {
    locale: string;
    shortcuts: true;
    theme: string;
    timezoneUTC: boolean;
}

const defaultSection = "display";

const timezones: Array<SelectOption> = [
    {label: "labels.timezone.local", value: false},
    {label: "labels.timezone.utc", value: true},
];

@Component({
    beforeRouteLeave: beforeRoute,
    beforeRouteUpdate: beforeRoute,
})
export default class Display extends Mixins<CustomMixins>(CustomMixins) {
    public section = defaultSection;

    public settings: Settings = {
        locale: "",
        shortcuts: true,
        theme: "",
        timezoneUTC: false,
    };

    public mounted(): void {
        this.reset(true);
    }

    public get locales(): Array<SelectOption> {
        return Object.keys(this.$i18n.messages)
            .sort()
            .map(name => ({label: this.$t("name", name) as string, value: name}));
    }

    public onHelpClick(e: MouseEvent): void {
        e.preventDefault();

        if ((e.target as HTMLElement).tagName === "A") {
            this.$components.modal("help", {tab: "keyboard"});
        }
    }

    @Watch("$route.hash", {immediate: true})
    public onRouteHash(to: string): void {
        this.section = to ? to.substr(1) : defaultSection;
    }

    public reset(apply = false): void {
        if (!apply) {
            this.$components.modal(
                "confirm",
                {
                    button: {
                        label: this.$t("labels.settings.reset"),
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

        this.settings = {
            locale: this.$store.getters.locale,
            shortcuts: this.$store.getters.shortcuts,
            theme: this.$store.getters.theme,
            timezoneUTC: this.$store.getters.timezoneUTC,
        };

        this.guardWatch("settings");
    }

    public save(): void {
        this.$store.commit("locale", this.settings.locale);
        this.$store.commit("shortcuts", this.settings.shortcuts);
        this.$store.commit("theme", this.settings.theme);
        this.$store.commit("timezoneUTC", this.settings.timezoneUTC);

        this.$components.notify(
            this.$t("messages.settings.saved", this.settings.locale as string) as string,
            "success",
        );

        this.unguard();
        location.reload();
    }

    public get themes(): Array<SelectOption> {
        return Object.keys(themes).map(key => ({label: themes[key].name, value: key}));
    }

    public get timezones(): Array<SelectOption> {
        return timezones.map(option => resolveOption(this, option));
    }
}
</script>

<style lang="scss" scoped>
.v-content {
    .v-form,
    .v-message {
        width: 33.33%;
    }

    .v-label {
        margin-top: 1rem;
    }

    .v-select {
        width: 50%;
    }
}
</style>
