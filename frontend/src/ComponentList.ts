import type { vo } from "../wailsjs/go/models";
import { Const } from "./Const";
import { DisplayPattern, StatsCategory } from "./enums";

export interface ComponentOption {
  column?: number;
  unit?: string;
  key1?: string;
  key2?: string;
}

export class ComponenInfo {
  key: string;
  component: any;
  option: ComponentOption;

  constructor(key: string, component: any, option: ComponentOption = {}) {
    this.key = key;
    this.component = component;
    this.option = option;
    option.column ??= 1;
    option.unit ??= "";
    option.key1 ??= "";
    option.key2 ??= "";
  }

  minColumnName(): string {
    return Const.COLUMN_NAMES[this.key].min;
  }

  shouldShowColumn(displays: vo.Displays, category: StatsCategory): boolean {
    return displays[category][this.key] === true;
  }

  shouldShowValue(
    displays: vo.Displays,
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

  columnCount(displays: vo.Displays): number {
    return this.list
      .filter((it) => it.shouldShowColumn(displays, this.category))
      .reduce((a, it) => a + it.option.column, 0);
  }
}
