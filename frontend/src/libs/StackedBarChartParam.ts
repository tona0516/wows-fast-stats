import type { ColorCode } from "./ColorCode";

export type StackedBarChartParam = {
  readonly label: string;
  readonly colorCode?: ColorCode;
  readonly value: number;
};
