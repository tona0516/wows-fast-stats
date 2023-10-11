import { DispName } from "src/lib/DispName";
import type { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import type { ColumnCategory } from "src/lib/types";

export class ColumnArray extends Array<AbstractColumn<any>> {
  constructor(
    private category: ColumnCategory,
    private columns: AbstractColumn<any>[],
  ) {
    super(...columns);
  }

  dispName(): string {
    const label = DispName.COLUMN_CATEGORIES.get(this.category);
    if (!label) {
      throw Error(`unexpected error: ColumnCategory: ${this.category}`);
    }

    return label;
  }

  columnCount(): number {
    return this.columns
      .filter((it) => it.shouldShowColumn())
      .reduce((a, it) => a + it.getInnerColumnNumber(), 0);
  }
}
