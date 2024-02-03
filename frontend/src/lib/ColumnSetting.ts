import type { DigitKey, OverallKey, ShipKey } from "src/lib/types";

export interface ColumnSetting {
  key: string;
  ship: {
    key?: ShipKey;
    value: boolean;
  };
  overall: {
    key?: OverallKey;
    value: boolean;
  };
  digit: {
    key?: DigitKey;
    value: number;
  };
}
