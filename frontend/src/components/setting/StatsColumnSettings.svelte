<script lang="ts">
  import Sortable from "sortablejs";
  import { STATS_COLUMN_INFO } from "@libs/constants";
  import { DEFAULT_DISPLAY_PREF } from "@libs/DisplayPref";
  import { storedDisplayPref } from "@libs/stores";
  import type { StatsKey } from "@libs/types";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import { get } from "svelte/store";
  import { onDestroy, onMount } from "svelte";

  const STATS_KEYS = Object.keys(STATS_COLUMN_INFO) as readonly StatsKey[];
  type StatsColumnCategory = "ship" | "overall";

  let shipListElement: HTMLUListElement | null = null;
  let shipHiddenListElement: HTMLUListElement | null = null;
  let overallListElement: HTMLUListElement | null = null;
  let overallHiddenListElement: HTMLUListElement | null = null;
  let sortables: Sortable[] = [];
  let shipVisibleOrder: StatsKey[] = [];
  let shipHiddenOrder: StatsKey[] = [];
  let overallVisibleOrder: StatsKey[] = [];
  let overallHiddenOrder: StatsKey[] = [];

  const normalizeOrder = (
    order: StatsKey[] | undefined,
    fallback: StatsKey[],
  ): StatsKey[] => {
    if (!order || order.length === 0) return fallback;
    const fallbackSet = new Set(fallback);
    const normalized = order.filter((key) => fallbackSet.has(key));
    const missing = fallback.filter((key) => !normalized.includes(key));
    return [...normalized, ...missing];
  };

  const getOrder = (category: StatsColumnCategory): StatsKey[] => {
    const pref = get(storedDisplayPref);
    if (!pref) return DEFAULT_DISPLAY_PREF.columnOrder[category];
    return normalizeOrder(
      pref.columnOrder?.[category],
      DEFAULT_DISPLAY_PREF.columnOrder[category],
    );
  };

  const persistOrder = async (
    category: StatsColumnCategory,
    visibleOrder: StatsKey[],
    hiddenOrder: StatsKey[],
  ) => {
    const pref = get(storedDisplayPref);
    if (!pref) return;
    if (!pref.columnOrder) {
      pref.columnOrder = {
        ship: [...DEFAULT_DISPLAY_PREF.columnOrder.ship],
        overall: [...DEFAULT_DISPLAY_PREF.columnOrder.overall],
      };
    }
    pref.columnOrder[category] = [...visibleOrder, ...hiddenOrder];

    const updateVisibility = (key: StatsKey, isVisible: boolean) => {
      const setting = pref[key];
      if (category === "ship" && "showShip" in setting) {
        setting.showShip = isVisible;
      }
      if (category === "overall" && "showOverall" in setting) {
        setting.showOverall = isVisible;
      }
    };

    visibleOrder.forEach((key) => updateVisibility(key, true));
    hiddenOrder.forEach((key) => updateVisibility(key, false));
    storedDisplayPref.set(pref);
    await SaveDisplayPref(JSON.stringify(pref, null, 2));
  };

  const createSortable = (
    element: HTMLUListElement | null,
    visibleElement: HTMLUListElement | null,
    hiddenElement: HTMLUListElement | null,
    category: StatsColumnCategory,
  ) => {
    if (!element) return null;
    return new Sortable(element, {
      animation: 150,
      group: category,
      onEnd: async () => {
        const visibleOrder = Array.from(
          visibleElement?.querySelectorAll<HTMLElement>("[data-key]") ?? [],
        )
          .map((item) => item.dataset.key)
          .filter((key): key is StatsKey => Boolean(key));
        const hiddenOrder = Array.from(
          hiddenElement?.querySelectorAll<HTMLElement>("[data-key]") ?? [],
        )
          .map((item) => item.dataset.key)
          .filter((key): key is StatsKey => Boolean(key));
        await persistOrder(category, visibleOrder, hiddenOrder);
      },
    });
  };

  const isSortable = (sortable: Sortable | null): sortable is Sortable =>
    Boolean(sortable);

  const splitOrder = (category: StatsColumnCategory) => {
    const pref = get(storedDisplayPref);
    const order = getOrder(category);
    if (!pref) {
      return { visible: order, hidden: [] };
    }

    const isVisible = (key: StatsKey) => {
      const setting = pref[key];
      if (category === "ship" && "showShip" in setting) {
        return setting.showShip === true;
      }
      if (category === "overall" && "showOverall" in setting) {
        return setting.showOverall === true;
      }
      return false;
    };

    return {
      visible: order.filter((key) => isVisible(key)),
      hidden: order.filter((key) => !isVisible(key)),
    };
  };

  onMount(() => {
    sortables = [
      createSortable(
        shipListElement,
        shipListElement,
        shipHiddenListElement,
        "ship",
      ),
      createSortable(
        shipHiddenListElement,
        shipListElement,
        shipHiddenListElement,
        "ship",
      ),
      createSortable(
        overallListElement,
        overallListElement,
        overallHiddenListElement,
        "overall",
      ),
      createSortable(
        overallHiddenListElement,
        overallListElement,
        overallHiddenListElement,
        "overall",
      ),
    ].filter(isSortable);
  });

  onDestroy(() => {
    sortables.forEach((sortable) => sortable.destroy());
    sortables = [];
  });

  $: ({ visible: shipVisibleOrder, hidden: shipHiddenOrder } =
    splitOrder("ship"));
  $: ({ visible: overallVisibleOrder, hidden: overallHiddenOrder } =
    splitOrder("overall"));

  const onChangePref = async () => {
    const pref = get(storedDisplayPref);
    if (!pref) return;
    await SaveDisplayPref(JSON.stringify(pref, null, 2));
  };
