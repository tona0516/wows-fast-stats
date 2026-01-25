<script lang="ts">
  import {
    STATS_EXTRAS,
    ZOOM_RATES,
    STATS_KEYS,
    STATS_COLUMN_INFO,
    PLAYER_NAME_COLORS,
  } from "@libs/constants";
  import {
    showToast,
    storedGameClientPathError,
    storedPref,
  } from "@libs/stores";
  import { Theme } from "@libs/Theme";
  import { LoadPref, SavePref, SelectInstallPath } from "@wails/go/main/App";
  import { onMount } from "svelte";
  import { get } from "svelte/store";
  import { themeChange } from "theme-change";

  onMount(async () => {
    themeChange(false);
  });

  const onClickSelectDirectory = async () => {
    try {
      await SelectInstallPath();
      const pref = await LoadPref();
      storedPref.set(pref);
      storedGameClientPathError.set("");

      showToast("ゲームクライアントパスを設定しました");
    } catch (error) {
      const errorString = error as string;
      if (errorString.includes("C103")) {
        return;
      }

      storedGameClientPathError.set(errorString);
    }
  };

  const onChange = async () => {
    const pref = get(storedPref);
    if (!pref) return;
    await SavePref(pref);
  };
</script>

<div class="container mx-auto max-w-3xl py-3 flex flex-col gap-4">
  <!-- ゲームクライアントパス設定 -->
  <div class="card bg-base-100 shadow-xl rounded-xl p-6">
    <div class="flex items-center mb-2">
      <span class="text-2xl font-bold">ゲームクライアントパス設定</span>
      <span class="ml-2 badge badge-outline badge-error">必須</span>
    </div>
    <p class="text-sm text-gray-500 mb-2">
      WorldOfWarships.exeが存在するフォルダを選択してください
    </p>
    {#if $storedPref.installPath}
      <div class="stats shadow w-full mb-2">
        <div class="stat">
          <div class="stat-title">パス</div>
          <div class="stat-value text-lg break-all">
            {$storedPref.installPath}
          </div>
        </div>
      </div>
    {/if}
    {#if $storedGameClientPathError}
      <div class="alert alert-error mb-2">
        <span>{$storedGameClientPathError}</span>
      </div>
    {/if}
    <button class="btn btn-primary w-full" on:click={onClickSelectDirectory}>
      フォルダ選択
    </button>
  </div>

  <!-- 全体表示設定 -->
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
        bind:value={$storedPref.zoomRate}
        on:change={onChange}
      >
        {#each ZOOM_RATES as zr}
          <option selected={zr === $storedPref.zoomRate} value={zr}
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
        bind:value={$storedPref.statsExtra}
        on:change={onChange}
      >
        {#each STATS_EXTRAS as se}
          <option selected={se[0] === $storedPref.statsExtra} value={se[0]}
            >{se[1]}</option
          >
        {/each}
      </select>
    </div>
  </div>

  <!-- カラム別表示設定 -->
  <div class="card bg-base-100 shadow-xl rounded-xl p-6">
    <span class="text-2xl font-bold mb-4">カラム別表示設定</span>
    <!-- プレイヤー名、艦名 -->
    <div
      class="overflow-x-auto rounded-xl shadow border border-base-300 bg-base-200 mb-6"
    >
      <table class="table w-full text-nowrap">
        <thead>
          <tr class="bg-base-300 text-base-content font-semibold">
            <th class="px-4 py-2">カラム名</th>
            <th class="px-4 py-2"></th>
          </tr>
        </thead>
        <tbody>
          <tr class="hover:bg-base-100">
            <td class="px-4 py-2 align-top">プレイヤー名</td>
            <td class="px-4 py-2">
              <div class="flex flex-col gap-2">
                <label class="flex items-center gap-2">
                  <input
                    class="toggle toggle-success"
                    type="checkbox"
                    bind:checked={$storedPref.column.player.enableNationFlag}
                    on:change={onChange}
                  />
                  <span
                    >クラン国籍の国旗を表示する（クラン説明から言語検出）</span
                  >
                </label>
                <label class="flex items-center gap-2">
                  <span>成績に基づく背景色</span>
                  <select
                    class="select select-sm select-bordered"
                    bind:value={$storedPref.column.player.colorPattern}
                    on:change={onChange}
                  >
                    {#each PLAYER_NAME_COLORS as color}
                      <option
                        selected={color[0] ===
                          $storedPref.column.player.colorPattern}
                        value={color[0]}>{color[1]}</option
                      >
                    {/each}
                  </select>
                </label>
              </div>
            </td>
          </tr>
          <tr class="hover:bg-base-100">
            <td class="px-4 py-2 align-top">艦名</td>
            <td class="px-4 py-2">
              <div class="flex flex-col gap-2">
                <label class="flex items-center gap-2">
                  <input
                    class="toggle toggle-success"
                    type="checkbox"
                    bind:checked={$storedPref.column.ship.enableNationFlag}
                    on:change={onChange}
                  />
                  <span>国旗を表示する</span>
                </label>
                <label class="flex items-center gap-2">
                  <input
                    class="toggle toggle-success"
                    type="checkbox"
                    bind:checked={$storedPref.column.ship.isColored}
                    on:change={onChange}
                  />
                  <span>艦種に基づく背景色にする</span>
                </label>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <div
      class="overflow-x-auto rounded-xl shadow border border-base-300 bg-base-200"
    >
      <table class="table w-full text-nowrap">
        <thead>
          <tr class="bg-base-300 text-base-content font-semibold">
            {#each ["カラム名", "艦成績", "総合成績", "小数点以下の桁数"] as columns}
              <th class="text-center px-4 py-2">{columns}</th>
            {/each}
          </tr>
        </thead>
        <tbody>
          {#each STATS_KEYS as statsKey}
            {@const info = STATS_COLUMN_INFO[statsKey]}
            <tr class="hover:bg-base-100">
              <td class="px-4 py-2">{info.full ?? statsKey}</td>
              {#if ["both", "ship"].includes(info.pattern)}
                <td class="text-center px-4 py-2">
                  <input
                    class="toggle toggle-success"
                    type="checkbox"
                    bind:checked={$storedPref.column.stats[statsKey].isShowShip}
                    on:change={onChange}
                  />
                </td>
              {:else}
                <td></td>
              {/if}
              {#if ["both", "overall"].includes(info.pattern)}
                <td class="text-center px-4 py-2">
                  <input
                    class="toggle toggle-success"
                    type="checkbox"
                    bind:checked={
                      $storedPref.column.stats[statsKey].isShowOverall
                    }
                    on:change={onChange}
                  />
                </td>
              {:else}
                <td></td>
              {/if}
              {#if statsKey !== "efficiencyBadge"}
                <td class="text-center px-4 py-2">
                  <select
                    class="select select-sm select-bordered"
                    bind:value={$storedPref.column.stats[statsKey].digit}
                    on:change={onChange}
                  >
                    {#each [0, 1, 2] as digit}
                      <option
                        selected={digit ===
                          $storedPref.column.stats[statsKey].digit}
                        value={digit}>{digit}</option
                      >
                    {/each}
                  </select>
                </td>
              {:else}
                <td></td>
              {/if}
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>

  <!-- その他 -->
  <div class="card bg-base-100 shadow-xl rounded-xl p-6">
    <span class="text-2xl font-bold mb-4">その他</span>
    <ul class="list">
      <li class="list-row flex items-center gap-2">
        <input
          class="toggle toggle-success"
          type="checkbox"
          bind:checked={$storedPref.isSendReport}
          on:change={onChange}
        />
        <span>アプリ改善のためのデータ送信を許可する</span>
      </li>
    </ul>
  </div>
</div>
