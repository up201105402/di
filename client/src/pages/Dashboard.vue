<template>
    <div class="wrapper">
        <side-bar :sidebarLinks="sidebarLinks">
        </side-bar>
        <div class="main-panel">
            <top-navbar></top-navbar>
            <div class="content">
                <div>
                    <div class="row">
                        <div class="col-12">
                            <card type="chart">
                                <div class="row" slot="header">
                                    <div class="col-sm-6 text-left">
                                        <h5 class="card-category">{{ $t('pages.home.name') }}</h5>
                                        <h2 class="card-title">{{ $t('pages.home.name') }}</h2>
                                    </div>
                                    <div class="col-sm-6">
                                        <div class="btn-group btn-group-toggle float-right" data-toggle="buttons">
                                            <label v-for="(option, index) in bigLineChartCategories" :key="option"
                                                class="btn btn-sm btn-primary btn-simple"
                                                :class="{ active: bigLineChart.activeIndex === index }" :id="index">
                                                <input type="radio" @click="initBigChart(index)" name="options"
                                                    autocomplete="off" :checked="bigLineChart.activeIndex === index">
                                                {{ option }}
                                            </label>
                                        </div>
                                    </div>
                                </div>
                                <div class="chart-area">
                                    <line-chart style="height: 100%" ref="bigChart" chart-id="big-line-chart"
                                        :chartId="'big-line-chart'" :chart-id="'big-line-chart'"
                                        :data="bigLineChart.chartData" :chart-data="bigLineChart.chartData"
                                        :gradient-colors="bigLineChart.gradientColors"
                                        :gradient-stops="bigLineChart.gradientStops"
                                        :extra-options="bigLineChart.extraOptions">
                                    </line-chart>
                                </div>
                            </card>
                        </div>
                    </div>
                    <div class="row">
                        <div class="col-lg-4" :class="{ 'text-right': false }">
                            <card type="chart">
                                <div slot="header">
                                    <h5 class="card-category">{{ $t('pages.dashboard.totalShipments') }}</h5>
                                    <h3 class="card-title"><i class="tim-icons icon-bell-55 text-primary "></i> 763,215</h3>
                                </div>
                                <div class="chart-area">
                                    <line-chart style="height: 100%" chart-id="purple-line-chart"
                                        :chart-id="'purple-line-chart'" :data="purpleLineChart.chartData"
                                        :chart-data="purpleLineChart.chartData"
                                        :gradient-colors="purpleLineChart.gradientColors"
                                        :gradient-stops="purpleLineChart.gradientStops"
                                        :extra-options="purpleLineChart.extraOptions">
                                    </line-chart>
                                </div>
                            </card>
                        </div>
                        <div class="col-lg-4" :class="{ 'text-right': false }">
                            <card type="chart">
                                <div slot="header">
                                    <h5 class="card-category">{{ $t('pages.dashboard.dailySales') }}</h5>
                                    <h3 class="card-title"><i class="tim-icons icon-delivery-fast text-info "></i> 3,500€
                                    </h3>
                                </div>
                                <div class="chart-area">
                                    <bar-chart style="height: 100%" chart-id="blue-bar-chart"
                                        :chart-data="blueBarChart.chartData" :chart-id="'blue-bar-chart'"
                                        :data="blueBarChart.chartData" :gradient-stops="blueBarChart.gradientStops"
                                        :extra-options="blueBarChart.extraOptions">
                                    </bar-chart>
                                </div>
                            </card>
                        </div>
                        <div class="col-lg-4" :class="{ 'text-right': false }">
                            <card type="chart">
                                <div slot="header">
                                    <h5 class="card-category">{{ $t('pages.dashboard.completedTasks') }}</h5>
                                    <h3 class="card-title"><i class="tim-icons icon-send text-success "></i> 12,100K</h3>
                                </div>
                                <div class="chart-area">
                                    <line-chart style="height: 100%" chart-id="green-line-chart"
                                        :chart-data="greenLineChart.chartData" :chart-id="'green-line-chart'"
                                        :data="greenLineChart.chartData" :gradient-stops="greenLineChart.gradientStops"
                                        :extra-options="greenLineChart.extraOptions">
                                    </line-chart>
                                </div>
                            </card>
                        </div>
                    </div>
                    <div class="row">
                        <!-- <div class="col-lg-6 col-md-12">
                    <card type="tasks" :header-classes="{ 'text-right': false }">
                        <template slot="header">
                            <h6 class="title d-inline">{{ $t('pages.dashboard.tasks', { count: 5 }) }}</h6>
                            <p class="card-category d-inline">{{ $t('pages.dashboard.today') }}</p>
                            <base-dropdown menu-on-right="" tag="div" title-classes="btn btn-link btn-icon"
                                aria-label="Settings menu" :class="{ 'float-left': false }">
                                <i slot="title" class="tim-icons icon-settings-gear-63"></i>
                                <a class="dropdown-item" href="#pablo">{{ $t('pages.dashboard.dropdown.action') }}</a>
                                <a class="dropdown-item" href="#pablo">{{ $t('pages.dashboard.dropdown.anotherAction') }}</a>
                                <a class="dropdown-item" href="#pablo">{{ $t('pages.dashboard.dropdown.somethingElse') }}</a>
                            </base-dropdown>
                        </template>
                        <div class="table-full-width table-responsive">
                            <task-list></task-list>
                        </div>
                    </card>
                </div> -->
                        <!-- <div class="col-lg-6 col-md-12">
                    <card class="card" :header-classes="{ 'text-right': false }">
                        <h4 slot="header" class="card-title">{{ $t('pages.dashboard.simpleTable') }}</h4>
                        <div class="table-responsive">
                            <user-table></user-table>
                        </div>
                    </card>
                </div> -->
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>
<script>
    import TopNavbar from '../components/TopNavbar.vue'
    import LineChart from '../components/charts/LineChart';
    import BarChart from '../components/charts/BarChart';
    import * as chartConfigs from '../components/charts/config';
    import TaskList from './dashboard/TaskList.vue';
    import UserTable from './dashboard/UserTable.vue';
    import config from './config';
    import { routes } from '../routes';

    export default {
        components: {
            TopNavbar,
            LineChart,
            BarChart,
            TaskList,
            UserTable
        },
        data() {
            return {
                bigLineChart: {
                    allData: [
                        [100, 70, 90, 70, 85, 60, 75, 60, 90, 80, 110, 100],
                        [80, 120, 105, 110, 95, 105, 90, 100, 80, 95, 70, 120],
                        [60, 80, 65, 130, 80, 105, 90, 130, 70, 115, 60, 130]
                    ],
                    activeIndex: 0,
                    chartData: {
                        datasets: [{}],
                        labels: ['JAN', 'FEB', 'MAR', 'APR', 'MAY', 'JUN', 'JUL', 'AUG', 'SEP', 'OCT', 'NOV', 'DEC'],
                    },
                    extraOptions: chartConfigs.purpleChartOptions,
                    gradientColors: config.colors.primaryGradient,
                    gradientStops: [1, 0.4, 0],
                    categories: []
                },
                purpleLineChart: {
                    extraOptions: chartConfigs.purpleChartOptions,
                    chartData: {
                        labels: ['JUL', 'AUG', 'SEP', 'OCT', 'NOV', 'DEC'],
                        datasets: [{
                            label: "Data",
                            fill: true,
                            borderColor: config.colors.primary,
                            borderWidth: 2,
                            borderDash: [],
                            borderDashOffset: 0.0,
                            pointBackgroundColor: config.colors.primary,
                            pointBorderColor: 'rgba(255,255,255,0)',
                            pointHoverBackgroundColor: config.colors.primary,
                            pointBorderWidth: 20,
                            pointHoverRadius: 4,
                            pointHoverBorderWidth: 15,
                            pointRadius: 4,
                            data: [80, 100, 70, 80, 120, 80],
                        }]
                    },
                    gradientColors: config.colors.primaryGradient,
                    gradientStops: [1, 0.2, 0],
                },
                greenLineChart: {
                    extraOptions: chartConfigs.greenChartOptions,
                    chartData: {
                        labels: ['JUL', 'AUG', 'SEP', 'OCT', 'NOV'],
                        datasets: [{
                            label: "My First dataset",
                            fill: true,
                            borderColor: config.colors.danger,
                            borderWidth: 2,
                            borderDash: [],
                            borderDashOffset: 0.0,
                            pointBackgroundColor: config.colors.danger,
                            pointBorderColor: 'rgba(255,255,255,0)',
                            pointHoverBackgroundColor: config.colors.danger,
                            pointBorderWidth: 20,
                            pointHoverRadius: 4,
                            pointHoverBorderWidth: 15,
                            pointRadius: 4,
                            data: [90, 27, 60, 12, 80],
                        }]
                    },
                    gradientColors: ['rgba(66,134,121,0.15)', 'rgba(66,134,121,0.0)', 'rgba(66,134,121,0)'],
                    gradientStops: [1, 0.4, 0],
                },
                blueBarChart: {
                    extraOptions: chartConfigs.barChartOptions,
                    chartData: {
                        labels: ['USA', 'GER', 'AUS', 'UK', 'RO', 'BR'],
                        datasets: [{
                            label: "Countries",
                            fill: true,
                            borderColor: config.colors.info,
                            borderWidth: 2,
                            borderDash: [],
                            borderDashOffset: 0.0,
                            data: [53, 20, 10, 80, 100, 45],
                        }]
                    },
                    gradientColors: config.colors.primaryGradient,
                    gradientStops: [1, 0.4, 0],
                },
                sidebarLinks: routes
            }
        },
        computed: {
            bigLineChartCategories() {
                return [
                    "Accounts",
                    "Purchases",
                    "Sessions"
                ];
            }
        },
        methods: {
            initBigChart(index) {
                let chartData = {
                    datasets: [{
                        fill: true,
                        borderColor: config.colors.primary,
                        borderWidth: 2,
                        borderDash: [],
                        borderDashOffset: 0.0,
                        pointBackgroundColor: config.colors.primary,
                        pointBorderColor: 'rgba(255,255,255,0)',
                        pointHoverBackgroundColor: config.colors.primary,
                        pointBorderWidth: 20,
                        pointHoverRadius: 4,
                        pointHoverBorderWidth: 15,
                        pointRadius: 4,
                        data: this.bigLineChart.allData[index]
                    }],
                    labels: ['JAN', 'FEB', 'MAR', 'APR', 'MAY', 'JUN', 'JUL', 'AUG', 'SEP', 'OCT', 'NOV', 'DEC'],
                }
                this.bigLineChart.chartData = chartData;
                this.bigLineChart.activeIndex = index;
            }
        },
        mounted() {
            console.error('dashboard mounted');
            this.i18n = this.$i18n;
            this.initBigChart(0);
        },
        setup() {
            // userAuth();
        }
    };
</script>
<style></style>