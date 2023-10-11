export abstract class AbstractColumn<T> {
  constructor(
    private displayKey: T,
    private minDisplayName: string,
    private fullDisplayName: string,
    private innerColumnNumber: number,
    private svelteComponent: any,
  ) {}

  getDisplayKey(): T {
    return this.displayKey;
  }

  getMinDisplayName(): string {
    return this.minDisplayName;
  }

  getFullDisplayName(): string {
    return this.fullDisplayName;
  }

  getInnerColumnNumber(): number {
    return this.innerColumnNumber;
  }

  getSvelteComponent(): any {
    return this.svelteComponent;
  }

  abstract shouldShowColumn(): boolean;
}
