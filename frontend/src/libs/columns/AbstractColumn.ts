import type { ColorCode } from "@libs/ColorCode";
import type { Optional } from "@libs/types";
import type { data } from "@wails/go/models";

export abstract class AbstractColumn {
  constructor(
    readonly key: string,
    readonly header: string,
  ) {}

  abstract needsShow(): boolean;
  abstract getTableDataComponent(): any;

  getTextColorCode(_: data.Player): Optional<ColorCode> {
    return undefined;
  }

  getBgColorCode(_: data.Player): Optional<ColorCode> {
    return undefined;
  }
}
