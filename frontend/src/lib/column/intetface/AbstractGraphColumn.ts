import StackedBarGraphTableData from "src/component/main/internal/table_data/StackedBarGraphTableData.svelte";
import type { IColumn } from "src/lib/column/intetface/IColumn";
import type { StackedBarGraphParam } from "src/lib/value_object/StackedBarGraphParam";
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
