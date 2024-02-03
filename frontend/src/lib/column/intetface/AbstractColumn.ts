export abstract class AbstractColumn {
  constructor(
    readonly key: string,
    readonly header: string,
    readonly innerColumnCount: number,
  ) {}

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  abstract svelteComponent(): any;
  abstract shouldShow(): boolean;
}
