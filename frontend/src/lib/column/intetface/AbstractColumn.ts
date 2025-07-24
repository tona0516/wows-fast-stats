import type { data } from "wailsjs/go/models";

export abstract class AbstractColumn {
  constructor(
    readonly key: string,
    readonly header: string,
  ) {}

  // biome-ignore lint/suspicious/noExplicitAny: <explanation>
  abstract svelteComponent(): any;
  abstract shouldShow(): boolean;
  abstract getBackgroundColorCode(player: data.Player): string | undefined;
}
