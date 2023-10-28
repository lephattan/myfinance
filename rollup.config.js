/**
 * @license
 * Copyright 2018 Google LLC
 * SPDX-License-Identifier: BSD-3-Clause
 */

import summary from "rollup-plugin-summary";
import resolve from "@rollup/plugin-node-resolve";
import replace from "@rollup/plugin-replace";
import { nodeResolve } from "@rollup/plugin-node-resolve";

const options = {
  onwarn(warning) {
    if (warning.code !== "THIS_IS_UNDEFINED") {
      console.error(`(!) ${warning.message}`);
    }
  },
  plugins: [
    replace({ "Reflect.decorate": "undefined" }),
    nodeResolve({
      exportConditions: ["development"],
    }),
    resolve(),
    summary(),
  ],
};

function component(name) {
  return {
    input: `build/${name}.js`,
    output: {
      file: `./assets/js/web-components/${name}.js`,
      format: "esm",
    },
    ...options,
  };
}

export default [component("pie-chart"), component("bar-line-chart")];
