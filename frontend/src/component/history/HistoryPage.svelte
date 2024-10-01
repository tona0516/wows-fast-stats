<script lang="ts">
  import StatsTable from "src/component/main/internal/StatsTable.svelte";
  import BattleMeta from "src/component/main/internal/BattleMeta.svelte";
  import { AllBattleHistories, BattleHistory } from "wailsjs/go/main/App";
  import type { data } from "wailsjs/go/models";
  import { storedConfig } from "src/stores";
  import UkTable from "src/component/common/uikit/UkTable.svelte";
  import UkIconButton from "src/component/common/uikit/UkIconButton.svelte";
  import { StatsTableOptions } from "src/lib/StatsTableOptions";

  const options = new StatsTableOptions(false, "pvp_all");

  let battleHistoryKeys: string[];
  let selectedBattle: data.Battle;

  (async () => {
    battleHistoryKeys = await AllBattleHistories();
  })();
</script>

<div class="uk-padding-small">
  {#if battleHistoryKeys}
    <div class="uk-margin-small uk-flex uk-flex-center">
      <UkTable>
        <thead>
          {#each ["日時", "戦闘タイプ", "自艦", "マップ", ""] as columns}
            <th class="uk-text-center">{columns}</th>
          {/each}
        </thead>
        <tbody>
          {#each battleHistoryKeys as key}
            {@const info = key.split("_")}
            {@const date = info[1]}
            {@const type = info[2]}
            {@const ownShip = info[3]}
            {@const arena = info[4]}
            <tr class="uk-text-center">
              <td>{date}</td>
              <td>{type}</td>
              <td>{ownShip}</td>
              <td>{arena}</td>
              <td>
                <UkIconButton
                  name="arrow-right"
                  onclick={async () =>
                    (selectedBattle = await BattleHistory(key))}
                />
              </td></tr
            >
          {/each}
        </tbody>
      </UkTable>
    </div>
  {/if}
</div>

<div class="uk-margin-small">
  {#if selectedBattle}
    {@const teams = selectedBattle.teams}
    {@const meta = selectedBattle.meta}
    {@const config = $storedConfig}

    <div class="uk-flex uk-flex-center">
      <StatsTable {teams} {config} {options} />
    </div>

    <div class="uk-flex uk-flex-center">
      <BattleMeta {meta} />
    </div>
  {/if}
</div>
