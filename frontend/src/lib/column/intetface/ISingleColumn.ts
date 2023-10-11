import type { domain } from "wailsjs/go/models";

export interface ISingleColumn {
  getTdClass(player: domain.Player): string;
  getDisplayValue(player: domain.Player): string;
  getTextColorCode(player: domain.Player): string;
}
