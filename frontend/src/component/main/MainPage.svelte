<script lang="ts">
  import BattleMeta from "src/component/main/internal/BattleMeta.svelte";
  import StatisticsTable from "src/component/main/internal/StatsTable.svelte";
  import {
    storedBattle,
    storedConfig,
    storedSummary,
    storedRequiredConfigError,
  } from "src/stores";
  import Menu from "./internal/Menu.svelte";
  import Summary from "./internal/Summary.svelte";
  import Ofuse from "./internal/Ofuse.svelte";
  import UkSpinner from "../common/uikit/UkSpinner.svelte";
  import { FetchProxy } from "src/lib/FetchProxy";
  import { Notifier } from "src/lib/Notifier";
  import { LogInfo } from "wailsjs/go/main/App";

  let isLoading = false;

  export const fetchBattle = async () => {
    try {
      isLoading = true;

      const start = new Date().getTime();
      await FetchProxy.getBattle();
      const elapsed = (new Date().getTime() - start) / 1000;

      Notifier.success(`データ取得完了: ${elapsed.toFixed(1)}秒`);

      LogInfo("fetch success", { "duration(s)": elapsed.toFixed(1) });
    } catch (error) {
      Notifier.failure(error);
    } finally {
      isLoading = false;
    }
  };
</script>

<!-- Note: Use the same color as that of body.  -->
<div class="uk-padding-small uk-light uk-background-secondary">
  <div class="uk-margin-small uk-flex uk-flex-center">
    <Menu />
  </div>

  <div class="uk-margin-small">
    {#if $storedBattle}
      {@const teams = $storedBattle.teams}
      {@const meta = $storedBattle.meta}
      {@const config = $storedConfig}

      <div class="uk-flex uk-flex-center">
        <StatisticsTable
          {teams}
          {config}
          on:EditAlertPlayer
          on:RemoveAlertPlayer
        />
      </div>

      <div class="uk-flex uk-flex-center">
        <BattleMeta {meta} />
      </div>

      {#if $storedSummary}
        {@const summary = $storedSummary}
        <div class="uk-flex uk-flex-center">
          <Summary {summary} />
        </div>
      {/if}
    {:else}
      <p class="uk-text-center">
        {#if $storedRequiredConfigError.valid}
          戦闘中ではありません。開始時に自動的にリロードします。
        {:else}
          設定画面から初期設定を行ってください。
        {/if}
      </p>
    {/if}
  </div>

  <div class="uk-margin-small">
    <Ofuse />
  </div>

  {#if isLoading}
    <div class="uk-overlay-default">
      <div class="uk-position-center">
        <UkSpinner />
      </div>
    </div>
  {/if}
</div>
