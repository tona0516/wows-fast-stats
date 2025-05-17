<script lang="ts">
import UkIcon from "src/component/common/uikit/UkIcon.svelte";
import UkSpinner from "src/component/common/uikit/UkSpinner.svelte";
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
  class="uk-select uk-form-width-medium uk-form-small"
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
  class="uk-button uk-button-primary uk-button-small"
  disabled={disableScreenshot}
  on:click={() => dispatch("ManualScreenshot")}
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
