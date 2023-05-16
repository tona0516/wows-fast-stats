<script lang="ts">
import { createEventDispatcher } from "svelte";
import { WindowReloadApp } from "../wailsjs/runtime/runtime";
import type { Page } from "./Page";
import type { vo } from "wailsjs/go/models";
import { Screenshot } from "./Screenshot";

const dispatch = createEventDispatcher();

type Func = "reload" | "screenshot";
type NavigationMenu = Page | Func;

export let config: vo.UserConfig;
export let currentPage: Page = "main";
export let battle: vo.Battle;
export let isFirstScreenshot: boolean = true;

let isLoadingScreenshot: boolean = false;

function onClickMenu(menu: NavigationMenu) {
  switch (menu) {
    case "main":
      currentPage = "main";
      break;
    case "config":
      currentPage = "config";
      break;
    case "appinfo":
      currentPage = "appinfo";
      break;
    case "reload":
      WindowReloadApp();
      break;
    case "screenshot":
      isLoadingScreenshot = true;
      const screenshot = new Screenshot(battle, isFirstScreenshot);
      screenshot
        .manual()
        .then(() => {
          dispatch("onScreenshotSuccess", {
            message: "スクリーンショットを保存しました。",
          });
        })
        .catch((error) => {
          // TODO handling based on type
          const errStr = error as string;
          if (errStr.includes("Canceled")) {
            return;
          }
          dispatch("onScreenshotFailure", { message: error });
        })
        .finally(() => {
          isFirstScreenshot = false;
          isLoadingScreenshot = false;
        });
      break;
    default:
      break;
  }
}

const pages: { title: string; name: Page; iconClass: string }[] = [
  { title: "ホーム", name: "main", iconClass: "bi bi-house" },
  { title: "設定", name: "config", iconClass: "bi bi-gear" },
  { title: "アプリ情報", name: "appinfo", iconClass: "bi bi-info-circle" },
];

const funcs: { title: string; name: Func; iconClass: string }[] = [
  { title: "リロード", name: "reload", iconClass: "bi bi-arrow-clockwise" },
  {
    title: "スクリーンショット",
    name: "screenshot",
    iconClass: "bi bi-camera",
  },
];
</script>

<nav class="navbar navbar-expand-sm navbar-light bg-light">
  <div class="container-fluid">
    <button
      class="navbar-toggler"
      type="button"
      data-bs-toggle="collapse"
      data-bs-target="#navbarNavAltMarkup"
      aria-controls="navbarNavAltMarkup"
      aria-expanded="false"
      aria-label="Toggle navigation"
      style="font-size: {config?.font_size || 'medium'};"
    >
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
      <div class="navbar-nav">
        {#each pages as page}
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary m-1 {currentPage ===
              page.name && 'active'}"
            title="{page.title}"
            style="font-size: {config?.font_size || 'medium'};"
            on:click="{() => onClickMenu(page.name)}"
          >
            <i class="{page.iconClass}"></i>
            {page.title}
          </button>
        {/each}
        {#if currentPage == "main"}
          {#each funcs as func}
            <button
              type="button"
              class="btn btn-sm btn-outline-success m-1"
              title="{func.title}"
              disabled="{func.name === 'screenshot' &&
                (battle === undefined || isLoadingScreenshot)}"
              style="font-size: {config?.font_size || 'medium'};"
              on:click="{() => onClickMenu(func.name)}"
            >
              {#if func.name === "screenshot" && isLoadingScreenshot}
                <span
                  class="spinner-border spinner-border-sm"
                  role="status"
                  aria-hidden="true"></span>
                読み込み中...
              {:else}
                <i class="{func.iconClass}"></i>
                {func.title}
              {/if}
            </button>
          {/each}
        {/if}
      </div>
    </div>
  </div>
</nav>
