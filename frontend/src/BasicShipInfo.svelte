<script lang="ts">
import type { vo } from "wailsjs/go/models";
import iconCV from "./assets/images/icon-cv.png";
import iconBB from "./assets/images/icon-bb.png";
import iconCL from "./assets/images/icon-cl.png";
import iconDD from "./assets/images/icon-dd.png";
import iconSS from "./assets/images/icon-ss.png";
import iconNone from "./assets/images/icon-none.png";
import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
import Const from "./Const";
export let player: vo.Player;

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

function shipURL(player: vo.Player): string {
  return (
    Const.BASE_NUMBERS_URL +
    "ship/" +
    player.ship_info.id +
    "," +
    player.ship_info.name.replaceAll(" ", "-")
  );
}
</script>

<td
  style="width: 1em; background-color: {Const.TYPE_S_COLORS[
    player.ship_info.type
  ]}"
>
  <img alt="" src="{shipIcon(player.ship_info.type)}" class="ship-icon" />
</td>

<td class="td-string">
  <!-- svelte-ignore a11y-invalid-attribute -->
  <a
    class="td-link"
    href="#"
    on:click="{() => BrowserOpenURL(shipURL(player))}"
  >
    <div class="omit">
      {tierString(player.ship_info.tier)}
      {player.ship_info.name}
    </div>
  </a>
</td>
