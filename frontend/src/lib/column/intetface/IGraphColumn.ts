import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
import type { model } from "wailsjs/go/models";

export interface IGraphColumn {
  getGraphParam(player: model.Player): StackedBarGraphParam;
}
