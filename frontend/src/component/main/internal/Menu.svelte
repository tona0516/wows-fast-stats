<script lang="ts">
import { DispName } from "src/lib/DispName";
import { Notifier } from "src/lib/Notifier";
import { storedConfig } from "src/stores";
import { UpdateUserConfig } from "wailsjs/go/main/App";
import { WindowReloadApp } from "wailsjs/runtime/runtime";

$: inputConfig = $storedConfig;

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

<button class="btn btn-primary" on:click={() => WindowReloadApp()}>
  <i class="bi bi-arrow-clockwise"></i>
</button>
