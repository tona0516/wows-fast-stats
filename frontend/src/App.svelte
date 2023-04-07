<script lang="ts">
  import SvelteTable from "svelte-table";
  import { BarLoader } from 'svelte-loading-spinners';
  import { toasts, ToastContainer, FlatToast }  from "svelte-toasts";
  import { ApplyConfig, Debug, GetConfig, GetTempArenaInfoHash, Load, SelectDirectory } from "../wailsjs/go/main/App.js";

  enum Page {
    MAIN,
    CONFIG
  }
  let currentPage = Page.MAIN

  enum State {
    STANDBY,
    FETCHING,
    ERROR,
  };
  let state = State.STANDBY
  let latestHash = ""

  let installPath = ""
  let appid = ""

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


  async function looper() {
    if (state === State.ERROR) {
        return
    }

    if (state === State.FETCHING) {
        return
    }

    let hash: string
    try {
        hash = await GetTempArenaInfoHash()
    } catch (error) {
        state = State.ERROR
        showErrorToast(error)
        return
    }

    if (hash === latestHash) {
        return
    }

    state = State.FETCHING
    try {
        const stats = await Load()
        friendRows = stats["friends"]
        enemyRows = stats["enemies"]
        latestHash = hash
        state = State.STANDBY
        showSuccessToast("updated!")
    } catch (error) {
        state = State.ERROR
        showErrorToast(error)
    }
}

setInterval(looper, 1000);

function clickMain() {
    currentPage = Page.MAIN
}

function clickConfig() {
    currentPage = Page.CONFIG
    GetConfig().then((config) => {
        installPath = config.install_path
        appid = config.appid
    }).catch((error) => {
        showErrorToast(error)
        installPath = ""
        appid = ""
    })
}

function clickApply() {
    ApplyConfig(installPath, appid).then(_ => {
        showSuccessToast("updated!")
    }).catch((error) => {
        showErrorToast(error)
    })
}

function selectDirectory() {
    SelectDirectory().then((result) => {
        installPath = result
    })
}

function showSuccessToast(message: string) {
    toasts.add({
      title: "Success",
      description: message,
      duration: 3000,
      placement: 'bottom-right',
      type: 'success',
      theme: 'dark',
      showProgress: true,
    });
}

function showErrorToast(message: string) {
    toasts.add({
      title: "Error",
      description: message,
      duration: 3000,
      placement: 'bottom-right',
      type: 'error',
      theme: 'dark',
      showProgress: true,
    });
}

</script>

<main>
  {#if currentPage === Page.CONFIG}
    <div>
      <input type=”text” bind:value={installPath} size="50" placeholder="World of Warshipsインストールフォルダ">
      <button on:click={selectDirectory}>フォルダ選択</button>
    </div>

    <div>
      <input type=”text” bind:value={appid} size="50" id="appid" placeholder="App ID">
    </div>

    <div>
      <button on:click={clickApply}>適用</button>
      <button on:click={clickMain}>戻る</button>
    </div>
  {/if}

  {#if currentPage === Page.MAIN}
    {#if state === State.FETCHING}
      <BarLoader color="#FF3E00" />
    {/if}

    {#if latestHash !== ""}
      <SvelteTable columns="{columns}" rows="{friendRows}"></SvelteTable>
      <SvelteTable columns="{columns}" rows="{enemyRows}"></SvelteTable>
    {/if}

    <button on:click={clickConfig}>設定</button>
  {/if}

  <ToastContainer let:data={data}>
    <FlatToast {data} />
  </ToastContainer>

</main>

<style>
</style>
