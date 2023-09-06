export interface Summary {
  shipColspan: number;
  overallColspan: number;
  labels: string[];
  friends: string[];
  enemies: string[];
  diffs: { value: string; colorClass: string }[];
}
