import type { data } from "wailsjs/go/models";

export interface ISummaryColumn {
  value(player: data.Player): number;
}
