<script lang="ts">
  import clone from "clone";
  import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
  import { FetchProxy } from "src/lib/FetchProxy";
  import { Notifier } from "src/lib/Notifier";
  import { storedConfig } from "src/stores";
  import { ApplyUserConfig } from "wailsjs/go/main/App";
  import { model } from "wailsjs/go/models";

  export let inputConfig: model.UserConfig;

  let isLoading = false;
  let teamAverage: model.TeamAverage = inputConfig.team_average;

  $: isValidMinShipBattles =
    teamAverage.min_ship_battles > 0 &&
    Number.isSafeInteger(teamAverage.min_ship_battles);
  $: isValidMinOverallBattles =
    teamAverage.min_overall_battles > 0 &&
    Number.isSafeInteger(teamAverage.min_overall_battles);
  $: isValidAll = isValidMinShipBattles && isValidMinOverallBattles;

  const clickApply = async () => {
    if (!isValidAll) {
      return;
    }

    try {
      isLoading = true;
      inputConfig.team_average = teamAverage;
      await ApplyUserConfig(inputConfig);
      await FetchProxy.getConfig();
      Notifier.success("設定を更新しました");
    } catch (error) {
      inputConfig = clone($storedConfig);
      Notifier.failure(error);
    } finally {
      isLoading = false;
    }
  };
</script>

<div class="uk-padding-small">
  <h5>チーム平均に含める最小戦闘数</h5>
  <div class="uk-margin-small-bottom">
    <div>艦戦闘数</div>
    <input
      class="uk-input uk-form-width-small"
      type="number"
      bind:value={teamAverage.min_ship_battles}
    />
    {#if !isValidMinShipBattles}
      <div class="uk-text-danger">1以上の整数を入力してください。</div>
    {/if}
  </div>

  <div class="uk-margin-small-bottom">
    <div>総合戦闘数</div>
    <input
      class="uk-input uk-form-width-small"
      type="number"
      bind:value={teamAverage.min_overall_battles}
    />
    {#if !isValidMinOverallBattles}
      <div class="uk-text-danger">1以上の整数を入力してください。</div>
    {/if}
  </div>
</div>

<div class="uk-padding-small">
  <div class="uk-flex">
    <button
      class="uk-button uk-button-primary uk-text-nowrap"
      disabled={isLoading || !isValidAll}
      on:click={clickApply}
    >
      {#if isLoading}
        <UkSpinner />
      {:else}
        保存
      {/if}
    </button>
  </div>
</div>
