import type { IColumn } from "src/lib/column/intetface/IColumn";
import type { BasicKey } from "src/lib/types";
import ShipInfoTableData from "src/tabledata_component/ShipInfoTableData.svelte";
import type { domain } from "wailsjs/go/models";

import FlagCommonWealth from "src/assets/images/flag_Commonwealth.png";
import FlagEurope from "src/assets/images/flag_Europe.png";
import FlagFrance from "src/assets/images/flag_France.png";
import FlagGermany from "src/assets/images/flag_Germany.png";
import FlagItaly from "src/assets/images/flag_Italy.png";
import FlagJapan from "src/assets/images/flag_Japan.png";
import FlagNetherlands from "src/assets/images/flag_Netherlands.png";
import FlagNone from "src/assets/images/flag_None.png";
import FlagPanAmerica from "src/assets/images/flag_Pan_America.png";
import FlagPanAsia from "src/assets/images/flag_Pan_Asia.png";
import FlagUssr from "src/assets/images/flag_Russia.png";
import FlagSpain from "src/assets/images/flag_Spain.png";
import FlagUsa from "src/assets/images/flag_USA.png";
import FlagUk from "src/assets/images/flag_United_Kingdom.png";

import ShipBB from "src/assets/images/ship_bb.png";
import ShipPremiumBB from "src/assets/images/ship_bb_premium.png";
import ShipCL from "src/assets/images/ship_cl.png";
import ShipPremiumCL from "src/assets/images/ship_cl_premium.png";
import ShipCV from "src/assets/images/ship_cv.png";
import ShipPremiumCV from "src/assets/images/ship_cv_premium.png";
import ShipDD from "src/assets/images/ship_dd.png";
import ShipPremiumDD from "src/assets/images/ship_dd_premium.png";
import ShipNone from "src/assets/images/ship_none.png";
import ShipSS from "src/assets/images/ship_ss.png";
import ShipPremiumSS from "src/assets/images/ship_ss_premium.png";

import { BASE_NUMBERS_URL } from "src/const";
import { isShipType, tierString } from "src/lib/util";

const FLAGS: { [key: string]: string } = {
  japan: FlagJapan,
  usa: FlagUsa,
  ussr: FlagUssr,
  germany: FlagGermany,
  uk: FlagUk,
  france: FlagFrance,
  italy: FlagItaly,
  pan_asia: FlagPanAsia,
  europe: FlagEurope,
  netherlands: FlagNetherlands,
  commonwealth: FlagCommonWealth,
  pan_america: FlagPanAmerica,
  spain: FlagSpain,
  none: FlagNone,
};

const SHIP_ICONS: { [key: string]: string } = {
  cv: ShipCV,
  bb: ShipBB,
  cl: ShipCL,
  dd: ShipDD,
  ss: ShipSS,
};

const PREMIUM_SHIP_ICONS: { [key: string]: string } = {
  cv: ShipPremiumCV,
  bb: ShipPremiumBB,
  cl: ShipPremiumCL,
  dd: ShipPremiumDD,
  ss: ShipPremiumSS,
};

export class ShipInfo implements IColumn<BasicKey> {
  constructor(private userConfig: domain.UserConfig) {}

  displayKey(): BasicKey {
    return "ship_info";
  }

  minDisplayName(): string {
    return "艦";
  }

  fullDisplayName(): string {
    return "艦情報";
  }

  shouldShowColumn(): boolean {
    return true;
  }

  countInnerColumn(): number {
    return 3;
  }

  svelteComponent() {
    return ShipInfoTableData;
  }

  displayValue(player: domain.Player): string {
    return `${tierString(player.ship_info.tier)} ${player.ship_info.name}`;
  }

  bgColorCode(player: domain.Player): string {
    const ownColor = this.userConfig.custom_color.ship_type.own;
    const type = player.ship_info.type

    if (!isShipType(type)) {
        return "#00000000";
    }

    return ownColor[type]
  }

  shipURL(player: domain.Player): string {
    return (
      BASE_NUMBERS_URL +
      "ship/" +
      player.ship_info.id +
      "," +
      player.ship_info.name.replaceAll(" ", "-")
    );
  }

  shipTypeIconPath(player: domain.Player): string {
    const shipInfo = player.ship_info;
    const type = shipInfo.type;

    if (!isShipType(type)) {
        return ShipNone;
    }

    return shipInfo.is_premium ? PREMIUM_SHIP_ICONS[type] : SHIP_ICONS[type];
  }

  nationIconPath(player: domain.Player): string {
    return FLAGS[player.ship_info.nation] ?? FLAGS.none;
  }
}
