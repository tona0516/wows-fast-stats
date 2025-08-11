<script lang="ts">
  import {
    storedBattle,
    storedInstallPathError,
    storedTonako,
  } from "src/stores";

  import BattleMetaInfo from "./internal/BattleMetaInfo.svelte";
  import MainStatsTable from "./internal/MainStatsTable.svelte";
  import MessagingTonako from "./internal/MessagingTonako.svelte";
  import { Tonako } from "./internal/Tonako";
</script>

<div>
  {#if $storedInstallPathError}
    <MessagingTonako
      tonako={Tonako.Pointing}
      message="設定画面から初期設定をおこなってください"
    />
  {:else if $storedTonako}
    <MessagingTonako
      tonako={$storedTonako.tonako}
      message={$storedTonako.message}
      showLoading={$storedTonako.isLoading}
    />
  {:else if $storedBattle}
    <div class="pt-2 flex flex-col items-center">
      <BattleMetaInfo meta={$storedBattle.meta} />
    </div>

    <div class="flex">
      <MainStatsTable teams={$storedBattle.teams} />
    </div>
  {/if}
</div>
