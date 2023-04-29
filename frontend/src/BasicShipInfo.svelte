<script lang="ts">
import type { vo } from "wailsjs/go/models";
import iconCV from "./assets/images/icon-cv.png";
import iconBB from "./assets/images/icon-bb.png";
import iconCL from "./assets/images/icon-cl.png";
import iconDD from "./assets/images/icon-dd.png";
import iconSS from "./assets/images/icon-ss.png";
import iconNone from "./assets/images/icon-none.png";
import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
export let config: vo.UserConfig;
export let player: vo.Player;
export let displayPattern: DisplayPattern;

function tierString(value: number): string {
  if (value === 11) return "★";

  const decimal = [10, 9, 5, 4, 1];
  const romanNumeral = ["X", "IX", "V", "IV", "I"];

  let romanized = "";

  for (var i = 0; i < decimal.length; i++) {
    while (decimal[i] <= value) {
      romanized += romanNumeral[i];
      value -= decimal[i];
    }
  }

  return romanized;
}

function shipIcon(shipType: string): string {
  switch (shipType) {
    case "AirCarrier":
      return iconCV;
    case "Battleship":
      return iconBB;
    case "Cruiser":
      return iconCL;
    case "Destroyer":
      return iconDD;
    case "Submarine":
      return iconSS;
    default:
      return iconNone;
  }
}
</script>

<td class="td-string">
  <!-- svelte-ignore a11y-invalid-attribute -->
  <a
    class="td-link"
    href="#"
    on:click="{() => BrowserOpenURL(player.ship_info.stats_url)}"
  >
    <div class="horizontal">
      <img
        alt=""
        src="{shipIcon(player.ship_info.type)}"
        class="ship-icon-scale"
      />
      <div class="omit">
        {tierString(player.ship_info.tier)}
        {player.ship_info.name}
      </div>
    </div>
  </a>
</td>
