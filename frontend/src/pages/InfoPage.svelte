<script lang="ts">
  import ExternalLink from "@components/ExternalLink.svelte";
  import { CurrentVersion, NewVersion } from "@wails/go/main/App";
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
  let updateMessage = "";
  let updateLink = "";
  let updateStatus: "idle" | "checking" | "success" | "warning" | "error" =
    "idle";

  $: statusDotClass =
    updateStatus === "warning"
      ? "bg-warning"
      : updateStatus === "success"
        ? "bg-success"
        : updateStatus === "error"
          ? "bg-error"
          : updateStatus === "checking"
            ? "bg-info"
            : "bg-base-content/20";

  const onClickCheckUpdate = async () => {
    if (checkingUpdate) return;

    checkingUpdate = true;
    updateStatus = "checking";
    updateMessage = "最新バージョンを確認しています...";
    updateLink = "";

    try {
      const newVersion = (await NewVersion()) as core.NewVersion | null;

      if (newVersion?.version && newVersion.downloadURL) {
        updateStatus = "warning";
        updateMessage = `新しいバージョン ${newVersion.version} が利用可能です。`;
        updateLink = newVersion.downloadURL;
        return;
      }

      updateStatus = "success";
      updateMessage = "現在利用中のバージョンが最新です。";
    } catch (error) {
      updateStatus = "error";
      updateMessage = `アップデートの確認に失敗しました: ${error}`;
    } finally {
      checkingUpdate = false;
    }
  };
</script>

<div class="min-h-full w-full bg-base-200/40 p-6">
  <div class="mx-auto flex w-full max-w-3xl flex-col gap-6">
    <div class="rounded-2xl bg-base-100 p-6 shadow-md">
      <div class="flex flex-col items-center gap-3">
        <img src={iconApp} alt="" width="128px" height="128px" />
        <div class="text-lg font-semibold tracking-wide">
          wows-fast-stats
          <span class="ml-2 text-base font-normal text-base-content/70">
            {#await CurrentVersion() then semver} {semver} {/await}
          </span>
        </div>
      </div>
    </div>

    <div class="grid gap-4 md:grid-cols-[1.2fr,0.8fr]">
      <div class="rounded-2xl bg-base-100 p-6 shadow-md">
        <div class="flex items-center justify-between">
          <div>
            <div class="text-base font-semibold">アップデート</div>
            <div class="text-sm text-base-content/60">バージョンの確認</div>
          </div>
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

        <div class="mt-4 rounded-xl border border-base-300 bg-base-200/50 p-4">
          <div class="flex items-start gap-3">
            <div
              class={`mt-1 h-2.5 w-2.5 rounded-full ${statusDotClass}`}
            ></div>
            <div class="space-y-1">
              <div class="text-sm font-medium">チェック結果</div>
              <div class="text-sm text-base-content/70">
                {#if updateMessage}
                  {updateMessage}
                {:else}
                  まだ確認していません。
                {/if}
              </div>
              {#if updateLink}
                <div class="pt-1 text-sm">
                  <ExternalLink url={updateLink}>
                    <i class="bi bi-download">ダウンロードページへ</i>
                  </ExternalLink>
                </div>
              {/if}
            </div>
          </div>
        </div>
      </div>

      <div class="rounded-2xl bg-base-100 p-6 shadow-md">
        <div class="text-base font-semibold">リンク</div>
        <div class="mt-3 flex flex-col gap-3">
          {#each LINKS as link}
            <ExternalLink url={link.url}>
              <div
                class="flex items-center gap-2 rounded-lg border border-base-300 px-3 py-2"
              >
                <i class="bi bi-{link.icon}"></i>
                <span class="text-sm">{link.text}</span>
              </div>
            </ExternalLink>
          {/each}
        </div>
      </div>
    </div>

    <div class="text-center text-xs text-base-content/60">
      Copyright © 2026 tona0516 All Rights Reserved.
    </div>
  </div>
</div>
