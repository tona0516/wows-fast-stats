<script lang="ts">
  import type { AbstractColumn } from "@libs/columns/AbstractColumn";
  import type { StatsCategory, StatsExtra } from "@libs/types";
  import type { core } from "@wails/go/models";
  import ColspanTableData from "./tabledata/ColspanTableData.svelte";
  import { AvgTierColumn } from "@libs/columns/AvgTierColumn";
  import { BattlesColumn } from "@libs/columns/BattlesColumn";
  import { DamageColumn } from "@libs/columns/DamageColumn";
  import { ExpColumn } from "@libs/columns/ExpColumn";
  import { HitRateColumn } from "@libs/columns/HitRateColumn";
  import { KDRateColumn } from "@libs/columns/KDRateColumn";
  import { KillColumn } from "@libs/columns/KillColumn";
  import { MaxDamageColumn } from "@libs/columns/MaxDamageColumn";
  import { PlanesKilledColumn } from "@libs/columns/PlanesKilledColumn";
  import { PlatoonRateColumn } from "@libs/columns/PlatoonRateColumn";
  import { PlayerNameColumn } from "@libs/columns/PlayerNameColumn";
  import { PRColumn } from "@libs/columns/PRColumn";
  import { WarshipColumn } from "@libs/columns/WarshipColumn";
  import { ShipTypeRateColumn } from "@libs/columns/ShipTypeRateColumn";
  import { SurvivedRateColumn } from "@libs/columns/SurvivedRateColumn";
  import { ThreatLevelColumn } from "@libs/columns/ThreatLevelColumn";
  import { TierRateColumn } from "@libs/columns/TierRateColumn";
  import { WinRateColumn } from "@libs/columns/WinRateColumn";
  import { EfficiencyBadgeColumn } from "@libs/columns/EfficiencyBadgeColumn";
  import { storedDisplayPref } from "@libs/stores";

  export let teams: core.Team[];

  type ColumnCategory = Readonly<"basic" | StatsCategory>;

  type RowPattern =
    | "no_column"
    | "private"
    | "no_stats"
    | "no_ship_stats"
    | "full";

  const CATEGORY_NAMES: Readonly<Map<ColumnCategory, string>> = new Map<
    ColumnCategory,
    string
  >([
    ["basic", "基本情報"],
    ["ship", "艦成績"],
    ["overall", "総合成績"],
  ]);

  class Category {
    constructor(
      public readonly value: ColumnCategory,
      public readonly columns: AbstractColumn[],
    ) {}

    header(): string {
      return CATEGORY_NAMES.get(this.value) ?? this.value;
    }

    showCount(): number {
      return this.columns.filter((col) => col.needsShow()).length;
    }
  }

  const basicCategory = new Category("basic", [
    new PlayerNameColumn(),
    new WarshipColumn(),
  ]);

  const shipCategory = new Category("ship", [
    new PRColumn("ship"),
    new WinRateColumn("ship"),
    new DamageColumn("ship"),
    new MaxDamageColumn("ship"),
    new KDRateColumn("ship"),
    new KillColumn("ship"),
    new ExpColumn("ship"),
    new BattlesColumn("ship"),
    new SurvivedRateColumn("ship"),
    new PlatoonRateColumn("ship"),
    new EfficiencyBadgeColumn("ship"),
    new PlanesKilledColumn(),
    new HitRateColumn(),
  ]);

  const overallCategory = new Category("overall", [
    new PRColumn("overall"),
    new WinRateColumn("overall"),
    new DamageColumn("overall"),
    new MaxDamageColumn("overall"),
    new KDRateColumn("overall"),
    new KillColumn("overall"),
    new ExpColumn("overall"),
    new BattlesColumn("overall"),
    new SurvivedRateColumn("overall"),
    new PlatoonRateColumn("overall"),
    new EfficiencyBadgeColumn("overall"),
    new ThreatLevelColumn(),
    new AvgTierColumn(),
    new TierRateColumn(),
    new ShipTypeRateColumn(),
  ]);

  const categories = [basicCategory, shipCategory, overallCategory];
  const allColumnCount =
    basicCategory.showCount() +
    shipCategory.showCount() +
    overallCategory.showCount();

  const getRowPattern = (
    player: core.Player,
    statsExtra: string,
    shipColumnCount: number,
    overallColumnCount: number,
  ): RowPattern => {
    if (shipColumnCount + overallColumnCount === 0) {
      return "no_column";
    }

    if (player.playerInfo.isHidden === true) {
      return "private";
    }

    const stats = player[statsExtra as StatsExtra];
    if (player.playerInfo.id === 0 || stats.overall.battles === 0) {
      return "no_stats";
    }

    if (stats.ship.battles === 0 && shipColumnCount > 0) {
      return "no_ship_stats";
    }

    return "full";
  };

  const getColumnText = (pattern: RowPattern): string => {
    switch (pattern) {
      case "private":
        return "PRIVAYE";
      case "no_stats":
        return "N/A";
      case "no_ship_stats":
        return "N/A";
      default:
        return "";
    }
  };
