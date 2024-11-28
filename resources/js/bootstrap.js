import "flyonui/flyonui";
import "./apexcharts-helper.js";
import "../../node_modules/lodash/lodash.min.js";
import axios from "axios";
import Alpine from "alpinejs";
import ApexCharts from "apexcharts";

window.ApexCharts = ApexCharts;
window.axios = axios;
window.axios.defaults.headers.common["X-Requested-With"] = "XMLHttpRequest";

window.Alpine = Alpine;

Alpine.start();
