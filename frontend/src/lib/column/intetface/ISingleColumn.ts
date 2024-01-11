import type { model } from "wailsjs/go/models";

export interface ISingleColumn {
  getTdClass(player: model.Player): string;
  getDisplayValue(player: model.Player): string;
  getTextColorCode(player: model.Player): string;
}
