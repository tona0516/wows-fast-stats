<script lang="ts">
import { createEventDispatcher } from "svelte";
import { WindowReloadApp } from "../wailsjs/runtime/runtime";
import type { Page } from "./Page";
import { Screenshot } from "./Screenshot";
import { ApplyUserConfig, LogError } from "../wailsjs/go/main/App";
import { get } from "svelte/store";
import {
  storedBattle,
  storedCurrentPage,
  storedIsFirstScreenshot,
  storedUserConfig,
} from "./stores";
import clone from "clone";

const dispatch = createEventDispatcher();

type Func = "reload" | "screenshot";
type NavigationMenu = Page | Func;

let battle = get(storedBattle);
storedBattle.subscribe((it) => (battle = it));

let currentPage = get(storedCurrentPage);
storedCurrentPage.subscribe((it) => (currentPage = it));

let isFirstScreenshot = get(storedIsFirstScreenshot);
storedIsFirstScreenshot.subscribe((it) => (isFirstScreenshot = it));

let userConfig = get(storedUserConfig);
storedUserConfig.subscribe((it) => (userConfig = it));

let isLoadingScreenshot: boolean = false;
let selectedStatsPattern: string;

function onClickMenu(menu: NavigationMenu) {
  switch (menu) {
    case "main":
      storedCurrentPage.set("main");
      break;
    case "config":
      storedCurrentPage.set("config");
      break;
    case "appinfo":
      storedCurrentPage.set("appinfo");
      break;
    case "alert_player":
      storedCurrentPage.set("alert_player");
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
          dispatch("Success", {
            message: "スクリーンショットを保存しました。",
          });
        })
        .catch((error: Error) => {
          if (error.message.includes("Canceled")) {
            return;
          }

          LogError(error.name + "," + error.message + "," + error.stack);
          dispatch("Failure", { message: error });
        })
        .finally(() => {
          storedIsFirstScreenshot.set(false);
          isLoadingScreenshot = false;
        });
      break;
    default:
      break;
  }
}

async function onStatsPatternChanged() {
  try {
    const config = clone(userConfig);
    config.stats_pattern = selectedStatsPattern;

    await ApplyUserConfig(config);
    storedUserConfig.set(config);

    dispatch("ChangeStatsPattern");
  } catch (error) {
    dispatch("Failure", { message: error });
    return;
  }
}

const pages: { title: string; name: Page; iconClass: string }[] = [
  { title: "ホーム", name: "main", iconClass: "bi bi-house" },
  { title: "設定", name: "config", iconClass: "bi bi-gear" },
  { title: "アプリ情報", name: "appinfo", iconClass: "bi bi-info-circle" },
  {
    title: "プレイヤーリスト",
    name: "alert_player",
    iconClass: "bi bi-person-lines-fill",
  },
];

const funcs: { title: string; name: Func; iconClass: string }[] = [
  { title: "リロード", name: "reload", iconClass: "bi bi-arrow-clockwise" },
  {
    title: "スクリーンショット",
    name: "screenshot",
    iconClass: "bi bi-camera",
  },
];

// TODO get from backend
const statsPattern: { titile: string; pattern: string }[] = [
  { titile: "ランダム戦", pattern: "pvp_all" },
  { titile: "ランダム戦(ソロ)", pattern: "pvp_solo" },
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
      style="font-size: {userConfig.font_size};"
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
            style="font-size: {userConfig.font_size};"
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
              style="font-size: {userConfig.font_size};"
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
          <select
            class="form-select form-select-sm m-1"
            style="font-size: {userConfig.font_size};"
            bind:value="{selectedStatsPattern}"
            on:change="{onStatsPatternChanged}"
          >
            {#each statsPattern as sp}
              {#if sp.pattern === userConfig.stats_pattern}
                <option selected value="{sp.pattern}">{sp.titile}</option>
              {:else}
                <option value="{sp.pattern}">{sp.titile}</option>
              {/if}
            {/each}
          </select>
        {/if}
      </div>
    </div>
  </div>
</nav>
