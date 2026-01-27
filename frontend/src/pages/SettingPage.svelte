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

<div class="container mx-auto max-w-3xl py-3 flex flex-col gap-4">
  <GameClientPathSection />
  <GlobalDisplaySection />
  <ColumnDisplaySection />
</div>
