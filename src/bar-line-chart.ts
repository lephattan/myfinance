import { LitElement, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ChartConfiguration } from "chart.js/auto";

import {
  Chart,
  BarController,
  BarElement,
  LineController,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Colors,
  Legend,
} from "chart.js";
import { responsiveChartStyles } from "./chart-styles";

Chart.register(
  BarController,
  BarElement,
  LineController,
  LineElement,
  PointElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Colors,
  Legend
);

@customElement("bar-line-chart")
export default class BarLineChart extends LitElement {
  static override styles = responsiveChartStyles;
  @property({ type: Array }) chartDatasets__bars = [];
  @property({ type: Array }) chartDatasets__lines = [];
  @property({ type: Array }) chartLabels = [];
  @property({ type: String }) chartTitle = "";
  @property({ type: String }) rightAxesTitle = "";
  @property({ type: String }) leftAxesTitle = "";
  @property({ type: String }) botAxesTitle = "";

  override firstUpdated() {
    const datasets: object[] = [];
    this.chartDatasets__bars.map((dataset: object) => {
      datasets.push({ ...dataset, type: "bar", order: 2 });
    });

    this.chartDatasets__lines.map((dataset: object) => {
      datasets.push({ ...dataset, type: "line", order: 1, yAxisID: "y2" });
    });

    const chartConfig = {
      type: "bar",
      data: {
        labels: this.chartLabels,
        datasets: datasets,
      },
      options: {
        plugins: {
          title: {
            display: true,
            text: this.chartTitle,
            padding: {
              top: 10,
              bottom: 10,
            },
          },
        },
        scales: {
          x: {
            stacked: true,
            title: {
              text: this.botAxesTitle,
              display: this.botAxesTitle !== "",
            },
          },
          y: {
            stacked: true,
            title: {
              text: this.leftAxesTitle,
              display: this.leftAxesTitle !== "",
            },
          },
          y2: {
            position: "right",
            title: {
              text: this.rightAxesTitle,
              display: true,
            },
          },
        },
        responsive: true,
        maintainAspectRatio: false,
      },
    } as ChartConfiguration;

    const ctx = this.renderRoot.querySelector(`#chart`)!;
    new Chart(ctx as HTMLCanvasElement, chartConfig);
  }

  override render() {
    return html`<div container>
      <canvas id="chart" canvas></canvas>
    </div> `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    "bar-line-chart": BarLineChart;
  }
}
