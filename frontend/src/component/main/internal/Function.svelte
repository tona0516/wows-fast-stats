<script lang="ts">
  import clone from "clone";
  import UkIcon from "src/component/common/uikit/UkIcon.svelte";
  import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
  import { DispName } from "src/lib/DispName";
  import { Notifier } from "src/lib/Notifier";
  import type { Screenshot } from "src/lib/Screenshot";
  import { storedBattle, storedConfig } from "src/stores";
  import { createEventDispatcher } from "svelte";
  import { ApplyUserConfig } from "wailsjs/go/main/App";
  import { WindowReloadApp } from "wailsjs/runtime/runtime";

  export let screenshot: Screenshot;

  $: isScreenshotting = false;
  $: disableScreenshot = isScreenshotting || $storedBattle?.meta === undefined;
  let selectedStatsPattern: string = $storedConfig.stats_pattern;

  const dispatch = createEventDispatcher();

  const takeScreenshot = async () => {
    if (!$storedBattle?.meta) {
      return;
    }

    try {
      isScreenshotting = true;
      if (await screenshot.manual($storedBattle.meta)) {
        Notifier.success("スクリーンショットを保存しました");
      }
    } catch (error) {
      Notifier.failure(error);
    } finally {
      isScreenshotting = false;
    }
  };

  const onStatsPatternChanged = async () => {
    try {
      let config = clone($storedConfig);
      config.stats_pattern = selectedStatsPattern;
      await ApplyUserConfig(config);
      storedConfig.set(config);
    } catch (error) {
      selectedStatsPattern = $storedConfig.stats_pattern;
      Notifier.failure(error);
      return;
    }
  };
</script>

<select
  class="uk-select uk-form-width-medium uk-form-small"
  bind:value={selectedStatsPattern}
  on:change={onStatsPatternChanged}
>
  {#each DispName.STATS_PATTERNS.toArray() as sp}
    <option selected={sp.key == $storedConfig.stats_pattern} value={sp.key}
      >{sp.value}</option
    >
  {/each}
</select>
<button
  class="uk-button uk-button-primary uk-button-small"
  disabled={disableScreenshot}
  on:click={() => takeScreenshot()}
>
  {#if isScreenshotting}
    <UkSpinner />
  {:else}
    <UkIcon name="camera" />
  {/if}
</button>
<button
  class="uk-button uk-button-primary uk-button-small"
  on:click={() => WindowReloadApp()}
>
  <UkIcon name="refresh" />
</button>
