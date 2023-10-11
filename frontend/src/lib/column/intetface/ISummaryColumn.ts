import type { StatsCategory } from "src/lib/types";
import type { domain } from "wailsjs/go/models";

export interface ISummaryColumn {
  getValue(player: domain.Player): number;
  getDigit(): number;
  getCategory(): StatsCategory;
}
