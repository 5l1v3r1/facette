<template>
    <v-app>
        <v-modal-help></v-modal-help>
        <v-notifier></v-notifier>

        <router-view name="toolbar"></router-view>
        <router-view name="sidebar"></router-view>
        <router-view></router-view>
    </v-app>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator";

@Component
export default class Main extends Vue {
    public mounted(): void {
        document.addEventListener("keydown", this.onModifier);
        document.addEventListener("keyup", this.onModifier);
        document.addEventListener("visibilitychange", this.onModifier);
    }

    public beforeDestroy(): void {
        document.removeEventListener("keydown", this.onModifier);
        document.removeEventListener("keyup", this.onModifier);
        document.removeEventListener("visibilitychange", this.onModifier);
    }

    private onModifier(e: Event | KeyboardEvent): void {
        switch (e.type) {
            case "keydown":
            case "keyup": {
                const ke = e as KeyboardEvent;

                if (ke.code === "AltLeft" || ke.code === "AltRight") {
                    const modifiers = Object.assign(this.$store.getters.modifiers, {
                        alt: e.type === "keydown",
                    });

                    this.$store.commit("modifiers", modifiers);
                } else if (ke.code === "ShiftLeft" || ke.code === "ShiftRight") {
                    const modifiers = Object.assign(this.$store.getters.modifiers, {
                        shift: e.type === "keydown",
                    });

                    this.$store.commit("modifiers", modifiers);
                }

                break;
            }

            case "visibilitychange": {
                this.$store.commit("modifiers", {alt: false, shift: false});
                break;
            }
        }
    }
}
</script>
