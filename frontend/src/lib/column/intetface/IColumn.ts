export interface IColumn<T> {
  displayKey(): T;
  minDisplayName(): string;
  fullDisplayName(): string;
  shouldShowColumn(): boolean;
  countInnerColumn(): number;
  svelteComponent(): any;
}
