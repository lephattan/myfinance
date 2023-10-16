import { LitElement, TemplateResult, html, css } from "lit";
import { customElement, property, state } from "lit/decorators.js";
import Chart, { ChartConfiguration } from "chart.js/auto";

@customElement("pie-chart")
export default class PieChart extends LitElement {
  static override styles = css`
    :host div[container] {
      position: relative;
      margin: auto;
      /* height: 80vh; */
      width: 100%;
    }
    ,
    :host [canvas] {
      margin: auto;
    }
  `;

  @property()
  holdings?: string;
  @property()
  chart_id?: string;
  @property()
  labelField?: string;
  @property()
  dataField?: string;
  @property({ type: String })
  chartTitle: string;

  @state()
  _chartData: object[] = [];

  constructor() {
    super();
    this.chartTitle = "";
  }

  readLabels(): string[] {
    if (this.labelField === undefined) {
      return [];
    }
    const label = this.labelField;
    return this._chartData.map((datum: Record<string, any>): string => {
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
    return this._chartData.map((datum: Record<string, any>): number => {
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
    const ctx = this.renderRoot.querySelector(`#${this.chart_id}`)!;
    new Chart(ctx as HTMLCanvasElement, chartConfig);
  }

  override render(): void | TemplateResult {
    function onEmpty() {
      return html`<p>No holdings data!</p>`;
    }

    function onError(err: any) {
      return html`<p>Error parsing holdings data: ${err}</p>`;
    }

    if (this.holdings === undefined) {
      return onEmpty();
    }

    try {
      const data = JSON.parse(this.holdings) as object[];
      this._chartData = data;
      return html`<div container>
        <canvas id="${this.chart_id}" canvas></canvas>
      </div> `;
    } catch (err) {
      return onError(err);
    }
  }
}

declare global {
  interface HTMLElementTagNameMap {
    "pie-chart": PieChart;
  }
}
