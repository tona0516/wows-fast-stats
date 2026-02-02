<script lang="ts">
  import BasicColumnSettings from "./BasicColumnSettings.svelte";
  import StatsColumnSettings from "./StatsColumnSettings.svelte";

  export let resetToken = 0;
</script>

<div class="card bg-base-100 shadow-xl rounded-xl p-6">
  <span class="mb-4 text-xl font-bold">カラム別表示設定</span>
  <BasicColumnSettings />
  {#key resetToken}
    <StatsColumnSettings />
  {/key}
</div>
