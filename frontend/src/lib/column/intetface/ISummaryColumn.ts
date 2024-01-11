import type { StatsCategory } from "src/lib/types";
import type { model } from "wailsjs/go/models";

export interface ISummaryColumn {
  getValue(player: model.Player): number;
  getDigit(): number;
  getCategory(): StatsCategory;
}
