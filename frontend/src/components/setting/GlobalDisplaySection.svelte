<script lang="ts">
  import { STATS_EXTRAS, ZOOM_RATES } from "@libs/constants";
  import { storedDisplayPref } from "@libs/stores";
  import { Theme } from "@libs/Theme";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import { get } from "svelte/store";

  const onChangePref = async () => {
    const pref = get(storedDisplayPref);
    if (!pref) return;
    await SaveDisplayPref(JSON.stringify(pref, null, 2));
  };
</script>

<div class="card bg-base-100 shadow-xl rounded-xl p-6">
  <span class="text-2xl font-bold mb-4">全体表示設定</span>
  <div class="form-control mb-2">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="label font-bold">テーマ</label>
    <select class="select select-bordered w-full" data-choose-theme>
      {#each Theme.getAll() as theme}
        <option value={theme}>{theme}</option>
      {/each}
    </select>
  </div>
  <div class="form-control mb-2">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="label font-bold">UIサイズ</label>
    <select
      class="select select-bordered w-full"
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
  <div class="form-control">
    <!-- svelte-ignore a11y-label-has-associated-control -->
    <label class="label font-bold">統計パターン</label>
    <select
      class="select select-bordered w-full my-2"
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
</div>
