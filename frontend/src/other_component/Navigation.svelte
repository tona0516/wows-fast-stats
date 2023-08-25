<script lang="ts">
  import clone from "clone";
  import { type NavigationItem, Const } from "src/Const";
  import type { Screenshot } from "src/Screenshot";
  import { Page, Func } from "src/enums";
  import {
    storedCurrentPage,
    storedUserConfig,
    storedBattle,
  } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import { Button, Spinner } from "sveltestrap";
  import { ApplyUserConfig, StatsPatterns } from "wailsjs/go/main/App";
  import { WindowReloadApp } from "wailsjs/runtime/runtime";

  export let screenshot: Screenshot;

  const dispatch = createEventDispatcher();

  let isScreenshotting: boolean = false;
  let selectedStatsPattern: string;

  function onSwitchPage(item: NavigationItem<Page>) {
    storedCurrentPage.set(item.name);
  }

  async function onClickFunc(item: NavigationItem<Func>) {
    switch (item.name) {
      case Func.Reload:
        reload();
        break;
      case Func.Screenshot:
        await takeScreenshot();
        break;
    }
  }

  function reload() {
    WindowReloadApp();
  }

  async function takeScreenshot() {
    try {
      isScreenshotting = true;
      if (await screenshot.manual($storedBattle.meta)) {
        dispatch("Success", { message: "スクリーンショットを保存しました。" });
      }
    } catch (error) {
      dispatch("Failure", { message: error });
    } finally {
      isScreenshotting = false;
    }
  }

  async function onStatsPatternChanged() {
    try {
      const config = clone($storedUserConfig);
      config.stats_pattern = selectedStatsPattern;

      await ApplyUserConfig(config);
      storedUserConfig.set(config);
    } catch (error) {
      dispatch("Failure", { message: error });
      return;
    }
  }
</script>

<nav class="navbar navbar-expand-sm sticky-top navbar-light bg-light p-1">
  <div class="container-fluid">
    <button
      class="navbar-toggler m-1"
      type="button"
      data-bs-toggle="collapse"
      data-bs-target="#navbarNavAltMarkup"
      aria-controls="navbarNavAltMarkup"
      aria-expanded="false"
      aria-label="Toggle navigation"
      style="font-size: {$storedUserConfig.font_size};"
    >
      <span class="navbar-toggler-icon" />
    </button>
    <div class="collapse navbar-collapse" id="navbarNavAltMarkup">
      <div class="navbar-nav">
        {#each Const.PAGES as page}
          <Button
            size="sm"
            color="secondary"
            outline
            class="m-1 {$storedCurrentPage === page.name && 'active'}"
            title={page.title}
            style="font-size: {$storedUserConfig.font_size};"
            on:click={() => onSwitchPage(page)}
          >
            <i class={page.iconClass} />
            {page.title}
          </Button>
        {/each}
        {#if $storedCurrentPage == Page.Main}
          {#each Const.FUNCS as func}
            <Button
              size="sm"
              color="success"
              outline
              class="m-1"
              title={func.title}
              disabled={func.name === Func.Screenshot &&
                ($storedBattle === undefined || isScreenshotting)}
              style="font-size: {$storedUserConfig.font_size};"
              on:click={() => onClickFunc(func)}
            >
              {#if func.name === Func.Screenshot && isScreenshotting}
                <Spinner size="sm" type="border" /> 読み込み中
              {:else}
                <i class={func.iconClass} />
                {func.title}
              {/if}
            </Button>
          {/each}

          <select
            class="form-select form-select-sm m-1"
            style="font-size: {$storedUserConfig.font_size};"
            bind:value={selectedStatsPattern}
            on:change={onStatsPatternChanged}
          >
            {#await StatsPatterns() then statsPatterns}
              {#each statsPatterns as sp}
                {@const label = Const.STATS_PATTERN[sp]}
                <option
                  selected={sp === $storedUserConfig.stats_pattern}
                  value={sp}>{label}</option
                >
              {/each}
            {/await}
          </select>
        {/if}
      </div>
    </div>
  </div>
</nav>