</script>

<div
  class="overflow-x-auto rounded-xl border border-base-300 bg-base-200 shadow-sm"
>
  <table class="table w-full text-nowrap">
    <thead>
      <tr class="bg-base-300 text-base-content font-semibold">
        {#each ["カラム名", "小数点以下の桁数", "K(キロ)表示"] as columns}
          <th
            class="px-4 py-3 text-center text-xs font-semibold uppercase tracking-wide"
          >
            {columns}
          </th>
        {/each}
      </tr>
    </thead>
    <tbody>
      {#each STATS_KEYS as statsKey}
        {@const info = STATS_COLUMN_INFO[statsKey]}
        <tr class="hover:bg-base-100">
          <td class="px-4 py-3 text-sm font-medium">{info.fullName}</td>

          {#if "digit" in $storedDisplayPref[statsKey] && typeof $storedDisplayPref[statsKey].digit === "number"}
            <td class="px-4 py-3 text-center">
              <select
                class="select select-sm select-bordered"
                bind:value={$storedDisplayPref[statsKey].digit}
                on:change={onChangePref}
              >
                {#each [0, 1, 2] as digit}
                  <option
                    selected={digit === $storedDisplayPref[statsKey].digit}
                    value={digit}>{digit}</option
                  >
                {/each}
              </select>
            </td>
          {:else}
            <td></td>
          {/if}
          {#if "siPrefix" in $storedDisplayPref[statsKey] && typeof $storedDisplayPref[statsKey].siPrefix === "boolean"}
            <td class="px-4 py-3 text-center">
              <input
                class="toggle toggle-success"
                type="checkbox"
                bind:checked={$storedDisplayPref[statsKey].siPrefix}
                on:change={onChangePref}
              />
            </td>
          {:else}
            <td></td>
          {/if}
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<div class="mt-6 grid grid-cols-1 gap-6">
  <div class="rounded-xl border border-base-300 bg-base-200 p-5 shadow-sm">
    <div class="mb-3 text-sm font-semibold">艦成績の表示項目</div>
    <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
      <div>
        <div class="mb-2 text-xs font-semibold text-base-content/70">
          表示中
        </div>
        <ul
          bind:this={shipListElement}
          class="space-y-2 min-h-60 rounded-lg border border-base-300 bg-base-100/60 p-2"
        >
          {#each shipVisibleOrder as key (key)}
            <li
              class="flex items-center gap-2 rounded-lg border border-base-300 bg-base-100 px-3 py-2 text-sm"
              data-key={key}
            >
              <span>{STATS_COLUMN_INFO[key].fullName}</span>
            </li>
          {/each}
        </ul>
      </div>
      <div>
        <div class="mb-2 text-xs font-semibold text-base-content/70">
          非表示
        </div>
        <ul
          bind:this={shipHiddenListElement}
          class="space-y-2 min-h-60 rounded-lg border border-base-300 bg-base-100/60 p-2"
        >
          {#each shipHiddenOrder as key (key)}
            <li
              class="flex items-center gap-2 rounded-lg border border-base-300 bg-base-100 px-3 py-2 text-sm"
              data-key={key}
            >
              <span>{STATS_COLUMN_INFO[key].fullName}</span>
            </li>
          {/each}
        </ul>
      </div>
    </div>
  </div>

  <div class="rounded-xl border border-base-300 bg-base-200 p-5 shadow-sm">
    <div class="mb-3 text-sm font-semibold">総合成績の表示項目</div>
    <div class="grid grid-cols-1 gap-5 md:grid-cols-2">
      <div>
        <div class="mb-2 text-xs font-semibold text-base-content/70">
          表示中
        </div>
        <ul
          bind:this={overallListElement}
          class="space-y-2 min-h-60 rounded-lg border border-base-300 bg-base-100/60 p-2"
        >
          {#each overallVisibleOrder as key (key)}
            <li
              class="flex items-center gap-2 rounded-lg border border-base-300 bg-base-100 px-3 py-2 text-sm"
              data-key={key}
            >
              <span>{STATS_COLUMN_INFO[key].fullName}</span>
            </li>
          {/each}
        </ul>
      </div>
      <div>
        <div class="mb-2 text-xs font-semibold text-base-content/70">
          非表示
        </div>
        <ul
          bind:this={overallHiddenListElement}
          class="space-y-2 min-h-60 rounded-lg border border-base-300 bg-base-100/60 p-2"
        >
          {#each overallHiddenOrder as key (key)}
            <li
              class="flex items-center gap-2 rounded-lg border border-base-300 bg-base-100 px-3 py-2 text-sm"
              data-key={key}
            >
              <span>{STATS_COLUMN_INFO[key].fullName}</span>
            </li>
          {/each}
        </ul>
      </div>
    </div>
  </div>
</div>
