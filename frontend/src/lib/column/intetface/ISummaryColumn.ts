import type { StatsCategory } from "src/lib/types";
import type { domain } from "wailsjs/go/models";

export interface ISummaryColumn {
  value(player: domain.Player): number;
  digit(): number;
  category(): StatsCategory;
}
