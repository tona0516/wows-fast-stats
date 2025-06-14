<script lang="ts">
import { DispName } from "src/lib/DispName";
import { Notifier } from "src/lib/Notifier";
import { storedBattle, storedConfig } from "src/stores";
import { createEventDispatcher } from "svelte";
import { UpdateUserConfig } from "wailsjs/go/main/App";
import { WindowReloadApp } from "wailsjs/runtime/runtime";

export let isScreenshotting: boolean;

$: inputConfig = $storedConfig;
$: disableScreenshot = isScreenshotting || $storedBattle?.meta === undefined;

const dispatch = createEventDispatcher();

const onStatsPatternChanged = async () => {
  try {
    await UpdateUserConfig(inputConfig);
  } catch (error) {
    inputConfig.stats_pattern = $storedConfig.stats_pattern;
    Notifier.failure(error);
    return;
  }
};
</script>

<select
  class="select"
  bind:value={inputConfig.stats_pattern}
  on:change={onStatsPatternChanged}
>
  {#each DispName.STATS_PATTERNS.toArray() as sp}
    <option selected={sp.key == $storedConfig.stats_pattern} value={sp.key}
      >{sp.value}</option
    >
  {/each}
</select>
<button
  class="btn btn-primary"
  disabled={disableScreenshot}
  on:click={() => dispatch("ManualScreenshot")}
>
  {#if isScreenshotting}
    <span class="loading loading-spinner"></span>
  {:else}
    <i class="bi bi-camera"></i>
  {/if}
</button>
<button class="btn btn-primary" on:click={() => WindowReloadApp()}>
  <i class="bi bi-arrow-clockwise"></i>
</button>
