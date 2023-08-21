export interface StackedBarGraphParam {
  digit: number;
  items: {
    label: string;
    color: string;
    value: number;
  }[];
}
