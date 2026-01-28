<script lang="ts">
  import GameClientPathSection from "@components/setting/GameClientPathSection.svelte";
  import { onMount } from "svelte";
  import { themeChange } from "theme-change";
  import GlobalDisplaySection from "@components/setting/GlobalDisplaySection.svelte";
  import ColumnDisplaySection from "@components/setting/ColumnDisplaySection.svelte";

  onMount(async () => {
    themeChange(false);
  });
</script>

<div class="min-h-full w-full bg-base-200/40 p-6">
  <div class="mx-auto flex w-full max-w-3xl flex-col gap-6">
    <GameClientPathSection />
    <GlobalDisplaySection />
    <ColumnDisplaySection />
  </div>
</div>
