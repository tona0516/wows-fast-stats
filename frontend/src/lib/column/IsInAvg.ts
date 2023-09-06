import type { IColumn } from "src/lib/column/intetface/IColumn";
import type { BasicKey } from "src/lib/types";
import AvgCheckboxTableData from "src/tabledata_component/AvgCheckboxTableData.svelte";
import type { domain } from "wailsjs/go/models";

export class IsInAvg implements IColumn<BasicKey> {
  displayKey(): BasicKey {
    return "is_in_avg";
  }

  minDisplayName(): string {
    return "";
  }

  fullDisplayName(): string {
    return "";
  }

  shouldShowColumn(): boolean {
    return true;
  }

  countInnerColumn(): number {
    return 1;
  }

  svelteComponent() {
    return AvgCheckboxTableData;
  }

  displayValue(player: domain.Player): boolean {
    return player.player_info.id !== 0;
  }
}
