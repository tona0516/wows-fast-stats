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

  import nationJapan from "../assets/images/nation-japan.png";
  import nationUsa from "../assets/images/nation-usa.png";
  import nationUssr from "../assets/images/nation-ussr.png";
  import nationGermany from "../assets/images/nation-germany.png";
  import nationUk from "../assets/images/nation-uk.png";
  import nationFrance from "../assets/images/nation-france.png";
  import nationItaly from "../assets/images/nation-italy.png";
  import nationPanAsia from "../assets/images/nation-pan-asia.png";
  import nationEurope from "../assets/images/nation-europe.png";
  import nationNetherlands from "../assets/images/nation-netherlands.png";
  import nationCommonWealth from "../assets/images/nation-commonwealth.png";
  import nationPanAmerica from "../assets/images/nation-pan-america.png";
  import nationSpain from "../assets/images/nation-spain.png";
  import nationNone from "../assets/images/nation-none.png";

  import type { domain } from "../../wailsjs/go/models";
  import { shipURL, tierString } from "../util";

  export let player: domain.Player;
  export let userConfig: domain.UserConfig;

  function nationIcon(shipInfo: domain.ShipInfo): string {
    switch (shipInfo.nation) {
      case "japan":
        return nationJapan;
      case "usa":
        return nationUsa;
      case "ussr":
        return nationUssr;
      case "germany":
        return nationGermany;
      case "uk":
        return nationUk;
      case "france":
        return nationFrance;
      case "italy":
        return nationItaly;
      case "pan_asia":
        return nationPanAsia;
      case "europe":
        return nationEurope;
      case "netherlands":
        return nationNetherlands;
      case "commonwealth":
        return nationCommonWealth;
      case "pan_america":
        return nationPanAmerica;
      case "spain":
        return nationSpain;
      default:
        return nationNone;
    }
  }

  function shipTypeIcon(shipInfo: domain.ShipInfo): string {
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

<td>
  <img alt="" src={nationIcon(player.ship_info)} class="nation-icon" />
</td>

<td style="background-color: {color}">
  <img alt="" src={shipTypeIcon(player.ship_info)} class="ship-icon" />
</td>

<td class="td-string omit">
  <!-- svelte-ignore a11y-invalid-attribute -->
  <a class="td-link" href="#" on:click={() => BrowserOpenURL(shipURL(player))}
    >{tierString(player.ship_info.tier)} {player.ship_info.name}
  </a>
</td>
