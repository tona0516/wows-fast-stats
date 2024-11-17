<script lang="ts">
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
  import { Notifier } from "src/lib/Notifier";
  import { storedConfig, storedInstallPathError } from "src/stores";
  import {
    SelectDirectory,
    StartWatching,
    UpdateInstallPath,
  } from "wailsjs/go/main/App";

  let isLoading = false;
  let inputInstallPathError: unknown = $storedInstallPathError;

  $: inputConfig = $storedConfig;

  const clickSelectDirectory = async () => {
    try {
      const path = await SelectDirectory();
      if (!path) return;
      inputConfig.install_path = path;
    } catch (error) {
      Notifier.failure(error);
    }
  };

  const clickApply = async () => {
    try {
      isLoading = true;
      await UpdateInstallPath(inputConfig.install_path);
      Notifier.success("設定を更新しました");
      StartWatching();
    } catch (error) {
      inputInstallPathError = error;
    } finally {
      isLoading = false;
    }
  };
</script>

<div class="uk-padding-small">
  <div class="uk-flex">
    <input
      class="uk-input"
      type="text"
      placeholder="World of Warshipsインストールフォルダ"
      bind:value={inputConfig.install_path}
    />
    <button
      class="uk-button uk-button-default uk-text-nowrap"
      on:click={clickSelectDirectory}>フォルダ選択</button
    >
  </div>
  <span>ゲームクライアントの実行ファイルがあるフォルダを選択してください。</span
  >
  {#if inputInstallPathError}
    <div class="uk-text-danger">
      <UkIcon name="warning" />
      <span class="uk-text-middle">{inputInstallPathError}</span>
    </div>
  {/if}
</div>

<div class="uk-padding-small">
  <div class="uk-flex">
    <button
      class="uk-button uk-button-primary uk-text-nowrap"
      disabled={isLoading}
      on:click={clickApply}
    >
      {#if isLoading}
        <UkSpinner />
      {:else}
        保存
      {/if}
    </button>
  </div>
</div>
