import type { data } from "wailsjs/go/models";

export abstract class AbstractColumn {
  constructor(
    readonly key: string,
    readonly header: string,
  ) {}

  abstract needsShow(): boolean;
  // biome-ignore lint/suspicious/noExplicitAny: <explanation>
  abstract getTableDataComponent(): any;
  abstract getBackgroundColorCode(player: data.Player): string | undefined;
}
