import { LitElement, css, html } from "lit";
import { customElement, property } from "lit/decorators.js";

@customElement("pie-chart")
export class PieChart extends LitElement {
  static override styles = css`
    :host {
      color: blue;
    }
  `;

  @property()
  name?: string = "there 2";

  override render() {
    return html`<p>Hello, ${this.name}!</p>`;
  }
}
