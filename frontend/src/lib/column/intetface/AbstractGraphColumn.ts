import type { StackedBarGraphParam } from "src/lib/StackedBarGraphParam";
import type { IColumn } from "src/lib/column/intetface/IColumn";
import StackedBarGraphTableData from "src/tabledata_component/StackedBarGraphTableData.svelte";
import type { domain } from "wailsjs/go/models";

export abstract class AbstractGraphColumn<T> implements IColumn<T> {
  // IColumn methods
  abstract displayKey(): T;
  abstract minDisplayName(): string;
  abstract fullDisplayName(): string;
  abstract shouldShowColumn(): boolean;

  countInnerColumn(): number {
    return 1;
  }

  svelteComponent() {
    return StackedBarGraphTableData;
  }

  // AbstractGraphColumn methods
  abstract displayValue(player: domain.Player): StackedBarGraphParam;
}
