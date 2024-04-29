export interface Cell {
  text: string;
  textColor?: string;
  bgColor?: string;
}

export type Row = Cell[];

export interface DispColorTableInfo {
  headers: string[];
  rows: Row[];
}
