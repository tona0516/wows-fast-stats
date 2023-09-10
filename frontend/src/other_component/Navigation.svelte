<script lang="ts">
  import clone from "clone";
  import { DispName } from "src/lib/DispName";
  import type { Screenshot } from "src/lib/Screenshot";
  import { Func, Page, ScreenshotType } from "src/lib/types";
  import {
    storedCurrentPage,
    storedUserConfig,
    storedBattle,
  } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import { Button, Spinner } from "sveltestrap";
  import { ApplyUserConfig } from "wailsjs/go/main/App";
  import { WindowReloadApp } from "wailsjs/runtime/runtime";

  export let screenshot: Screenshot;

  let isScreenshotting = false;
  let meta = $storedBattle?.meta;
  storedBattle.subscribe((it) => {
    meta = it?.meta;
  });
  let selectedStatsPattern: string = $storedUserConfig.stats_pattern;

  $: disableScreenshotButton = meta === undefined || isScreenshotting;

  const dispatch = createEventDispatcher();
  const navID = "navbarNavAltMarkup";

  const onClickFunc = (func: Func) => {
    switch (func) {
      case Func.RELOAD:
        WindowReloadApp();
        break;
      case Func.SCREENSHOT:
        takeScreenshot();
        break;
    }
  };

  const takeScreenshot = async () => {
    if (!meta) {
      return;
    }

    try {
      isScreenshotting = true;
      if (await screenshot.take(ScreenshotType.manual, meta)) {
        dispatch("Success", { message: "スクリーンショットを保存しました。" });
      }
    } catch (error) {
      dispatch("Failure", { message: error });
    } finally {
      isScreenshotting = false;
    }
  };

  const onStatsPatternChanged = async () => {
    // Note: for the following sveltestrap bug
    // https://github.com/bestguy/sveltestrap/issues/461
    await new Promise((resolve) => setTimeout(resolve, 100));

    try {
      let config = clone($storedUserConfig);
      config.stats_pattern = selectedStatsPattern;
      await ApplyUserConfig(config);
      storedUserConfig.set(config);
    } catch (error) {
      dispatch("Failure", { message: error });
      return;
    }
  };
</script>

<!-- Note: doesn't show buttons with sveltestrap -->
<nav class="navbar navbar-expand-sm sticky-top navbar-light bg-light p-1">
  <div class="container-fluid">
    <button
      class="navbar-toggler"
      type="button"
      data-bs-toggle="collapse"
      data-bs-target="#{navID}"
      aria-controls={navID}
      aria-expanded="false"
    >
      <span class="navbar-toggler-icon" />
    </button>
    <div class="collapse navbar-collapse" id={navID}>
      <div class="navbar-nav">
        {#each DispName.PAGES as page}
          <Button
            size="sm"
            color="secondary"
            outline
            class="m-1 {$storedCurrentPage === page.first ? 'active' : ''}"
            style="font-size: {$storedUserConfig.font_size};"
            on:click={() => storedCurrentPage.set(page.first)}
          >
            <i class={page.third} />
            {page.second}
          </Button>
        {/each}
        {#if $storedCurrentPage == Page.MAIN}
          {#each DispName.FUNCS as func}
            <Button
              size="sm"
              color="success"
              outline
              class="m-1"
              disabled={func.first === Func.SCREENSHOT &&
                disableScreenshotButton}
              style="font-size: {$storedUserConfig.font_size};"
              on:click={() => onClickFunc(func.first)}
            >
              {#if func.first === Func.SCREENSHOT && isScreenshotting}
                <Spinner size="sm" /> 読み込み中
              {:else}
                <i class={func.third} />
                {func.second}
              {/if}
            </Button>
          {/each}

          <!-- Note: sveltestrap "input" binds empty value when page changed -->
          <select
            class="form-select form-select-sm m-1"
            style="font-size: {$storedUserConfig.font_size};"
            bind:value={selectedStatsPattern}
            on:change={onStatsPatternChanged}
          >
            {#each DispName.STATS_PATTERNS as pair}
              <option
                selected={pair.first == $storedUserConfig.stats_pattern}
                value={pair.first}>{pair.second}</option
              >
            {/each}
          </select>
        {/if}
      </div>
    </div>
  </div>
</nav>
