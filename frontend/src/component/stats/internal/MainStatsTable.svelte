<script lang="ts">
  import {
    storedColumnmSettings,
    storedStatsExtra,
    storedTeamThreatLevels,
  } from "src/stores";

  import { RowPattern } from "src/lib/RowPattern";
  import type { data } from "wailsjs/go/models";
  import ColspanTableData from "./table_data/ColspanTableData.svelte";
  import { ThreatLevel } from "src/lib/column/model/ThreatLevel";
  import { PR } from "src/lib/column/model/PR";
  import { WinRate } from "src/lib/column/model/WinRate";
  import { Damage } from "src/lib/column/model/Damage";
  import { MaxDamage } from "src/lib/column/model/MaxDamage";
  import { KDRate } from "src/lib/column/model/KDRate";
  import { Kill } from "src/lib/column/model/Kill";
  import { Exp } from "src/lib/column/model/Exp";
  import { Battles } from "src/lib/column/model/Battles";
  import { PlanesKilled } from "src/lib/column/model/PlanesKilled";
  import { SurvivedRate } from "src/lib/column/model/SurvivedRate";
  import { HitRate } from "src/lib/column/model/HitRate";
  import { PlatoonRate } from "src/lib/column/model/PlatoonRate";
  import { AvgTier } from "src/lib/column/model/AvgTier";
  import { UsingTierRate } from "src/lib/column/model/UsingTierRate";
  import { UsingShipTypeRate } from "src/lib/column/model/UsingShipTypeRate";
  import { PlayerName } from "src/lib/column/model/PlayerName";
  import { ShipInfo } from "src/lib/column/model/ShipInfo";
  import type { ColumnCategory } from "src/lib/types";
  import type { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
  import { AppConstants } from "src/lib/AppConstants";

  export let teams: data.Team[];

  class Category {
    constructor(
      public readonly value: ColumnCategory,
      public readonly columns: AbstractColumn[],
    ) {}

    header(): string {
      return AppConstants.CATEGORY_NAMES.get(this.value) ?? this.value;
    }

    showCount(): number {
      return this.columns.filter((col) => col.needsShow()).length;
    }
  }

  const basicCategory = new Category("basic", [
    new PlayerName(),
    new ShipInfo(),
  ]);

  const shipCategory = new Category("ship", [
    new PR("ship"),
    new WinRate("ship"),
    new Damage("ship"),
    new MaxDamage("ship"),
    new KDRate("ship"),
    new Kill("ship"),
    new Exp("ship"),
    new Battles("ship"),
    new SurvivedRate("ship"),
    new PlatoonRate("ship"),
    new PlanesKilled(),
    new HitRate(),
  ]);

  const overallCategory = new Category("overall", [
    new PR("overall"),
    new WinRate("overall"),
    new Damage("overall"),
    new MaxDamage("overall"),
    new KDRate("overall"),
    new Kill("overall"),
    new Exp("overall"),
    new Battles("overall"),
    new SurvivedRate("overall"),
    new PlatoonRate("overall"),
    new ThreatLevel(),
    new AvgTier(),
    new UsingTierRate(),
    new UsingShipTypeRate(),
  ]);

  const categories = [basicCategory, shipCategory, overallCategory];
  const allColumnCount =
    basicCategory.showCount() +
    shipCategory.showCount() +
    overallCategory.showCount();
  const showThreatLevel = $storedColumnmSettings["threat_level"].overall;
</script>

<div class="overflow-x-auto w-screen pb-4">
  <table class="table text-nowrap">
    {#each teams as team, i}
      {#if team.players.length !== 0}
        <thead>
          {#if showThreatLevel && $storedTeamThreatLevels && $storedTeamThreatLevels[i]}
            {@const teamThreatLevel = $storedTeamThreatLevels[i]}
            <tr>
              <th colspan={allColumnCount}>
                戦力評価値平均: <span class="text-lg"
                  >{teamThreatLevel.average.toFixed(0)}</span
                >
                [確度:
                <span class="text-lg"
                  >{teamThreatLevel.accuracy.toFixed(0)}</span
                >%] [介護指数:
                <span class="text-lg"
                  >{teamThreatLevel.dissociationDegree.toFixed(0)}</span
                >%]
              </th>
            </tr>
          {/if}

          <tr>
            {#each categories as category}
              {#if category.showCount() > 0}
                <th class="p-1 text-center" colspan={category.showCount()}>
                  {category.header()}
                </th>
              {/if}
            {/each}
          </tr>
          <tr>
            {#each categories as category}
              {#each category.columns as column}
                {#if column.needsShow()}
                  <th class="p-1 text-center">{column.header}</th>
                {/if}
              {/each}
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each team.players as player}
            {@const rowPattern = RowPattern.derive(
              player,
              $storedStatsExtra,
              shipCategory.showCount(),
              overallCategory.showCount(),
            )}
            <tr>
              {#each basicCategory.columns as column}
                <svelte:component
                  this={column.getTableDataComponent()}
                  {column}
                  {player}
                />
              {/each}

              {#if [RowPattern.NO_COLUMN, RowPattern.PRIVATE, RowPattern.NO_STATS].includes(rowPattern)}
                <ColspanTableData
                  colspan={allColumnCount}
                  text={RowPattern.getColumnText(rowPattern)}
                />
              {:else if rowPattern === RowPattern.NO_SHIP_STATS}
                <ColspanTableData
                  colspan={shipCategory.showCount()}
                  text={RowPattern.getColumnText(rowPattern)}
                />
                {#each shipCategory.columns as column}
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
</div>
