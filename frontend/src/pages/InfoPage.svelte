<script lang="ts">
  import ExternalLink from "@components/ExternalLink.svelte";
  import {
    CurrentVersion,
    NewVersion,
    ShowMessageDialog,
  } from "@wails/go/main/App";
  import type { core } from "@wails/go/models";
  import iconApp from "src/assets/images/appicon.png";

  const LINKS = [
    {
      icon: "question-circle",
      url: "https://github.com/tona0516/wows-fast-stats/wiki/FAQ",
      text: "FAQ",
    },
    {
      icon: "twitter-x",
      url: "https://x.com/tonango_wows",
      text: "@tonango_0516",
    },
    {
      icon: "github",
      url: "https://github.com/tona0516/wows-fast-stats",
      text: "tona0516/wows-fast-stats",
    },
  ];

  let checkingUpdate = false;

  const onClickCheckUpdate = async () => {
    if (checkingUpdate) return;

    checkingUpdate = true;

    try {
      const newVersion = (await NewVersion()) as core.NewVersion | null;

      if (newVersion?.version && newVersion.url) {
        await ShowMessageDialog(
          `新しいバージョン ${newVersion.version} が利用可能です。\n${newVersion.url}`,
        );
        return;
      }

      await ShowMessageDialog("現在利用中のバージョンが最新です。");
    } catch (error) {
      await ShowMessageDialog(`アップデートの確認に失敗しました: ${error}`);
    } finally {
      checkingUpdate = false;
    }
  };
</script>

<div class="p-4 flex flex-col items-center">
  <img src={iconApp} alt="" width="128px" height="128px" />
  <div class="pt-1">
    wows-fast-stats {#await CurrentVersion() then semver} {semver} {/await}
  </div>
  <div class="pt-3">
    <button
      class="btn btn-primary"
      on:click={onClickCheckUpdate}
      disabled={checkingUpdate}
    >
      {#if checkingUpdate}
        <span class="loading loading-spinner loading-xs"></span>
        <span class="ml-2">確認中...</span>
      {:else}
        アップデートを確認
      {/if}
    </button>
  </div>
  <div class="pt-2">
    {#each LINKS as link}
      <div class="flex flex-col items-center">
        <ExternalLink url={link.url}>
          <i class="bi bi-{link.icon}">{link.text}</i>
        </ExternalLink>
      </div>
    {/each}
  </div>
  <div class="pt-2">Copyright © 2025 tona0516 All Rights Reserved.</div>
</div>
