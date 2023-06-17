<script lang="ts">
import { get } from "svelte/store";
import { storedBattle, storedUserConfig, storedSummaryResult } from "../stores";
import AvgCheckboxTableData from "../tabledata_component/AvgCheckboxTableData.svelte";
import PlayerNameTableData from "../tabledata_component/PlayerNameTableData.svelte";
import ShipInfoTableData from "../tabledata_component/ShipInfoTableData.svelte";
import GenericTableData from "../tabledata_component/GenericTableData.svelte";
import PairTableData from "../tabledata_component/PairTableData.svelte";
import ShipTypeRateTableData from "../tabledata_component/ShipTypeRateTableData.svelte";
import TierRateTableData from "../tabledata_component/TierRateTableData.svelte";
import NoData from "../tabledata_component/NoData.svelte";
import { decideDisplayPattern } from "../util";
import { ComponenInfo, ComponentList } from "../ComponentList";
import { StatsCategory } from "../enums";
import ColorDescription from "../other_component/ColorDescription.svelte";
import Ofuse from "../other_component/Ofuse.svelte";

let battle = get(storedBattle);
storedBattle.subscribe((it) => (battle = it));

let userConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => (userConfig = it));

let summaryResult = get(storedSummaryResult);
storedSummaryResult.subscribe((it) => {
  summaryResult = it;
});

const basicComponents = new ComponentList(StatsCategory.Basic, [
  new ComponenInfo("is_in_avg", AvgCheckboxTableData),
  new ComponenInfo("player_name", PlayerNameTableData),
  new ComponenInfo("ship_info", ShipInfoTableData, { column: 2 }),
]);

const shipComponents = new ComponentList(StatsCategory.Ship, [
  new ComponenInfo("pr", GenericTableData),
  new ComponenInfo("damage", GenericTableData),
  new ComponenInfo("win_rate", GenericTableData, { unit: "%" }),
  new ComponenInfo("kd_rate", GenericTableData),
  new ComponenInfo("kill", GenericTableData),
  new ComponenInfo("exp", GenericTableData),
  new ComponenInfo("battles", GenericTableData),
  new ComponenInfo("survived_rate", PairTableData, {
    unit: "%",
    key1: "win_survived_rate",
    key2: "lose_survived_rate",
  }),
  new ComponenInfo("hit_rate", PairTableData, {
    unit: "%",
    key1: "main_battery_hit_rate",
    key2: "torpedoes_hit_rate",
  }),
  new ComponenInfo("planes_killed", GenericTableData),
]);

const overallComponents = new ComponentList(StatsCategory.Overall, [
  new ComponenInfo("damage", GenericTableData),
  new ComponenInfo("win_rate", GenericTableData, { unit: "%" }),
  new ComponenInfo("kd_rate", GenericTableData),
  new ComponenInfo("kill", GenericTableData),
  new ComponenInfo("death", GenericTableData),
  new ComponenInfo("exp", GenericTableData),
  new ComponenInfo("battles", GenericTableData),
  new ComponenInfo("survived_rate", PairTableData, {
    unit: "%",
    key1: "win_survived_rate",
    key2: "lose_survived_rate",
  }),
  new ComponenInfo("avg_tier", GenericTableData),
  new ComponenInfo("using_ship_type_rate", ShipTypeRateTableData),
  new ComponenInfo("using_tier_rate", TierRateTableData),
]);
</script>

