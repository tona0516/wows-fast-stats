<script lang="ts">
  import { STATS_COLUMN_INFO } from "@libs/constants";
  import type { StatsKey } from "@libs/types";

  const visibleLabel = "表示中";
  const hiddenLabel = "非表示";

  export let visibleOrder: StatsKey[] = [];
  export let hiddenOrder: StatsKey[] = [];
  export let visibleListElement: HTMLUListElement | null = null;
  export let hiddenListElement: HTMLUListElement | null = null;
</script>

<div class="grid grid-cols-1 gap-5 md:grid-cols-2">
  <div>
    <div class="mb-2 text-xs font-semibold text-base-content/70">
      {visibleLabel}
    </div>
    <ul
      bind:this={visibleListElement}
      class="space-y-2 min-h-60 rounded-lg border border-base-300 bg-base-100/60 p-2"
    >
      {#each visibleOrder as key (key)}
        <li
          class="flex items-center gap-2 rounded-lg border border-base-300 bg-base-100 px-3 py-2 text-sm"
          data-key={key}
        >
          <span>{STATS_COLUMN_INFO[key].minName}</span>
        </li>
      {/each}
    </ul>
  </div>
  <div>
    <div class="mb-2 text-xs font-semibold text-base-content/70">
      {hiddenLabel}
    </div>
    <ul
      bind:this={hiddenListElement}
      class="space-y-2 min-h-60 rounded-lg border border-base-300 bg-base-100/60 p-2"
    >
      {#each hiddenOrder as key (key)}
        <li
          class="flex items-center gap-2 rounded-lg border border-base-300 bg-base-100 px-3 py-2 text-sm"
          data-key={key}
        >
          <span>{STATS_COLUMN_INFO[key].minName}</span>
        </li>
      {/each}
    </ul>
  </div>
</div>
