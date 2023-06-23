<script lang="ts">
import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
import iconCV from "../assets/images/icon-cv.png";
import iconBB from "../assets/images/icon-bb.png";
import iconCL from "../assets/images/icon-cl.png";
import iconDD from "../assets/images/icon-dd.png";
import iconSS from "../assets/images/icon-ss.png";
import iconNone from "../assets/images/icon-none.png";
import type { vo } from "../../wailsjs/go/models";
import { shipURL, tierString } from "../util";

export let player: vo.Player;
export let userConfig: vo.UserConfig;

function shipIcon(shipType: string): string {
  switch (shipType) {
    case "cv":
      return iconCV;
    case "bb":
      return iconBB;
    case "cl":
      return iconCL;
    case "dd":
      return iconDD;
    case "ss":
      return iconSS;
    default:
      return iconNone;
  }
}

$: color =
  userConfig.custom_color.ship_type.own[player.ship_info.type] ?? "#00000000";
</script>

<td style="width: 1em; background-color: {color}">
  <img alt="" src="{shipIcon(player.ship_info.type)}" class="ship-icon" />
</td>

<td class="td-string omit">
  <!-- svelte-ignore a11y-invalid-attribute -->
  <a class="td-link" href="#" on:click="{() => BrowserOpenURL(shipURL(player))}"
    >{tierString(player.ship_info.tier)} {player.ship_info.name}
  </a>
</td>
