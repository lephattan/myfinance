import { css } from "lit";

export const responsiveChartStyles = css`
  :host div[container] {
    position: relative;
    margin: auto;
    /* height: 80vh; */
    width: 100%;
    height: 100%;
  }
  ,
  :host [canvas] {
    margin: auto;
  }
`;
