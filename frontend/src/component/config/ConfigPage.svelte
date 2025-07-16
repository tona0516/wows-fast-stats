<script lang="ts">
  import { DispName } from "src/lib/DispName";
  import { deriveColumnSettings } from "src/lib/util";
  import { storedConfig, storedInstallPathError } from "src/stores";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import {
    OpenDirectory,
    SelectDirectory,
    ShowMessageDialog,
    UpdateInstallPath,
    UpdateUserConfig,
  } from "wailsjs/go/main/App";

  $: inputConfig = $storedConfig;
  $: installPath = $storedConfig.install_path;
  $: columnSettings = deriveColumnSettings(inputConfig);
  $: applyInstallPathButton = {
    label: "フォルダ設定保存",
    disabled: false,
  };

  onMount(() => {
    themeChange(false);
  });

  const onClickSelectDirectory = async () => {
    SelectDirectory().then((path) => {
      if (!path) {
        return;
      }

      installPath = path;
      applyInstallPathButton = {
        label: "フォルダ設定保存",
        disabled: false,
      };
    });
  };

  const onClickSaveInstallPath = async () => {
    try {
      await UpdateInstallPath(installPath);
      storedInstallPathError.set("");
      applyInstallPathButton = {
        label: "保存しました！",
        disabled: true,
      };
    } catch (error) {
      storedInstallPathError.set(error as string);
    }
  };

  const onClickOpenDirectory = async (path: string) => {
    OpenDirectory(path).catch((error) => {
      if (error instanceof Error) {
        ShowMessageDialog(error.message);
      }
    });
  };

  const onChange = async () => {
    try {
      await UpdateUserConfig(inputConfig);
    } catch (error) {
      inputConfig = $storedConfig;
      if (error instanceof Error) {
        ShowMessageDialog(error.message);
      }
    }
  };
</script>

<div>
  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">インストールフォルダ設定</p>
    <p class="text-sm text-nowrap">
      ゲームクライアントの実行ファイルがあるフォルダを設定してください。
    </p>

    <div class="stats shadow">
      <div class="stat">
        <div class="stat-title">現在の入力値（設定値）</div>
        <div class="stat-value text-lg">{inputConfig.install_path}</div>
      </div>
    </div>

    <div>
      <button class="btn btn-neutral" on:click={onClickSelectDirectory}
        >フォルダ選択</button
      >
      <button
        class="btn btn-primary ml-2"
        on:click={onClickSaveInstallPath}
        disabled={applyInstallPathButton.disabled}
        >{applyInstallPathButton.label}</button
      >
    </div>
    {#if $storedInstallPathError}
      <div>
        <i class="bi bi-warning">{$storedInstallPathError}</i>
      </div>
    {/if}
  </div>

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">統計パターン</p>
    <select
      class="select my-2"
      bind:value={inputConfig.stats_pattern}
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
      {#each ["light", "dark", "retro", "night"] as theme}
        <option value={theme}>{theme}</option>
      {/each}
    </select>
  </div>

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">UIサイズ</p>
    <select
      class="select my-2"
      bind:value={inputConfig.font_size}
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
                  bind:checked={inputConfig.display.ship[column.ship.key]}
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
                  bind:checked={inputConfig.display.overall[column.overall.key]}
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
                  bind:value={inputConfig.digit[column.digit.key]}
                  on:change={onChange}
                >
                  {#each [0, 1, 2] as digit}
                    <option
                      selected={digit === inputConfig.digit[column.digit.key]}
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
    <p class="text-xl font-bold">各種カラー</p>
    <div class="grid xl:grid-cols-3 gap-4">
      <div>
        <table class="table max-w-xs text-nowrap my-2">
          <thead>
            <tr>
              {#each ["スキル", "文字色"] as column}
                <th class="text-center">{column}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each DispName.SKILL_LEVELS.toArray() as sl}
              <tr>
                <td class="text-center">{sl.value}</td>
                <td>
                  <input
                    class="input w-[64px]"
                    type="color"
                    bind:value={inputConfig.color.skill.text[sl.key]}
                    on:change={onChange}
                  />
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      <div>
        <table class="table max-w-xs text-nowrap my-2">
          <thead>
            <tr>
              {#each ["Tier", "使用艦", "非使用艦"] as column}
                <th class="text-center">{column}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each DispName.TIER_GROUPS.toArray() as tg}
              <tr>
                <td class="text-center">{tg.value}</td>
                <td>
                  <input
                    class="input w-[64px]"
                    type="color"
                    bind:value={inputConfig.color.tier.own[tg.key]}
                    on:change={onChange}
                  />
                </td>

                <td>
                  <input
                    class="input w-[64px]"
                    type="color"
                    bind:value={inputConfig.color.tier.other[tg.key]}
                    on:change={onChange}
                  />
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
      <div>
        <table class="table max-w-xs text-nowrap my-2">
          <thead>
            <tr>
              {#each ["艦種", "使用艦", "非使用艦"] as column}
                <th class="text-center">{column}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each DispName.SHIP_TYPES.toArray() as st}
              <tr>
                <td class="text-center">{st.value}</td>
                <td>
                  <input
                    class="input w-[64px]"
                    type="color"
                    bind:value={inputConfig.color.ship_type.own[st.key]}
                    on:change={onChange}
                  />
                </td>

                <td>
                  <input
                    class="input w-[64px]"
                    type="color"
                    bind:value={inputConfig.color.ship_type.other[st.key]}
                    on:change={onChange}
                  />
                </td>
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
  </div>

  <div class="p-4 flex flex-col items-center">
    <p class="text-xl font-bold">プレイヤー名の背景色</p>
    <select
      class="select my-2"
      bind:value={inputConfig.color.player_name}
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
          bind:checked={inputConfig.show_language_frag}
          on:change={onChange}
        />
        クラン国籍を表示する（クラン説明から言語検出）
      </li>
      <li class="list-row">
        <input
          class="toggle toggle-success"
          type="checkbox"
          bind:checked={inputConfig.send_report}
          on:change={onChange}
        />アプリ改善のためのデータ送信を許可する
      </li>
      <li class="list-row">
        <input
          class="toggle toggle-success"
          type="checkbox"
          bind:checked={inputConfig.save_temp_arena_info}
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
