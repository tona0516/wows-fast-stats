<script lang="ts">
  import clone from "clone";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
  import { FetchProxy } from "src/lib/FetchProxy";
  import { Notifier } from "src/lib/Notifier";
  import { storedConfig } from "src/stores";
  import { SelectDirectory, StartWatching } from "wailsjs/go/main/App";
  import { data } from "wailsjs/go/models";

  export let inputConfig: data.UserConfigV2;

  let isLoading = false;
  let requiredConfigError: data.RequiredConfigError;

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

      requiredConfigError = await FetchProxy.applyRequiredConfig(
        inputConfig.install_path,
        inputConfig.appid,
      );

      if (!requiredConfigError.valid) {
        return;
      }

      await FetchProxy.getConfig();

      Notifier.success("設定を更新しました");
      StartWatching();
    } catch (error) {
      inputConfig = clone($storedConfig);
      Notifier.failure(error);
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
  {#if requiredConfigError?.install_path}
    <div class="uk-text-danger">
      <UkIcon name="warning" />
      <span class="uk-text-middle">{requiredConfigError.install_path}</span>
    </div>
  {/if}
</div>

<div class="uk-padding-small">
  <input
    class="uk-input"
    type="text"
    placeholder="アプリケーションID"
    bind:value={inputConfig.appid}
  />
  <span>
    <ExternalLink url="https://developers.wargaming.net/"
      >Developer Room</ExternalLink
    >で作成したIDを入力してください。</span
  >
  {#if requiredConfigError?.appid}
    <div class="uk-text-danger">
      <UkIcon name="warning" />
      <span class="uk-text-middle">{requiredConfigError.appid}</span>
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
