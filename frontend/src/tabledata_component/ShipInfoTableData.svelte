<script lang="ts">
  import { BrowserOpenURL } from "../../wailsjs/runtime/runtime";
  import iconCV from "../assets/images/icon-cv.png";
  import iconBB from "../assets/images/icon-bb.png";
  import iconCL from "../assets/images/icon-cl.png";
  import iconDD from "../assets/images/icon-dd.png";
  import iconSS from "../assets/images/icon-ss.png";
  import iconCVPremium from "../assets/images/icon-cv-premium.png";
  import iconBBPremium from "../assets/images/icon-bb-premium.png";
  import iconCLPremium from "../assets/images/icon-cl-premium.png";
  import iconDDPremium from "../assets/images/icon-dd-premium.png";
  import iconSSPremium from "../assets/images/icon-ss-premium.png";
  import iconNone from "../assets/images/icon-none.png";
  import type { domain } from "../../wailsjs/go/models";
  import { shipURL, tierString } from "../util";

  export let player: domain.Player;
  export let userConfig: domain.UserConfig;

  function shipIcon(shipInfo: domain.ShipInfo): string {
    switch (shipInfo.type) {
      case "cv":
        return shipInfo.is_premium ? iconCVPremium : iconCV;
      case "bb":
        return shipInfo.is_premium ? iconBBPremium : iconBB;
      case "cl":
        return shipInfo.is_premium ? iconCLPremium : iconCL;
      case "dd":
        return shipInfo.is_premium ? iconDDPremium : iconDD;
      case "ss":
        return shipInfo.is_premium ? iconSSPremium : iconSS;
      default:
        return iconNone;
    }
  }

  $: color =
    userConfig.custom_color.ship_type.own[player.ship_info.type] ?? "#00000000";
</script>

<td style="background-color: {color}">
  <img alt="" src={shipIcon(player.ship_info)} class="ship-icon" />
</td>

<td class="td-string omit">
  <!-- svelte-ignore a11y-invalid-attribute -->
  <a class="td-link" href="#" on:click={() => BrowserOpenURL(shipURL(player))}
    >{tierString(player.ship_info.tier)} {player.ship_info.name}
  </a>
</td>
