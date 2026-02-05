<script lang="ts">
  import ColumnSettingTable from "./ColumnSettingTable.svelte";
  import StatsColumnOrderSettings from "./StatsColumnOrderSettings.svelte";

  export let resetToken = 0;
</script>

{#key resetToken}
  <div class="bg-base-100 shadow-xl rounded-xl p-4">
    <span class="text-xl font-bold">カラム別表示設定</span>
    <div
      class="mt-4 overflow-x-auto rounded-xl border border-base-300 bg-base-200 shadow-sm"
    >
      <ColumnSettingTable />
    </div>

    <StatsColumnOrderSettings />
  </div>
{/key}
