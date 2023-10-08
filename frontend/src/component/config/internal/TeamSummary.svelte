<script lang="ts">
  import clone from "clone";
  import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
  import { storedUserConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import { ApplyUserConfig, UserConfig } from "wailsjs/go/main/App";
  import { domain } from "wailsjs/go/models";

  export let inputUserConfig: domain.UserConfig;

  let isLoading = false;
  let teamAverage: domain.TeamAverage = inputUserConfig.team_average;

  $: isValidMinShipBattles =
    teamAverage.min_ship_battles > 0 &&
    Number.isSafeInteger(teamAverage.min_ship_battles);
  $: isValidMinOverallBattles =
    teamAverage.min_overall_battles > 0 &&
    Number.isSafeInteger(teamAverage.min_overall_battles);
  $: isValidAll = isValidMinShipBattles && isValidMinOverallBattles;

  const dispatch = createEventDispatcher();

  const clickApply = async () => {
    if (!isValidAll) {
      return;
    }

    try {
      isLoading = true;
      inputUserConfig.team_average = teamAverage;
      await ApplyUserConfig(inputUserConfig);

      const latest = await UserConfig();
      storedUserConfig.set(latest);
      dispatch("UpdateSuccess");
    } catch (error) {
      inputUserConfig = clone($storedUserConfig);
      dispatch("Failure", { message: error });
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
