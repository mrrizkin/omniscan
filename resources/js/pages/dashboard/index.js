import "apexcharts";
import { buildChart, buildTooltip } from "../../charts/apexcharts-helper";
import { onMounted } from "../../lib/utils";

onMounted(function () {
    // Apex Single Area Chart (Start)
    buildChart("#apex-single-area-chart", (mode) => ({
        chart: {
            height: 300,
            type: "area",
            toolbar: {
                show: false,
            },
            zoom: {
                enabled: false,
            },
        },
        series: [
            {
                name: "Units",
                data: [0, 100, 50, 125, 70, 150, 100, 170, 120, 175, 100, 200],
            },
        ],
        legend: {
            show: false,
        },
        dataLabels: {
            enabled: false,
        },
        stroke: {
            curve: "straight",
            width: 2,
        },
        grid: {
            strokeDashArray: 2,
            borderColor: "oklch(var(--bc) / 0.4)",
        },
        colors: ["oklch(var(--p))"], // color var
        fill: {
            type: "gradient",
            gradient: {
                shadeIntensity: 1,
                opacityFrom: 0.7,
                gradientToColors: ["oklch(var(--b1))"],
                opacityTo: 0.3,
                stops: [0, 90, 100],
            },
        },
        xaxis: {
            type: "category",
            tickPlacement: "on",
            categories: [
                "1 March 2024",
                "2 March 2024",
                "3 March 2024",
                "4 March 2024",
                "5 March 2024",
                "6 March 2024",
                "7 March 2024",
                "8 March 2024",
                "9 March 2024",
                "10 March 2024",
                "11 March 2024",
                "12 March 2024",
            ],
            axisBorder: {
                show: false,
            },
            axisTicks: {
                show: false,
            },
            tooltip: {
                enabled: false,
            },
            labels: {
                style: {
                    colors: "oklch(var(--bc) / 0.9)",
                    fontSize: "12px",
                    fontWeight: 400,
                },
                formatter: (title) => {
                    let t = title;

                    if (t) {
                        const newT = t.split(" ");
                        t = `${newT[0]} ${newT[1].slice(0, 3)}`;
                    }

                    return t;
                },
            },
        },
        yaxis: {
            labels: {
                align: "left",
                minWidth: 0,
                maxWidth: 140,
                style: {
                    colors: "oklch(var(--bc) / 0.9)",
                    fontSize: "12px",
                    fontWeight: 400,
                },
                formatter: (value) => (value >= 1000 ? `${value / 1000}k` : value),
            },
        },
        tooltip: {
            x: {
                format: "MMMM yyyy",
            },
            y: {
                formatter: (value) => `${value >= 1000 ? `${value / 1000}k` : value}`,
            },
            custom: function (props) {
                const { categories } = props.ctx.opts.xaxis;
                const { dataPointIndex } = props;
                const title = categories[dataPointIndex].split(" ");
                const newTitle = `${title[0]} ${title[1]}`;

                return buildTooltip(props, {
                    title: newTitle,
                    mode,
                    valuePrefix: "",
                    hasTextLabel: true,
                    markerExtClasses: "bg-primary",
                    wrapperExtClasses: "",
                });
            },
        },
        responsive: [
            {
                breakpoint: 568,
                options: {
                    chart: {
                        height: 300,
                    },
                    labels: {
                        style: {
                            fontSize: "10px",
                            colors: "oklch(var(--bc) / 0.9)",
                        },
                        offsetX: -2,
                        formatter: (title) => title.slice(0, 3),
                    },
                    yaxis: {
                        labels: {
                            align: "left",
                            minWidth: 0,
                            maxWidth: 140,
                            style: {
                                fontSize: "10px",
                                colors: "oklch(var(--bc) / 0.9)",
                            },
                            formatter: (value) => (value >= 1000 ? `${value / 1000}k` : value),
                        },
                    },
                },
            },
        ],
    }));

    buildChart("#lineAreaChart", (mode) => ({
        chart: {
            height: 300,
            width: "100%",
            type: "area",
            parentHeightOffset: 0,
            toolbar: {
                show: false,
            },
        },
        dataLabels: {
            enabled: false,
        },
        stroke: {
            show: false,
            curve: "straight",
        },
        legend: {
            show: true,
            position: "top",
            horizontalAlign: "start",
            labels: {
                colors: "oklch(var(--bc) / 0.9)",
                useSeriesColors: false,
            },
        },
        grid: {
            borderColor: "oklch(var(--bc) / 0.4)",
            xaxis: {
                lines: {
                    show: true,
                },
            },
        },
        colors: ["oklch(var(--su) / 0.3)", "oklch(var(--su) / 0.6)", "oklch(var(--su) / 0.9)"],
        series: [
            {
                name: "Visits",
                data: [100, 120, 90, 170, 130, 160, 140, 240, 220, 180, 270, 280, 375],
            },
            {
                name: "Clicks",
                data: [60, 80, 70, 110, 80, 100, 90, 180, 160, 140, 200, 220, 275],
            },
            {
                name: "Sales",
                data: [20, 40, 30, 70, 40, 60, 50, 140, 120, 100, 140, 180, 220],
            },
        ],
        xaxis: {
            categories: ["7/12", "8/12", "9/12", "10/12", "11/12", "12/12", "13/12", "14/12", "15/12", "16/12", "17/12", "18/12", "19/12", "20/12"],
            axisBorder: {
                show: false,
            },
            axisTicks: {
                show: false,
            },
            labels: {
                style: {
                    colors: "oklch(var(--bc) / 0.9)",
                    fontSize: "13px",
                },
            },
        },
        yaxis: {
            labels: {
                style: {
                    colors: "oklch(var(--bc) / 0.9)",
                    fontSize: "13px",
                },
            },
        },
        fill: {
            opacity: 1,
            type: "solid",
        },
        tooltip: {
            shared: false,
        },
    }));

    buildChart("#apex-column-bar-chart", (mode) => ({
        chart: {
            type: "bar",
            height: 300,
            toolbar: {
                show: false,
            },
            zoom: {
                enabled: false,
            },
        },
        series: [
            {
                name: "Investment",
                data: [25000, 39000, 65000, 45000, 79000, 80000, 69000, 63000, 60000, 66000, 90000, 78000],
            },
        ],
        plotOptions: {
            bar: {
                horizontal: false,
                columnWidth: "12px",
            },
        },
        legend: {
            show: false,
        },
        dataLabels: {
            enabled: false,
        },
        colors: ["oklch(var(--p))", "oklch(var(--b1))"],
        xaxis: {
            categories: ["Cook", "Erin", "Jack", "Will", "Gayle", "Megan", "John", "Luke", "Ellis", "Mason", "Elvis", "Liam"],
            axisBorder: {
                show: false,
            },
            axisTicks: {
                show: false,
            },
            labels: {
                style: {
                    colors: "oklch(var(--bc) / 0.9)",
                    fontSize: "12px",
                    fontWeight: 400,
                },
            },
        },
        yaxis: {
            labels: {
                align: "left",
                minWidth: 0,
                maxWidth: 140,
                style: {
                    colors: "oklch(var(--bc) / 0.9)",
                    fontSize: "12px",
                    fontWeight: 400,
                },
                formatter: (value) => (value >= 1000 ? `${value / 1000}k` : value),
            },
        },
        states: {
            hover: {
                filter: {
                    type: "darken",
                    value: 0.9,
                },
            },
        },
        tooltip: {
            y: {
                formatter: (value) => `$${value >= 1000 ? `${value / 1000}k` : value}`,
            },
            custom: function (props) {
                const { categories } = props.ctx.opts.xaxis;
                const { dataPointIndex } = props;
                const title = categories[dataPointIndex];
                const newTitle = `${title}`;

                return buildTooltip(props, {
                    title: newTitle,
                    mode,
                    hasTextLabel: true,
                    wrapperExtClasses: "min-w-28",
                    labelDivider: ":",
                    labelExtClasses: "ms-2",
                });
            },
        },
        responsive: [
            {
                breakpoint: 568,
                options: {
                    chart: {
                        height: 300,
                    },
                    plotOptions: {
                        bar: {
                            columnWidth: "10px",
                        },
                    },
                    stroke: {
                        width: 8,
                    },
                    labels: {
                        style: {
                            colors: "oklch(var(--bc) / 0.9)",
                            fontSize: "10px",
                        },
                        formatter: (title) => title.slice(0, 3),
                    },
                    yaxis: {
                        labels: {
                            align: "left",
                            minWidth: 0,
                            maxWidth: 140,
                            style: {
                                colors: "oklch(var(--bc) / 0.9)",
                                fontSize: "10px",
                            },
                            formatter: (value) => (value >= 1000 ? `${value / 1000}k` : value),
                        },
                    },
                },
            },
        ],
    }));

    buildChart("#apex-doughnut-chart", (mode) => ({
        chart: {
            height: 300,
            type: "donut",
        },
        plotOptions: {
            pie: {
                donut: {
                    size: "70%",
                    labels: {
                        show: true,
                        name: {
                            fontSize: "2rem",
                        },
                        value: {
                            fontSize: "1.5rem",
                            color: "oklch(var(--bc) / 0.9)",
                            formatter: function (val) {
                                return parseInt(val, 10) + "%";
                            },
                        },
                        total: {
                            show: true,
                            fontSize: "1rem",
                            label: "Operational",
                            formatter: function (w) {
                                return "42%";
                            },
                        },
                    },
                },
            },
        },
        series: [42, 7, 25, 25],
        labels: ["Operational", "Networking", "Hiring", "R&D"],
        legend: {
            show: true,
            position: "bottom",
            markers: { offsetX: -3 },
            labels: {
                useSeriesColors: true,
            },
        },
        dataLabels: {
            enabled: false,
        },
        stroke: {
            show: false,
            curve: "straight",
        },
        colors: ["oklch(var(--p))", "oklch(var(--su))", "oklch(var(--er))", "oklch(var(--n))"],
        states: {
            hover: {
                filter: {
                    type: "none",
                },
            },
        },
        tooltip: {
            enabled: true,
        },
    }));
});
