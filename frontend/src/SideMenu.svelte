<script lang="ts">
  import { WindowReloadApp } from "wailsjs/runtime/runtime";
  import type { Page } from "./lib/types";
  export let page: Page;

  const MENU_ITEMS: { page: Page; icon: string; tooltip: string }[] = [
    { page: "stats", icon: "bi-table", tooltip: "戦闘統計" },
    { page: "config", icon: "bi-gear", tooltip: "設定" },
    { page: "info", icon: "bi-info-circle", tooltip: "アプリ情報" },
  ];
</script>

<ul class="menu rounded-box sticky top-0">
  <li>
    {#each MENU_ITEMS as item}
      <button
        class="btn btn-ghost {page === item.page &&
          'bg-current/[15%]'} tooltip tooltip-right"
        data-tip={item.tooltip}
        on:click={() => (page = item.page)}
      >
        <i class="bi {item.icon}"></i>
      </button>
    {/each}
  </li>

  <li class="mt-8">
    <button
      class="btn btn-primary tooltip tooltip-right"
      data-tip="リロード"
      on:click={() => WindowReloadApp()}
    >
      <i class="bi bi-arrow-clockwise"></i>
    </button>
  </li>
</ul>
