import { LitElement, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import Chart, { ChartConfiguration } from "chart.js/auto";
import { responsiveChartStyles } from "./chart-styles";

// TODO: optimize chartjs import to reduce bundle size

@customElement("bar-line-chart")
export default class BarLineChart extends LitElement {
  static override styles = responsiveChartStyles;
  @property({ type: Array }) chartDatasets__bars = [];
  @property({ type: Array }) chartLabels = [];
  @property({ type: String }) chartTitle = "";

  override firstUpdated() {
    const chartConfig = {
      type: "bar",
      data: {
        datasets: this.chartDatasets__bars,
        labels: this.chartLabels,
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
          },
          y: {
            stacked: true,
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
    console.debug("render called");
    console.debug(this.chartDatasets__bars);
    console.debug(this.chartLabels);
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
