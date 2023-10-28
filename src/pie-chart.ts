import { LitElement, TemplateResult, html } from "lit";
import { customElement, property } from "lit/decorators.js";
import { ChartConfiguration } from "chart.js";
import {
  Chart,
  DoughnutController,
  ArcElement,
  Tooltip,
  Colors,
  Legend,
} from "chart.js";
import { responsiveChartStyles } from "./chart-styles";

Chart.register(DoughnutController, ArcElement, Tooltip, Colors, Legend);

@customElement("pie-chart")
export default class PieChart extends LitElement {
  static override styles = responsiveChartStyles;
  @property({ type: Array }) holdings = [];
  @property({ type: String }) chartTitle = "";
  @property({ type: String }) labelField = undefined;
  @property({ type: String }) dataField = undefined;

  constructor() {
    super();
  }

  readLabels(): string[] {
    if (this.labelField === undefined) {
      return [];
    }
    const label = this.labelField;
    return this.holdings.map((datum: Record<string, any>): string => {
      if (datum[label] !== undefined) {
        return datum[label];
      }
      return "";
    });
  }

  readData(): number[] {
    if (this.dataField === undefined) {
      console.log("no dataField provided");
      return [];
    }
    const label = this.dataField;
    return this.holdings.map((datum: Record<string, any>): number => {
      if (datum[label] !== undefined) {
        return datum[label];
      }
      return 0;
    });
  }

  override firstUpdated() {
    const holding_cost = this.readData();
    const totalCost = holding_cost.reduce((a, v) => {
      return a + v;
    }, 0);

    const dataset = {
      data: holding_cost,
      hoverOffset: 30,
    };
    const labels = this.readLabels();
    const chartLabels = labels.map((label, idx) => {
      let labelData;
      try {
        labelData = holding_cost[idx];
      } catch {
        labelData = 0;
      }
      const percent = Math.round((labelData * 100 * 100) / totalCost) / 100;
      return `${label.toUpperCase()} (${percent}%)`;
    });
    const chartConfig = {
      type: "doughnut",
      data: {
        datasets: [dataset],
        labels: chartLabels,
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
        responsive: true,
        maintainAspectRatio: false,
      },
    } as ChartConfiguration;
    const ctx = this.renderRoot.querySelector(`#chart`)!;
    new Chart(ctx as HTMLCanvasElement, chartConfig);
  }

  override render(): void | TemplateResult {
    return html`<div container>
      <canvas id="chart" canvas></canvas>
    </div> `;
  }
}

declare global {
  interface HTMLElementTagNameMap {
    "pie-chart": PieChart;
  }
}
