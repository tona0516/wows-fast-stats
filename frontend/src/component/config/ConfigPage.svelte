<script lang="ts">
import AlertPlayer from "./internal/AlertPlayer.svelte";
import Display from "./internal/Display.svelte";
import Other from "./internal/Other.svelte";
import Required from "./internal/Required.svelte";
import TeamSummary from "./internal/TeamSummary.svelte";

const tabName = "tab_config";
</script>

<div class="tabs tabs-lift">
  <input
    type="radio"
    name={tabName}
    class="tab"
    aria-label="必須設定"
    checked={true}
  />
  <div class="tab-content border-base-300 p-6">
    <Required />
  </div>

  <input type="radio" name={tabName} class="tab" aria-label="表示設定" />
  <div class="tab-content border-base-300 p-6">
    <Display />
  </div>

  <input
    type="radio"
    name={tabName}
    class="tab"
    aria-label="チームサマリー設定"
  />
  <div class="tab-content border-base-300 p-6">
    <TeamSummary />
  </div>

  <input
    type="radio"
    name={tabName}
    class="tab"
    aria-label="プレイヤーリスト設定"
  />
  <div class="tab-content border-base-300 p-6">
    <AlertPlayer on:AddAlertPlayer on:EditAlertPlayer on:RemoveAlertPlayer />
  </div>

  <input type="radio" name={tabName} class="tab" aria-label="その他設定" />
  <div class="tab-content border-base-300 p-6">
    <Other />
  </div>
</div>
