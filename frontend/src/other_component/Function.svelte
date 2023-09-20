<script lang="ts">
  import clone from "clone";
  import { DispName } from "src/lib/DispName";
  import type { Screenshot } from "src/lib/Screenshot";
  import { ScreenshotType } from "src/lib/types";
  import { storedBattle, storedUserConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import { ApplyUserConfig } from "wailsjs/go/main/App";
  import { WindowReloadApp } from "wailsjs/runtime/runtime";

  export let screenshot: Screenshot;

  $: isScreenshotting = false;
  $: disableScreenshot = isScreenshotting || $storedBattle?.meta === undefined;
  let selectedStatsPattern: string = $storedUserConfig.stats_pattern;

  const dispatch = createEventDispatcher();

  const takeScreenshot = async () => {
    if (!$storedBattle?.meta) {
      return;
    }

    try {
      isScreenshotting = true;
      if (await screenshot.take(ScreenshotType.manual, $storedBattle.meta)) {
        dispatch("ScreenshotSaved");
      }
    } catch (error) {
      dispatch("Failure", { message: error });
    } finally {
      isScreenshotting = false;
    }
  };

  const onStatsPatternChanged = async () => {
    try {
      let config = clone($storedUserConfig);
      config.stats_pattern = selectedStatsPattern;
      await ApplyUserConfig(config);
      storedUserConfig.set(config);
    } catch (error) {
      selectedStatsPattern = $storedUserConfig.stats_pattern;
      dispatch("Failure", { message: error });
      return;
    }
  };
</script>

<div class="uk-flex uk-flex-center">
  <select
    class="uk-select uk-form-width-medium uk-form-small
    "
    bind:value={selectedStatsPattern}
    on:change={onStatsPatternChanged}
  >
    {#each DispName.STATS_PATTERNS.toArray() as sp}
      <option
        selected={sp.key == $storedUserConfig.stats_pattern}
        value={sp.key}>{sp.value}</option
      >
    {/each}
  </select>
  <button
    class="uk-button uk-button-primary uk-button-small"
    disabled={disableScreenshot}
    on:click={() => takeScreenshot()}
  >
    {#if isScreenshotting}
      <div uk-spinner />
    {:else}
      <span uk-icon="icon: camera" />
    {/if}
  </button>
  <button
    class="uk-button uk-button-primary uk-button-small"
    on:click={() => WindowReloadApp()}
  >
    <span uk-icon="icon: refresh" />
  </button>
</div>
