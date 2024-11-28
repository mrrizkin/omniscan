import "flyonui/flyonui";
import "../../node_modules/apexcharts/dist/apexcharts.css";
import "../../node_modules/apexcharts/dist/apexcharts.min.js";
import "./apexcharts-helper.js";
import "../../node_modules/lodash/lodash.min.js";
import axios from "axios";
import Alpine from "alpinejs";

window.axios = axios;
window.axios.defaults.headers.common["X-Requested-With"] = "XMLHttpRequest";

window.Alpine = Alpine;

Alpine.start();
