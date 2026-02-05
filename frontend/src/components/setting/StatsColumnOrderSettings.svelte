<script lang="ts">
  import Sortable from "sortablejs";
  import { DEFAULT_DISPLAY_PREF } from "@libs/DisplayPref";
  import { storedDisplayPref } from "@libs/stores";
  import type { StatsKey } from "@libs/types";
  import { SaveDisplayPref } from "@wails/go/main/App";
  import StatsColumnOrderSection from "@components/setting/StatsColumnOrderSection.svelte";
  import { get } from "svelte/store";
  import { onDestroy, onMount } from "svelte";

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
</script>

<div class="mt-6 grid grid-cols-1 gap-6">
  <StatsColumnOrderSection
    title="艦成績の表示項目"
    visibleOrder={shipVisibleOrder}
    hiddenOrder={shipHiddenOrder}
    bind:visibleListElement={shipListElement}
    bind:hiddenListElement={shipHiddenListElement}
  />

  <StatsColumnOrderSection
    title="総合成績の表示項目"
    visibleOrder={overallVisibleOrder}
    hiddenOrder={overallHiddenOrder}
    bind:visibleListElement={overallListElement}
    bind:hiddenListElement={overallHiddenListElement}
  />
</div>
