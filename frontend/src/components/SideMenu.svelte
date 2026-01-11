<script lang="ts">
  import type { Page } from "@libs/types";
  import { WindowReloadApp } from "@wails/runtime/runtime";

  export let page: Page;

  const MENU_ITEMS: { page: Page; icon: string; tooltip: string }[] = [
    { page: "stats", icon: "bi-table", tooltip: "戦闘統計" },
    { page: "pref", icon: "bi-gear", tooltip: "設定" },
    { page: "info", icon: "bi-info-circle", tooltip: "アプリ情報" },
  ];
</script>

<ul class="menu menu-sm rounded-box sticky top-0 h-screen">
  <li class="gap-1">
    {#each MENU_ITEMS as item}
      <button
        class="btn btn-sm btn-ghost tooltip tooltip-right {page === item.page &&
          'bg-current/15'} "
        data-tip={item.tooltip}
        on:click={() => (page = item.page)}
      >
        <i class="bi {item.icon}"></i>
      </button>
    {/each}
  </li>

  <li class="mt-8">
    <button
      class="btn btn-sm btn-primary tooltip tooltip-right"
      data-tip="リロード"
      on:click={WindowReloadApp}
    >
      <i class="bi bi-arrow-clockwise"></i>
    </button>
  </li>
</ul>
