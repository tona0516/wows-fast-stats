import SingleTableData from "src/component/main/internal/table_data/SingleTableData.svelte";
import type { IColumn } from "src/lib/column/intetface/IColumn";
import type { domain } from "wailsjs/go/models";

export abstract class AbstractSingleColumn<T> implements IColumn<T> {
  // IColumn methods
  abstract displayKey(): T;
  abstract minDisplayName(): string;
  abstract fullDisplayName(): string;
  abstract shouldShowColumn(): boolean;

  countInnerColumn(): number {
    return 1;
  }

  svelteComponent() {
    return SingleTableData;
  }

  // AbstractSingleColumn methods
  abstract tdClass(player: domain.Player): string;
  abstract displayValue(player: domain.Player): string;
  abstract textColorCode(player: domain.Player): string;
}
