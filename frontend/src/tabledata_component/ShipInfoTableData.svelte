<script lang="ts">
  import { Const } from "src/Const";
  import { shipTypeIcon, shipURL, tierString } from "src/util";
  import type { domain } from "wailsjs/go/models";
  import { BrowserOpenURL } from "wailsjs/runtime/runtime";

  export let player: domain.Player;
  export let userConfig: domain.UserConfig;

  $: color =
    userConfig.custom_color.ship_type.own[player.ship_info.type] ?? "#00000000";
</script>

<td class="td-icon">
  <img
    alt=""
    src={Const.NATION_ICON[player.ship_info.nation] ?? Const.NATION_ICON.none}
    class="nation-icon"
  />
</td>

<td class="td-icon" style="background-color: {color}">
  <img alt="" src={shipTypeIcon(player.ship_info)} class="ship-icon" />
</td>

<td class="td-string omit">
  <!-- svelte-ignore a11y-invalid-attribute -->
  <a class="td-link" href="#" on:click={() => BrowserOpenURL(shipURL(player))}
    >{tierString(player.ship_info.tier)} {player.ship_info.name}
  </a>
</td>
