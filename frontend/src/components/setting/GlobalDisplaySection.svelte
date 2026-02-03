<script lang="ts">
  import { storedDisplayPref } from "@libs/stores";
  import { Theme } from "@libs/Theme";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import { get } from "svelte/store";

  const ZOOM_RATES = [
    50, 67, 75, 80, 90, 100, 110, 120, 125, 133, 150, 167, 175, 200,
  ] as const;

  const STATS_EXTRAS = new Map<string, string>([
    ["pvpAll", "ランダム戦"],
    ["pvpSolo", "ランダム戦(ソロ)"],
    ["rankSolo", "ランク戦"],
  ]);

  const onChangePref = async () => {
    const pref = get(storedDisplayPref);
    if (!pref) return;
    await SaveDisplayPref(JSON.stringify(pref, null, 2));
  };
</script>

<div class="bg-base-100 shadow-xl rounded-xl p-4">
  <span class="text-xl font-bold">全体表示設定</span>
  <table class="table">
    <tbody>
      <tr>
        <td class="font-bold">テーマ</td>
        <td class="text-right">
          <select class="select" data-choose-theme>
            {#each Theme.getAll() as theme}
              <option value={theme}>{theme}</option>
            {/each}
          </select>
        </td>
      </tr>
      <tr>
        <td class="font-bold">UIサイズ</td>
        <td class="text-right">
          <select
            class="select"
            bind:value={$storedDisplayPref.zoomRate}
            on:change={onChangePref}
          >
            {#each ZOOM_RATES as zr}
              <option value={zr}>{zr}%</option>
            {/each}
          </select>
        </td>
      </tr>
      <tr>
        <td class="font-bold">統計パターン</td>
        <td class="text-right">
          <select
            class="select"
            bind:value={$storedDisplayPref.statsExtra}
            on:change={onChangePref}
          >
            {#each STATS_EXTRAS as [value, label]}
              <option {value}>{label}</option>
            {/each}
          </select>
        </td>
      </tr>
      <tr>
        <td class="font-bold">テーブルの枠線を表示する</td>
        <td class="text-right">
          <input
            class="toggle toggle-success"
            type="checkbox"
            bind:checked={$storedDisplayPref.showBoarder}
            on:change={onChangePref}
          />
        </td>
      </tr>
    </tbody>
  </table>
</div>
