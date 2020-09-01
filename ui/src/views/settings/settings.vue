<template>
    <v-content>
        <v-toolbar clip="content">
            <v-button icon="save" @click="save">{{ i18n.t("labels.settings.save") }}</v-button>

            <v-divider vertical></v-divider>

            <v-button icon="undo" :disabled="!routeGuarded" @click="reset()">{{ i18n.t("labels.reset") }}</v-button>
        </v-toolbar>

        <template v-if="$route.params.section === 'display'">
            <h1>{{ i18n.t("labels.settings.display") }}</h1>

            <v-form class="half">
                <v-label>{{ i18n.t("labels.language._") }}</v-label>
                <v-select
                    class="half"
                    :options="localeOptions"
                    :placeholder="i18n.t('labels.language.select')"
                    v-model:value="settings.locale"
                ></v-select>

                <v-label>{{ i18n.t("labels.theme._") }}</v-label>
                <v-select
                    class="half"
                    :options="themeOptions"
                    :placeholder="i18n.t('labels.theme.select')"
                    v-model:value="settings.theme"
                ></v-select>

                <v-label>{{ i18n.t("labels.timezone._") }}</v-label>
                <v-select
                    class="half"
                    :options="timezoneOptions"
                    :placeholder="i18n.t('labels.timezone.select')"
                    v-model:value="settings.timezoneUTC"
                ></v-select>
            </v-form>
        </template>

        <template v-else-if="$route.params.section === 'keyboard'">
            <h1>{{ i18n.t("labels.settings.keyboard") }}</h1>

            <v-form class="third">
                <v-markdown
                    :content="i18n.t('help.keyboard.shortcuts', ['#keyboard'])"
                    @click="onHelpClick"
                ></v-markdown>

                <v-checkbox type="toggle" v-model:value="settings.shortcuts">
                    {{ i18n.t("labels.keyboard.shortcuts.enable") }}
                </v-checkbox>
            </v-form>
        </template>

        <v-message class="half" type="warning" v-if="routeGuarded">{{ i18n.t("messages.settings.reload") }}</v-message>
    </v-content>
</template>

<script lang="ts">
import clone from "lodash/clone";
import {computed, onBeforeMount, onMounted, ref} from "vue";
import {useI18n} from "vue-i18n";
import {onBeforeRouteLeave, onBeforeRouteUpdate} from "vue-router";
import {useStore} from "vuex";

import {SelectOption} from "types/ui";

import common from "@/common";
import {ModalConfirmParams} from "@/components/modal/confirm.vue";
import {useUI} from "@/components/ui";
import themes from "@/components/ui/themes";
import {State} from "@/store";

interface Settings {
    locale: string;
    shortcuts: boolean;
    theme: string | null;
    timezoneUTC: boolean;
}

const defaultSettings = {
    locale: "",
    shortcuts: true,
    theme: null,
    timezoneUTC: false,
};

export default {
    setup(): Record<string, unknown> {
        const i18n = useI18n();
        const store = useStore<State>();
        const ui = useUI();

        const {applyRouteParams, beforeRoute, routeGuarded, unwatchGuard, watchGuard} = common;

        const settings = ref<Settings>(clone(defaultSettings));

        const localeOptions = computed<Array<SelectOption>>(() => {
            return Object.keys(i18n.messages.value as Record<string, unknown>)
                .sort()
                .map(name => ({label: i18n.t("name", name), value: name}));
        });

        const themeOptions = computed<Array<SelectOption>>(() => {
            return Array<SelectOption>().concat(
                {label: i18n.t("labels.theme.auto"), value: null},
                ...Object.keys(themes).map(key => ({label: themes[key].name, value: key})),
            );
        });

        const timezoneOptions = ref<Array<SelectOption>>([
            {label: i18n.t("labels.timezone.local"), value: false},
            {label: i18n.t("labels.timezone.utc"), value: true},
        ]);

        const onHelpClick = (ev: MouseEvent): void => {
            if ((ev.target as HTMLElement).tagName === "A") {
                ev.preventDefault();
                ui.modal("help");
            }
        };

        const reset = async (force = false): Promise<void> => {
            if (!force) {
                const ok = await ui.modal<boolean>("confirm", {
                    button: {
                        label: i18n.t("labels.settings.reset"),
                        danger: true,
                    },
                    message: i18n.t("messages.unsavedLost"),
                } as ModalConfirmParams);

                if (!ok) {
                    return;
                }
            }

            settings.value = {
                locale: store.state.locale,
                shortcuts: store.state.shortcuts,
                theme: store.state.theme,
                timezoneUTC: store.state.timezoneUTC,
            };

            watchGuard(settings);
        };

        const save = (): void => {
            store.commit("locale", settings.value.locale);
            store.commit("shortcuts", settings.value.shortcuts);
            store.commit("theme", settings.value.theme);
            store.commit("timezoneUTC", settings.value.timezoneUTC);

            store.commit("pendingNotification", {
                text: i18n.t("messages.settings.saved", null, {locale: settings.value.locale}),
                type: "success",
            });

            unwatchGuard();
            location.reload();
        };

        onBeforeRouteLeave(beforeRoute);

        onBeforeRouteUpdate(beforeRoute);

        onBeforeMount(() => applyRouteParams());

        onMounted(() => reset(true));

        return {
            i18n,
            localeOptions,
            onHelpClick,
            reset,
            routeGuarded,
            save,
            settings,
            themeOptions,
            timezoneOptions,
        };
    },
};
</script>

<style lang="scss" scoped>
@import "../mixins";

.v-content {
    @include content;

    .v-message.warning {
        margin-top: 3em;
    }
}
</style>
