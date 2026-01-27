<script lang="ts">
    import BasicColumnSettings from "./BasicColumnSettings.svelte";
    import StatsColumnSettings from "./StatsColumnSettings.svelte";
</script>

<div class="card bg-base-100 shadow-xl rounded-xl p-6">
  <span class="text-2xl font-bold mb-4">カラム別表示設定</span>
  <BasicColumnSettings />
  <StatsColumnSettings />
</div>
