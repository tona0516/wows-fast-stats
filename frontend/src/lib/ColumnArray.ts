import { DispName } from "src/lib/DispName";
import type { IColumn } from "src/lib/column/intetface/IColumn";
import type { ColumnCategory } from "src/lib/types";

export class ColumnArray<T> extends Array<IColumn<any>> {
  constructor(
    private category: ColumnCategory,
    private items: Array<IColumn<any>>,
  ) {
    super(...items);
  }

  categoryName(): string {
    const label = DispName.COLUMN_CATEGORIES.find(
      (pair) => pair.first === this.category,
    );
    if (!label) {
      throw Error(`unexpected error: ColumnCategory: ${this.category}`);
    }

    return label.second;
  }

  columnCount(): number {
    return this.items
      .filter((it) => it.shouldShowColumn())
      .reduce((a, it) => a + it.countInnerColumn(), 0);
  }
}
