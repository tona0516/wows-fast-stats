export interface StackedBarGraphItem {
  label: string;
  colorCode: string;
  value: number;
}

export interface StackedBarGraphParam {
  digit: number;
  items: StackedBarGraphItem[];
}
