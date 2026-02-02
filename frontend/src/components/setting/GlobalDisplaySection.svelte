<script lang="ts">
  import { storedDisplayPref } from "@libs/stores";
  import { Theme } from "@libs/Theme";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import { get } from "svelte/store";

  const ZOOM_RATES = [
    50, 67, 75, 80, 90, 100, 110, 120, 125, 133, 150, 167, 175, 200,
  ] as const;

  const STATS_EXTRAS: Readonly<Map<string, string>> = new Map<string, string>([
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
  <span class="text-xl font-bold mb-6">全体表示設定</span>
  <div class="form-control mt-4">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="label font-bold">テーマ</label>
    <select class="select" data-choose-theme>
      {#each Theme.getAll() as theme}
        <option value={theme}>{theme}</option>
      {/each}
    </select>
  </div>
  <div class="form-control mt-4">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="label font-bold">UIサイズ</label>
    <select
      class="select"
      bind:value={$storedDisplayPref.zoomRate}
      on:change={onChangePref}
    >
      {#each ZOOM_RATES as zr}
        <option selected={zr === $storedDisplayPref.zoomRate} value={zr}
          >{zr}%</option
        >
      {/each}
    </select>
  </div>
  <div class="form-control mt-4">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="label font-bold">統計パターン</label>
    <select
      class="select select-bordered"
      bind:value={$storedDisplayPref.statsExtra}
      on:change={onChangePref}
    >
      {#each STATS_EXTRAS as se}
        <option selected={se[0] === $storedDisplayPref.statsExtra} value={se[0]}
          >{se[1]}</option
        >
      {/each}
    </select>
  </div>
  <div class="form-control mt-4">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="label cursor-pointer">
      <input
        class="toggle toggle-success"
        type="checkbox"
        bind:checked={$storedDisplayPref.showBoarder}
        on:change={onChangePref}
      />
      <span class="label-text font-bold">テーブルの枠線を表示する</span>
    </label>
  </div>
</div>
