<script lang="ts">
  import { DispName } from "src/lib/DispName";
  import type { StatsTableOptions } from "src/lib/StatsTableOptions";

  export let options: StatsTableOptions;
</script>

<select
  class="uk-select uk-form-width-medium uk-form-small"
  id="stats_pattern"
  bind:value={options.statsPattern}
>
  {#each DispName.STATS_PATTERNS.toArray() as sp}
    <option selected={sp.key == options.statsPattern} value={sp.key}
      >{sp.value}</option
    >
  {/each}
</select>
