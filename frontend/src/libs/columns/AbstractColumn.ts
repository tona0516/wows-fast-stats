import type { ColorCode } from "@libs/ColorCode";
import type { Optional } from "@libs/types";
import type { core } from "@wails/go/models";

export abstract class AbstractColumn {
  constructor(
    readonly key: string,
    readonly header: string,
  ) {}

  abstract needsShow(): boolean;
  abstract getTableDataComponent(): any;

  getTextColorCode(_: core.Player): Optional<ColorCode> {
    return undefined;
  }

  getBgColorCode(_: core.Player): Optional<ColorCode> {
    return undefined;
  }
}