{#if battle}
  <div class="m-2">
    <table class="table table-sm table-bordered table-text-color">
      {#each battle.teams as team}
        {@const displays = userConfig.displays}
        {@const basicColspan = basicComponents.columnCount(displays)}
        {@const shipColspan = shipComponents.columnCount(displays)}
        {@const overallColspan = overallComponents.columnCount(displays)}

        <thead>
          <tr>
            {#if basicColspan > 0}
              <th colspan="{basicColspan}">{basicComponents.minColumnName()}</th
              >
            {/if}
            {#if shipColspan > 0}
              <th colspan="{shipColspan}">{shipComponents.minColumnName()}</th>
            {/if}
            {#if overallColspan > 0}
              <th colspan="{overallColspan}"
                >{overallComponents.minColumnName()}</th
              >
            {/if}
          </tr>
          <tr>
            <!-- basic -->
            {#each basicComponents.list as c}
              {#if c.shouldShowColumn(displays, basicComponents.category)}
                <th colspan="{c.option.column}">{c.minColumnName()}</th>
              {/if}
            {/each}

            <!-- ship -->
            {#each shipComponents.list as c}
              {#if c.shouldShowColumn(displays, shipComponents.category)}
                <th colspan="{c.option.column}">{c.minColumnName()}</th>
              {/if}
            {/each}

            <!-- overall -->
            {#each overallComponents.list as c}
              {#if c.shouldShowColumn(displays, overallComponents.category)}
                <th colspan="{c.option.column}">{c.minColumnName()}</th>
              {/if}
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each team.players as player}
            {@const statsPattern = userConfig.stats_pattern}
            {@const displayPattern = decideDisplayPattern(player, statsPattern)}
            <tr>
              <!-- basic -->
              {#each basicComponents.list as c}
                <svelte:component
                  this="{c.component}"
                  player="{player}"
                  statsPattern="{statsPattern}"
                  on:UpdateAlertPlayer
                  on:RemoveAlertPlayer
                  on:CheckPlayer
                />
              {/each}

              <NoData
                shipColspan="{shipColspan}"
                overallColspan="{overallColspan}"
                displayPattern="{displayPattern}"
              />

              <!-- ship -->
              {#each shipComponents.list as c}
                {#if c.shouldShowValue(displays, shipComponents.category, displayPattern)}
                  <svelte:component
                    this="{c.component}"
                    player="{player}"
                    statsPattern="{statsPattern}"
                    statsCatetory="{shipComponents.category}"
                    key="{c.key}"
                    option="{c.option}"
                  />
                {/if}
              {/each}

              <!-- overall -->
              {#each overallComponents.list as c}
                {#if c.shouldShowValue(displays, overallComponents.category, displayPattern)}
                  <svelte:component
                    this="{c.component}"
                    player="{player}"
                    statsPattern="{statsPattern}"
                    statsCatetory="{overallComponents.category}"
                    key="{c.key}"
                    option="{c.option}"
                  />
                {/if}
              {/each}
            </tr>
          {/each}
        </tbody>
      {/each}
    </table>
  </div>

  {#if summaryResult}
    <div class="mx-4 d-flex flex-row centerize">
      <div class="mx-2">
        <h6 class="text-center">戦闘情報</h6>

        <table class="table table-sm table-text-color w-auto">
          <tbody>
            <tr>
              <td class="td-string">日時</td>
              <td class="td-string">{battle.meta.date}</td>
            </tr>

            <tr>
              <td class="td-string">戦闘タイプ</td>
              <td class="td-string">{battle.meta.type}</td>
            </tr>

            <tr>
              <td class="td-string">マップ</td>
              <td class="td-string">{battle.meta.arena}</td>
            </tr>
          </tbody>
        </table>
      </div>

      <div class="mx-2">
        <h6 class="text-center">チーム平均</h6>

        <table class="table table-sm table-text-color w-auto">
          <thead>
            <tr>
              <th></th>
              <th colspan="{summaryResult.shipColspan}">艦成績</th>
              <th colspan="{summaryResult.overallColspan}">総合成績</th>
            </tr>
            <tr>
              <th></th>
              {#each summaryResult.labels as label}
                <th>{label}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            <tr>
              <td class="td-string">{battle.teams[0].name}</td>
              {#each summaryResult.friends as friend}
                <td class="td-number">{friend}</td>
              {/each}
            </tr>

            <tr>
              <td class="td-string">{battle.teams[1].name}</td>
              {#each summaryResult.enemies as enemy}
                <td class="td-number">{enemy}</td>
              {/each}
            </tr>

            <tr>
              <td class="td-string">差</td>
              {#each summaryResult.diffs as diff}
                <td class="td-number {diff.colorClass}">{diff.value}</td>
              {/each}
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  {/if}
{/if}

<div class="m-2 centerize">
  <div>
    <ColorDescription />
  </div>
</div>

<div class="m-2 centerize">
  <div>
    <Ofuse />
  </div>
</div>