</script>

<div class="overflow-x-auto">
  <table
    class="{$storedDisplayPref.showBoarder
      ? 'border border-gray-500'
      : ''} border-collapse table text-nowrap w-full"
  >
    {#each teams as team, i}
      {#if team.players.length !== 0}
        <thead>
          <tr class="bg-base-300 text-xs">
            {#each categories as category}
              {#if category.showCount() > 0}
                <th
                  class="{$storedDisplayPref.showBoarder
                    ? 'border border-gray-500'
                    : ''} p-2 text-center font-bold tracking-wide"
                  colspan={category.showCount()}
                  scope="colgroup"
                >
                  {category.header()}
                </th>
              {/if}
            {/each}
          </tr>
          <tr class="bg-base-200 text-[11px]">
            {#each categories as category}
              {#each category.columns as column}
                {#if column.needsShow()}
                  <th
                    class="{$storedDisplayPref.showBoarder
                      ? 'border border-gray-500'
                      : ''} p-1 text-center font-medium whitespace-nowrap"
                    scope="col">{column.header}</th
                  >
                {/if}
              {/each}
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each team.players as player}
            {@const rowPattern = getRowPattern(
              player,
              $storedDisplayPref.statsExtra,
              shipCategory.showCount(),
              overallCategory.showCount(),
            )}
            <tr class="hover:bg-base-100/70">
              {#each basicCategory.columns as column}
                <svelte:component
                  this={column.getTableDataComponent()}
                  {column}
                  {player}
                />
              {/each}

              {#if ["no_column", "private", "no_stats"].includes(rowPattern)}
                <ColspanTableData
                  colspan={allColumnCount}
                  text={getColumnText(rowPattern)}
                />
              {:else if rowPattern === "no_ship_stats"}
                <ColspanTableData
                  colspan={shipCategory.showCount()}
                  text={getColumnText(rowPattern)}
                />
                {#each overallCategory.columns as column}
                  <svelte:component
                    this={column.getTableDataComponent()}
                    {column}
                    {player}
                  />
                {/each}
              {:else}
                {#each shipCategory.columns as column}
                  <svelte:component
                    this={column.getTableDataComponent()}
                    {column}
                    {player}
                  />
                {/each}
                {#each overallCategory.columns as column}
                  <svelte:component
                    this={column.getTableDataComponent()}
                    {column}
                    {player}
                  />
                {/each}
              {/if}
            </tr>
          {/each}
        </tbody>
      {/if}
    {/each}
  </table>

  {#if teams.every((t) => t.players.length === 0)}
    <div class="alert alert-info">
      <span>表示可能なプレイヤーデータがありません。</span>
    </div>
  {/if}
</div>
