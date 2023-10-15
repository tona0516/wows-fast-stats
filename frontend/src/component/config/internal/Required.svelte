<script lang="ts">
  import clone from "clone";
  import ExternalLink from "src/component/common/ExternalLink.svelte";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
  import { storedRequiredConfigError, storedConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import {
    ApplyRequiredUserConfig,
    SelectDirectory,
    UserConfig,
  } from "wailsjs/go/main/App";
  import { domain, vo } from "wailsjs/go/models";

  export let inputConfig: domain.UserConfig;

  let isLoading = false;
  let requiredConfigError: vo.RequiredConfigError;

  const dispatch = createEventDispatcher();

  const clickSelectDirectory = async () => {
    try {
      const path = await SelectDirectory();
      if (!path) return;
      inputConfig.install_path = path;
    } catch (error) {
      dispatch("Failure", { message: error });
    }
  };

  const clickApply = async () => {
    try {
      isLoading = true;
      requiredConfigError = await ApplyRequiredUserConfig(
        inputConfig.install_path,
        inputConfig.appid,
      );

      if (!requiredConfigError.valid) {
        return;
      }

      const latest = await UserConfig();
      storedConfig.set(latest);
      storedRequiredConfigError.set(requiredConfigError);
      dispatch("UpdateSuccess");
    } catch (error) {
      inputConfig = clone($storedConfig);
      dispatch("Failure", { message: error });
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
