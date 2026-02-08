<script lang="ts">
  import Section from "@components/commons/Section.svelte";
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

  const rows = [
    { label: "テーマ", type: "theme" as const },
    { label: "UIサイズ", type: "zoom" as const },
    { label: "統計パターン", type: "stats" as const },
    { label: "テーブルの枠線を表示", type: "border" as const },
    { label: "K(キロ)表示", type: "prefix" as const },
  ];
</script>

<Section title="全体表示設定">
  <table class="table">
    <tbody>
      {#each rows as row}
        <tr>
          <td class="font-bold">{row.label}</td>
          <td class="text-right">
            {#if row.type === "theme"}
              <select class="select" data-choose-theme>
                {#each Theme.getAll() as theme}
                  <option value={theme}>{theme}</option>
                {/each}
              </select>
            {:else if row.type === "zoom"}
              <select
                class="select"
                bind:value={$storedDisplayPref.zoomRate}
                on:change={onChangePref}
              >
                {#each ZOOM_RATES as zr}
                  <option value={zr}>{zr}%</option>
                {/each}
              </select>
            {:else if row.type === "stats"}
              <select
                class="select"
                bind:value={$storedDisplayPref.statsExtra}
                on:change={onChangePref}
              >
                {#each STATS_EXTRAS as [value, label]}
                  <option {value}>{label}</option>
                {/each}
              </select>
            {:else if row.type === "border"}
              <input
                class="toggle toggle-success"
                type="checkbox"
                bind:checked={$storedDisplayPref.showBoarder}
                on:change={onChangePref}
              />
            {:else if row.type === "prefix"}
              <input
                class="toggle toggle-success"
                type="checkbox"
                bind:checked={$storedDisplayPref.showSiPrefix}
                on:change={onChangePref}
              />
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</Section>
