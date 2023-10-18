export abstract class AbstractColumn<T> {
  constructor(
    readonly displayKey: T,
    readonly minDisplayName: string,
    readonly fullDisplayName: string,
    readonly innerColumnNumber: number,
  ) {}

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  abstract getSvelteComponent(): any;
  abstract shouldShowColumn(): boolean;
}
