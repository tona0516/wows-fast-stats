<script lang="ts">
import clone from "clone";
import { createEventDispatcher } from "svelte";
import { get } from "svelte/store";
import { ApplyUserConfig, StatsPatterns } from "../../wailsjs/go/main/App";
import { WindowReloadApp, LogError } from "../../wailsjs/runtime/runtime";
import { Screenshot } from "../Screenshot";
import {
  storedBattle,
  storedCurrentPage,
  storedIsFirstScreenshot,
  storedUserConfig,
} from "../stores";
import { Func, Page } from "../enums";
import { Const } from "../Const";

const dispatch = createEventDispatcher();

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

function onSwitchPage(page: Page) {
  storedCurrentPage.set(page);
}

async function onClickFunc(func: Func) {
  switch (func) {
    case Func.Reload:
      reload();
      break;
    case Func.Screenshot:
      await screenshot();
      break;
    default:
      break;
  }
}

function reload() {
  WindowReloadApp();
}

async function screenshot() {
  try {
    isLoadingScreenshot = true;
    const screenshot = new Screenshot(battle, isFirstScreenshot);
    await screenshot.manual();
    dispatch("Success", {
      message: "スクリーンショットを保存しました。",
    });
  } catch (error) {
    if (error.message.includes("Canceled")) {
      return;
    }

    LogError(error.name + "," + error.message + "," + error.stack);
    dispatch("Failure", { message: error });
  } finally {
    storedIsFirstScreenshot.set(false);
    isLoadingScreenshot = false;
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
        {#each Const.PAGES as page}
          <button
            type="button"
            class="btn btn-sm btn-outline-secondary m-1 {currentPage ===
              page.name && 'active'}"
            title="{page.title}"
            style="font-size: {userConfig.font_size};"
            on:click="{() => onSwitchPage(page.name)}"
          >
            <i class="{page.iconClass}"></i>
            {page.title}
          </button>
        {/each}
        {#if currentPage == Page.Main}
          {#each Const.FUNCS as func}
            <button
              type="button"
              class="btn btn-sm btn-outline-success m-1"
              title="{func.title}"
              disabled="{func.name === Func.Screenshot &&
                (battle === undefined || isLoadingScreenshot)}"
              style="font-size: {userConfig.font_size};"
              on:click="{() => onClickFunc(func.name)}"
            >
              {#if func.name === Func.Screenshot && isLoadingScreenshot}
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
            {#await StatsPatterns() then statsPatterns}
              {#each statsPatterns as sp}
                {@const label = Const.STATS_PATTERN[sp]}
                <option
                  selected="{sp === userConfig.stats_pattern}"
                  value="{sp}">{label}</option
                >
              {/each}
            {/await}
          </select>
        {/if}
      </div>
    </div>
  </div>
</nav>