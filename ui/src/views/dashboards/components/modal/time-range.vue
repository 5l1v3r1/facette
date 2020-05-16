<template>
    <v-modal name="time-range">
        <v-form slot-scope="modal">
            <v-datetime :label="$t('labels.timeRange.from')" v-autofocus.select v-model="timeRange.from"></v-datetime>

            <v-datetime :label="$t('labels.timeRange.to')" :min="timeRange.from" v-model="timeRange.to"></v-datetime>

            <template slot="bottom">
                <v-button @click="modal.close(false)">{{ $t("labels.cancel") }}</v-button>
                <v-button primary :disabled="!timeRange.from || !timeRange.to" @click="modal.close(timeRange)">
                    {{ $t("labels.timeRange.set") }}
                </v-button>
            </template>
        </v-form>
    </v-modal>
</template>

<script lang="ts">
import {Component, Vue} from "vue-property-decorator";

interface TimeRangeParams {
    timeRange: TimeRange;
}

@Component
export default class ModalTimeRangeComponent extends Vue {
    public timeRange: TimeRange = {
        from: "",
        to: "",
    };

    public onModalShow(params: TimeRangeParams): void {
        this.timeRange = params.timeRange;
    }
}
</script>

<style lang="scss" scoped>
.v-modal {
    .v-datetime {
        height: 100%;
    }

    .v-datetime + .v-datetime {
        margin-left: 1.5rem;
    }
}
</style>
