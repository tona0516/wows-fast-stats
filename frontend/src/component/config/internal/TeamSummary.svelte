<script lang="ts">
import { Notifier } from "src/lib/Notifier";
import { storedConfig } from "src/stores";
import { UpdateUserConfig } from "wailsjs/go/main/App";

let isLoading = false;

$: inputConfig = $storedConfig;
$: isValidMinShipBattles =
  inputConfig.team_summary.min_ship_battles > 0 &&
  Number.isSafeInteger(inputConfig.team_summary.min_ship_battles);
$: isValidMinOverallBattles =
  inputConfig.team_summary.min_overall_battles > 0 &&
  Number.isSafeInteger(inputConfig.team_summary.min_overall_battles);
$: isValidAll = isValidMinShipBattles && isValidMinOverallBattles;

const clickApply = async () => {
  if (!isValidAll) {
    return;
  }

  try {
    isLoading = true;
    await UpdateUserConfig(inputConfig);
    Notifier.success("設定を更新しました");
  } catch (error) {
    inputConfig.team_summary = $storedConfig.team_summary;
    Notifier.failure(error);
  } finally {
    isLoading = false;
  }
};
</script>

<div>
  <h5>チーム平均に含める最小戦闘数</h5>
  <div>
    <div>艦戦闘数</div>
    <input
      class="input"
      type="number"
      bind:value={inputConfig.team_summary.min_ship_battles}
    />
    {#if !isValidMinShipBattles}
      <div>1以上の整数を入力してください。</div>
    {/if}
  </div>

  <div>
    <div>総合戦闘数</div>
    <input
      class="input"
      type="number"
      bind:value={inputConfig.team_summary.min_overall_battles}
    />
    {#if !isValidMinOverallBattles}
      <div>1以上の整数を入力してください。</div>
    {/if}
  </div>
</div>

<div>
  <div class="flex">
    <button
      class="btn btn-primary"
      disabled={isLoading || !isValidAll}
      on:click={clickApply}
    >
      {#if isLoading}
        <span class="loading loading-spinner"></span>
      {:else}
        保存
      {/if}
    </button>
  </div>
</div>
