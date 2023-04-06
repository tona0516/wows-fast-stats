<script lang="ts">
  import SvelteTable from "svelte-table";
  import { BarLoader } from 'svelte-loading-spinners';
  import { Debug, GetTempArenaInfoHash, Load } from "../wailsjs/go/main/App.js";

  let latestHash = ""
  enum State {
    STANDBY,
    FETCHING,
};
  let state = State.STANDBY
  let friendRows = [];
  let enemyRows = [];
  let columns = [
    {
        title: "Clan",
        value: v => v.player_player_info.clan
    },
    {
        title: "Player name",
        value: v => v.player_player_info.name
    },
    {
        title: "CP",
        value: v => v.player_ship_stats.combat_power
    },
    {
        title: "PR",
        value: v => v.player_ship_stats.personal_rating
    },
    {
        title: "Ship name",
        value: v => v.player_ship_info.name
    },
    {
        title: "Tier",
        value: v => v.player_ship_info.tier
    },
    {
        title: "Ship dmg",
        value: v => v.player_ship_stats.avg_damage
    },
    {
        title: "Ship win",
        value: v => v.player_ship_stats.win_rate
    },
    {
        title: "Ship exp",
        value: v => v.player_ship_stats.avg_exp
    },
    {
        title: "Ship dmg",
        value: v => v.player_ship_stats.avg_damage
    },
    {
        title: "Ship battles",
        value: v => v.player_ship_stats.battles
    },
    {
        title: "Player dmg",
        value: v => v.player_player_stats.avg_damage
    },
    {
        title: "Player win",
        value: v => v.player_player_stats.win_rate
    },
    {
        title: "Player exp",
        value: v => v.player_player_stats.avg_exp
    },
    {
        title: "Player dmg",
        value: v => v.player_player_stats.avg_damage
    },
    {
        title: "Player battles",
        value: v => v.player_player_stats.battles
    },
    {
        title: "Tier avg",
        value: v => v.player_player_stats.avg_tier
    },
  ];


const looper = async () => {
    if (state === State.FETCHING) {
        return
    }

    const hash = await GetTempArenaInfoHash()
    if (hash === "") {
        return
    }
    if (hash !== latestHash) {
        state = State.FETCHING
        const stats = await Load()
        friendRows = stats["friends"]
        enemyRows = stats["enemies"]
        latestHash = hash
        state = State.STANDBY
    }
}

setInterval(looper, 1000);

</script>

<main>
  {#if state === State.FETCHING}
    <BarLoader color="#FF3E00" />
  {/if}

  {#if latestHash !== ""}
    <SvelteTable columns="{columns}" rows="{friendRows}"></SvelteTable>
    <SvelteTable columns="{columns}" rows="{enemyRows}"></SvelteTable>
  {/if}
</main>

<style>
</style>
