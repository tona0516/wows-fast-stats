export abstract class AbstractColumn<T> {
  constructor(
    readonly displayKey: T,
    readonly minDisplayName: string,
    readonly fullDisplayName: string,
    readonly innerColumnNumber: number,
    readonly svelteComponent: any,
  ) {}

  abstract shouldShowColumn(): boolean;
}
