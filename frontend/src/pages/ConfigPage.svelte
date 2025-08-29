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
    storedPlayerNameColumnSettings,
    storedShipInfoColumnSettings,
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

<div class="w-3/4">
  <div class="grid grid-cols-1 gap-4">
    <!-- インストールパス設定 -->
    <div class="grid grid-cols-1 gap-2">
      <div class="flex items-center">
        <span class="text-xl font-bold">インストールパス設定</span>
        <span class="ml-1 badge badge-outline badge-error">必須</span>
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
        <div role="alert" class="alert alert-error alert-soft">
          <span>{$storedInstallPathError}</span>
        </div>
      {/if}

      <div>
        <button class="btn btn-neutral" on:click={onClickSelectDirectory}
          >フォルダ選択</button
        >
      </div>
    </div>

    <!-- 全体表示設定 -->
    <div class="grid grid-cols-1 gap-2">
      <span class="text-xl font-bold">全体表示設定</span>
      <div>
        <p class="font-bold">テーマ</p>
        <select class="select" data-choose-theme>
          {#each Theme.getAll() as theme}
            <option value={theme}>{theme}</option>
          {/each}
        </select>
      </div>

      <div>
        <p class="font-bold">UIサイズ</p>
        <select class="select" bind:value={$storedZoomRate}>
          {#each ZOOM_RATES as zr}
            <option selected={zr === $storedZoomRate} value={zr}>{zr}%</option>
          {/each}
        </select>
      </div>

      <div>
        <p class="font-bold">統計パターン</p>
        <select class="select my-2" bind:value={$storedStatsExtra}>
          {#each STATS_EXTRAS as se}
            <option selected={se[0] === $storedStatsExtra} value={se[0]}
              >{se[1]}</option
            >
          {/each}
        </select>
      </div>
    </div>
  </div>

  <!-- カラム別表示設定 -->
  <div class="grid grid-cols-1 gap-2">
    <span class="text-xl font-bold">カラム別表示設定</span>

    <!-- プレイヤー名、艦名 -->
    <table class="table max-w-md text-nowrap">
      <thead>
        <tr>
          <th>カラム名</th>
          <th></th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>プレイヤー名</td>
          <td>
            <div class="grid grid-cols-1 gap-2">
              <div>
                クラン国籍の国旗を表示する（クラン説明から言語検出）
                <input
                  class="toggle toggle-success"
                  type="checkbox"
                  bind:checked={
                    $storedPlayerNameColumnSettings.enableNationFlag
                  }
                />
              </div>
              <div>
                成績に基づく背景色
                <select
                  class="select select-sm"
                  bind:value={$storedPlayerNameColumnSettings.colorPattern}
                >
                  {#each PLAYER_NAME_COLORS as color}
                    <option
                      selected={color[0] ===
                        $storedPlayerNameColumnSettings.colorPattern}
                      value={color[0]}>{color[1]}</option
                    >
                  {/each}
                </select>
              </div>
            </div>
          </td>
        </tr>
        <tr>
          <td>艦名</td>
          <td>
            <div class="grid grid-cols-1 gap-2">
              <div>
                国旗を表示する
                <input
                  class="toggle toggle-success"
                  type="checkbox"
                  bind:checked={$storedShipInfoColumnSettings.enableNationFlag}
                />
              </div>
              <div>
                艦種に基づく背景色にする
                <input
                  class="toggle toggle-success"
                  type="checkbox"
                  bind:checked={$storedShipInfoColumnSettings.enableColorized}
                />
              </div>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

    <table class="table max-w-md text-nowrap">
      <thead>
        <tr>
          {#each ["カラム名", "艦成績", "総合成績", "小数点以下の桁数"] as columns}
            <th class="text-center">{columns}</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each STATS_KEYS as statsKey}
          {@const info = STATS_COLUMN_INFO[statsKey]}
          <tr>
            <td>
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

  <!-- その他 -->
  <div class="grid grid-cols-1 gap-2">
    <span class="text-xl font-bold">その他</span>
    <ul class="list">
      <li class="list-row">
        アプリ改善のためのデータ送信を許可する
        <input
          class="toggle toggle-success"
          type="checkbox"
          bind:checked={sendReport}
          on:change={() => UpdateSendReport(sendReport)}
        />
      </li>
    </ul>
  </div>
</div>
