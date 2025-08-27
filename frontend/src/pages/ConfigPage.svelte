<script lang="ts">
  import {
    STATS_EXTRAS,
    ZOOM_RATES,
    STATS_KEYS,
    STATS_COLUMN_INFO,
    PLAYER_NAME_COLORS,
  } from "@libs/constants";
  import {
    storedInstallPathError,
    showToast,
    storedStatsExtra,
    storedZoomRate,
    storedColumnmSettings,
  } from "@libs/stores";
  import { Theme } from "@libs/Theme";
  import {
    SendReport,
    InstallPath,
    TrySaveInstallPath,
    SubscribeBattle,
    UpdateSendReport,
  } from "@wails/go/main/App";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";

  let installPath: string = "";
  let sendReport: boolean = false;

  onMount(async () => {
    themeChange(false);
    sendReport = await SendReport();
    installPath = await InstallPath();
  });

  const onClickSelectDirectory = async () => {
    try {
      const isSuccess = await TrySaveInstallPath();
      if (isSuccess) {
        installPath = await InstallPath();
        storedInstallPathError.set("");

        showToast("インストールパスを設定しました");
        SubscribeBattle();
      }
    } catch (error) {
      storedInstallPathError.set(error as string);
    }
  };
</script>

<div>
  <div class="p-4 flex flex-col items-center">
    <div class="flex items-center">
      <span class="text-xl font-bold">インストールパス設定</span>
      <span class="ml-2 badge badge-outline badge-error">必須</span>
    </div>

    <p class="text-sm text-nowrap">
      WorldOfWarships.exeが存在するフォルダを選択してください
    </p>

    {#if installPath}
      <div class="stats shadow w-3/4">
        <div class="stat">
          <div class="stat-title">ゲームクライアント インストールパス</div>
          <div class="stat-value text-lg">{installPath}</div>
        </div>
      </div>
    {/if}

    {#if $storedInstallPathError}
      <div role="alert" class="mt-2 alert alert-error alert-soft">
        <span>{$storedInstallPathError}</span>
      </div>
    {/if}

    <div class="mt-2">
      <button class="btn btn-neutral" on:click={onClickSelectDirectory}
        >フォルダ選択</button
      >
    </div>
  </div>

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">統計パターン</p>
    <select class="select my-2" bind:value={$storedStatsExtra}>
      {#each STATS_EXTRAS as se}
        <option selected={se[0] === $storedStatsExtra} value={se[0]}
          >{se[1]}</option
        >
      {/each}
    </select>
  </div>

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">テーマ</p>
    <select class="select my-2" data-choose-theme>
      {#each Theme.getAll() as theme}
        <option value={theme}>{theme}</option>
      {/each}
    </select>
  </div>

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">UIサイズ</p>
    <select class="select my-2" bind:value={$storedZoomRate}>
      {#each ZOOM_RATES as zr}
        <option selected={zr === $storedZoomRate} value={zr}>{zr}%</option>
      {/each}
    </select>
  </div>

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">表示項目</p>
    <table class="table max-w-md text-nowrap my-2">
      <thead>
        <tr>
          {#each ["項目", "艦成績", "総合成績", "小数点以下の桁数"] as columns}
            <th class="text-center">{columns}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each STATS_KEYS as statsKey}
          {@const info = STATS_COLUMN_INFO[statsKey]}
          <tr>
            <td class="text-center">
              {info.full ?? statsKey}
            </td>

            {#if ["both", "ship"].includes(info.pattern)}
              <td>
                <input
                  class="toggle toggle-success"
                  type="checkbox"
                  bind:checked={$storedColumnmSettings[statsKey].ship}
                />
              </td>
            {:else}
              <td></td>
            {/if}

            {#if ["both", "overall"].includes(info.pattern)}
              <td>
                <input
                  class="toggle toggle-success"
                  type="checkbox"
                  bind:checked={$storedColumnmSettings[statsKey].overall}
                />
              </td>
            {:else}
              <td></td>
            {/if}

            {#if statsKey !== "ship_badge"}
              <td>
                <select
                  class="select select-sm"
                  bind:value={$storedColumnmSettings[statsKey].digit}
                >
                  {#each [0, 1, 2] as digit}
                    <option
                      selected={digit ===
                        $storedColumnmSettings[statsKey].digit}
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

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">その他</p>
    <ul class="list my-2">
      <li class="list-row">
        <input
          class="toggle toggle-success"
          type="checkbox"
          bind:checked={sendReport}
          on:change={() => UpdateSendReport(sendReport)}
        />アプリ改善のためのデータ送信を許可する
      </li>
    </ul>
  </div>
</div>
