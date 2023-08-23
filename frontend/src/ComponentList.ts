import { Const } from "src/Const";
import { DisplayPattern, StatsCategory } from "src/enums";
import type { domain } from "wailsjs/go/models";

export interface ComponentOption {
  column?: number;
  unit?: string;
  key1?: string;
  key2?: string;
}

export class ComponenInfo {
  columnKey: string;
  component: any;
  option: ComponentOption;

  constructor(key: string, component: any, option: ComponentOption = {}) {
    this.columnKey = key;
    this.component = component;
    this.option = option;
    option.column ??= 1;
    option.unit ??= "";
    option.key1 ??= "";
    option.key2 ??= "";
  }

  minColumnName(): string {
    return Const.COLUMN_NAMES[this.columnKey].min;
  }

  shouldShowColumn(
    displays: domain.Displays,
    category: StatsCategory
  ): boolean {
    return displays[category][this.columnKey] === true;
  }

  shouldShowValue(
    displays: domain.Displays,
    category: StatsCategory,
    displayPattern: DisplayPattern
  ): boolean {
    if (!this.shouldShowColumn(displays, category)) {
      return false;
    }

    if (
      [DisplayPattern.Private, DisplayPattern.NoData].includes(displayPattern)
    ) {
      return false;
    }

    if (
      category === StatsCategory.Ship &&
      displayPattern === DisplayPattern.NoShipStats
    ) {
      return false;
    }

    return true;
  }
}

export class ComponentList {
  category: StatsCategory;
  list: ComponenInfo[];

  constructor(category: StatsCategory, list: ComponenInfo[]) {
    (this.category = category), (this.list = list);
  }

  minColumnName(): string {
    return Const.COLUMN_NAMES[this.category].min;
  }

  columnCount(displays: domain.Displays): number {
    return this.list
      .filter((it) => it.shouldShowColumn(displays, this.category))
      .reduce((a, it) => a + it.option.column, 0);
  }
}
