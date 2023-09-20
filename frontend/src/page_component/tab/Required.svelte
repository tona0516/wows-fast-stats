<script lang="ts">
  import clone from "clone";
  import ExternalLink from "src/other_component/ExternalLink.svelte";
  import { storedUserConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import {
    ApplyRequiredUserConfig,
    SelectDirectory,
    UserConfig,
  } from "wailsjs/go/main/App";
  import { domain, vo } from "wailsjs/go/models";

  export let inputUserConfig: domain.UserConfig;

  let isLoading = false;
  let validatedResult: vo.ValidatedResult;

  const dispatch = createEventDispatcher();

  const clickSelectDirectory = async () => {
    try {
      const path = await SelectDirectory();
      if (!path) return;
      inputUserConfig.install_path = path;
    } catch (error) {
      dispatch("Failure", { message: error });
    }
  };

  const clickApply = async () => {
    try {
      isLoading = true;
      validatedResult = await ApplyRequiredUserConfig(
        inputUserConfig.install_path,
        inputUserConfig.appid,
      );

      const results = Object.values(validatedResult);
      const isValid = results.filter((it) => !it).length === results.length;
      if (!isValid) {
        return;
      }

      const latest = await UserConfig();
      storedUserConfig.set(latest);
      dispatch("UpdateSuccess");
    } catch (error) {
      inputUserConfig = clone($storedUserConfig);
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
      bind:value={inputUserConfig.install_path}
    />
    <button
      class="uk-button uk-button-default uk-text-nowrap"
      on:click={clickSelectDirectory}>フォルダ選択</button
    >
  </div>
  <span>ゲームクライアントの実行ファイルがあるフォルダを選択してください。</span
  >
  {#if validatedResult?.install_path}
    <div class="uk-text-danger">
      <span uk-icon="icon: warning" />
      <span class="uk-text-middle">{validatedResult.install_path}</span>
    </div>
  {/if}
</div>

<div class="uk-padding-small">
  <input
    class="uk-input"
    type="text"
    placeholder="アプリケーションID"
    bind:value={inputUserConfig.appid}
  />
  <span>
    <ExternalLink url="https://developers.wargaming.net/"
      >Developer Room</ExternalLink
    >で作成したIDを入力してください。</span
  >
  {#if validatedResult?.appid}
    <div class="uk-text-danger">
      <span uk-icon="icon: warning"></span>
      <span class="uk-text-middle">{validatedResult.appid}</span>
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
        <div uk-spinner />
      {:else}
        保存
      {/if}
    </button>
  </div>
</div>
