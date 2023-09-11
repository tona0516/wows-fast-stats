import { type DigitKey, type OverallKey, type ShipKey } from "src/lib/types";

export interface DisplayItem {
  name: string;
  digitKey: DigitKey;
  shipKey?: ShipKey;
  overallKey?: OverallKey;
}
