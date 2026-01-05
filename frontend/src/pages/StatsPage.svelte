<script lang="ts">
  import BattleMetadata from "@components/BattleMetadata.svelte";
  import MainStatsTable from "@components/MainStatsTable.svelte";
  import TeamAverageStatsDisplay from "@components/TeamAverageStatsDisplay.svelte";
  import MessagingTonako from "@components/MessagingTonako.svelte";
  import { storedTonako, storedBattle, storedUserConfig } from "@libs/stores";
</script>

<div>
  {#if $storedTonako}
    <MessagingTonako
      tonako={$storedTonako.tonako}
      message={$storedTonako.message}
      showLoading={$storedTonako.isLoading}
    />
  {:else if $storedBattle}
    <div class="pt-2 flex flex-col items-center">
      <BattleMetadata metadata={$storedBattle.metadata} />
    </div>

    <div class="flex justify-center">
      <MainStatsTable teams={$storedBattle.teams} />
    </div>

    {#if $storedBattle.teams.length === 2}
      <div class="flex justify-center px-4">
        <TeamAverageStatsDisplay
          friendTeam={$storedBattle.teams[0]}
          enemyTeam={$storedBattle.teams[1]}
        />
      </div>
    {/if}
  {/if}
</div>
