<script lang="ts">
import type { vo } from "../../wailsjs/go/models";
import AvgCheckboxTableData from "../tabledata_component/AvgCheckboxTableData.svelte";
import PlayerNameTableData from "../tabledata_component/PlayerNameTableData.svelte";
import ShipInfoTableData from "../tabledata_component/ShipInfoTableData.svelte";
import GenericTableData from "../tabledata_component/GenericTableData.svelte";
import PairTableData from "../tabledata_component/PairTableData.svelte";
import ShipTypeRateTableData from "../tabledata_component/ShipTypeRateTableData.svelte";
import TierRateTableData from "../tabledata_component/TierRateTableData.svelte";
import { ComponenInfo, ComponentList } from "../ComponentList";
import { StatsCategory } from "../enums";
import NoData from "../tabledata_component/NoData.svelte";
import { decideDisplayPattern } from "../util";

export let teams: vo.Team[];
export let userConfig: vo.UserConfig;
export let alertPlayers: vo.AlertPlayer[];

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
  new ComponenInfo("planes_killed", GenericTableData),
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

$: displays = userConfig.displays;
$: basicColspan = basicComponents.columnCount(displays);
$: shipColspan = shipComponents.columnCount(displays);
$: overallColspan = overallComponents.columnCount(displays);
</script>

<table class="table table-sm table-bordered table-text-color">
  {#each teams as team}
    <thead>
      <tr>
        {#if basicColspan > 0}
          <th colspan="{basicColspan}">{basicComponents.minColumnName()}</th>
        {/if}
        {#if shipColspan > 0}
          <th colspan="{shipColspan}">{shipComponents.minColumnName()}</th>
        {/if}
        {#if overallColspan > 0}
          <th colspan="{overallColspan}">{overallComponents.minColumnName()}</th
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
              userConfig="{userConfig}"
              statsPattern="{statsPattern}"
              alertPlayers="{alertPlayers}"
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
                columnKey="{c.columnKey}"
                option="{c.option}"
                customColor="{userConfig.custom_color}"
                customDigit="{userConfig.custom_digit}"
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
                columnKey="{c.columnKey}"
                option="{c.option}"
                customColor="{userConfig.custom_color}"
                customDigit="{userConfig.custom_digit}"
              />
            {/if}
          {/each}
        </tr>
      {/each}
    </tbody>
  {/each}
</table>
