<template>
    <div class="v-formblock">
        <v-checkbox toggle v-model="axis.show">{{ $t("labels.show") }}</v-checkbox>

        <div class="reveal" :class="{visible: axis.show}">
            <v-label>{{ $tc("labels.labels", 1) }}</v-label>
            <v-input :delay="500" :help="$t('help.charts.axes.label')" v-model="axis.label"></v-input>

            <v-flex>
                <v-flex direction="column">
                    <v-label>{{ $t("labels.charts.axes.min") }}</v-label>
                    <v-input
                        type="number"
                        :delay="500"
                        :help="$t('help.charts.axes.min')"
                        v-model.number="axis.min"
                    ></v-input>
                </v-flex>

                <v-flex direction="column">
                    <v-label>{{ $t("labels.charts.axes.max") }}</v-label>
                    <v-input
                        type="number"
                        :delay="500"
                        :help="$t('help.charts.axes.max')"
                        v-model.number="axis.max"
                    ></v-input>
                </v-flex>
            </v-flex>
        </div>
    </div>
</template>

<script lang="ts">
import {Component, Prop, Vue} from "vue-property-decorator";

@Component
export default class FormChartYAxisComponent extends Vue {
    @Prop({required: true, type: Object})
    public axis!: ChartYAxis;
}
</script>

<style lang="scss" scoped>
.v-formblock {
    .reveal {
        visibility: hidden;

        &.visible {
            visibility: visible;
        }
    }

    .v-checkbox:first-child {
        margin-top: 0 !important;
    }

    .v-input.number {
        width: 8rem !important;
    }
}
</style>
