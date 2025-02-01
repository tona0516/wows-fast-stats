import type { model } from "wailsjs/go/models";

export interface ISummaryColumn {
  value(player: model.Player): number;
}
