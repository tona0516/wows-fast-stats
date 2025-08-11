<script lang="ts">
  import { DispName } from "src/lib/DispName";
  import { Theme } from "src/lib/Theme";
  import { deriveColumnSettings } from "src/lib/util";
  import { showToast, storedConfig, storedInstallPathError } from "src/stores";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import {
    OpenDirectory,
    ShowMessageDialog,
    SubscribeBattle,
    TrySaveInstallPath,
    UpdateUserConfig,
  } from "wailsjs/go/main/App";

  $: columnSettings = deriveColumnSettings($storedConfig);

  onMount(() => {
    themeChange(false);
  });

  const onClickSelectDirectory = async () => {
    try {
      const isSuccess = await TrySaveInstallPath();
      if (isSuccess) {
        showToast("インストールパスを設定しました");
        storedInstallPathError.set("");
        SubscribeBattle();
      }
    } catch (error) {
      storedInstallPathError.set(error as string);
    }
  };

  const onClickOpenDirectory = async (path: string) => {
    OpenDirectory(path).catch((error) => {
      ShowMessageDialog(error as string);
    });
  };

  const onChange = async () => {
    const beforeConfig = structuredClone($storedConfig);

    try {
      await UpdateUserConfig($storedConfig);
    } catch (error) {
      storedConfig.set(beforeConfig);
      ShowMessageDialog(error as string);
    }
  };
</script>

<div>
  <div class="p-4 flex flex-col items-center">
    <div class="flex items-center">
      <span class="text-xl font-bold">インストールパス設定</span>
      <span class="ml-2 badge badge-outline badge-error">必須</span>
    </div>

    <div class="stats shadow w-3/4">
      <div class="stat {$storedInstallPathError && 'input-error'}">
        <div class="stat-title">ゲームクライアント インストールパス</div>
        <div class="stat-value text-lg">{$storedConfig.install_path}</div>
      </div>
    </div>

    <p class="text-sm text-nowrap">
      WorldOfWarships.exeが存在するフォルダを選択してください
    </p>

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
    <select
      class="select my-2"
      bind:value={$storedConfig.stats_pattern}
      on:change={onChange}
    >
      {#each DispName.STATS_PATTERNS.toArray() as sp}
        <option selected={sp.key == $storedConfig.stats_pattern} value={sp.key}
          >{sp.value}</option
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
    <select
      class="select my-2"
      bind:value={$storedConfig.font_size}
      on:change={onChange}
    >
      {#each DispName.FONT_SIZES.toArray() as fs}
        <option selected={fs.key === $storedConfig.font_size} value={fs.key}
          >{fs.value}</option
        >
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
        {#each columnSettings as column}
          <tr>
            <td class="text-center">
              {DispName.FULL_COLUMN_NAMES.get(column.key) ?? column.key}
            </td>

            {#if column.ship.key}
              <td>
                <input
                  class="toggle toggle-success"
                  type="checkbox"
                  bind:checked={$storedConfig.display.ship[column.ship.key]}
                  on:change={onChange}
                />
              </td>
            {:else}
              <td></td>
            {/if}

            {#if column.overall.key}
              <td>
                <input
                  class="toggle toggle-success"
                  type="checkbox"
                  bind:checked={
                    $storedConfig.display.overall[column.overall.key]
                  }
                  on:change={onChange}
                />
              </td>
            {:else}
              <td></td>
            {/if}

            {#if column.digit.key}
              <td>
                <select
                  class="select select-sm"
                  bind:value={$storedConfig.digit[column.digit.key]}
                  on:change={onChange}
                >
                  {#each [0, 1, 2] as digit}
                    <option
                      selected={digit === $storedConfig.digit[column.digit.key]}
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
    <p class="text-xl font-bold">プレイヤー名の背景色</p>
    <select
      class="select my-2"
      bind:value={$storedConfig.color.player_name}
      on:change={onChange}
    >
      {#each DispName.PLAYER_NAME_COLORS.toArray() as pnc}
        <option
          selected={pnc.key === $storedConfig.color.player_name}
          value={pnc.key}>{pnc.value}</option
        >
      {/each}
    </select>
  </div>

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">その他</p>
    <ul class="list my-2">
      <li class="list-row">
        <input
          class="toggle toggle-success"
          type="checkbox"
          bind:checked={$storedConfig.show_language_frag}
          on:change={onChange}
        />
        クラン国籍を表示する（クラン説明から言語検出）
      </li>
      <li class="list-row">
        <input
          class="toggle toggle-success"
          type="checkbox"
          bind:checked={$storedConfig.send_report}
          on:change={onChange}
        />アプリ改善のためのデータ送信を許可する
      </li>
      <li class="list-row">
        <input
          class="toggle toggle-success"
          type="checkbox"
          bind:checked={$storedConfig.save_temp_arena_info}
          on:change={onChange}
        />
        【開発用】自動で戦闘情報(tempArenaInfo.json)を保存する

        <!-- svelte-ignore a11y-invalid-attribute -->
        <a href="#" on:click={() => onClickOpenDirectory("temp_arena_info")}>
          <i class="bi bi-folder">保存フォルダを開く</i>
        </a>
      </li>
    </ul>
  </div>
</div>
