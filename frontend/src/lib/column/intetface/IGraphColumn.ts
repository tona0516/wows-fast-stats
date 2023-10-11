import type { StackedBarGraphParam } from "src/lib/column/StackedBarGraphParam";
import type { domain } from "wailsjs/go/models";

export interface IGraphColumn {
  getGraphParam(player: domain.Player): StackedBarGraphParam;
}
