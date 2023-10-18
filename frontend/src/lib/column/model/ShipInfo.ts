import { AbstractColumn } from "src/lib/column/intetface/AbstractColumn";
import { type BasicKey } from "src/lib/types";
import type { domain } from "wailsjs/go/models";

import FlagCommonWealth from "src/assets/images/flag_Commonwealth.png";
import FlagEurope from "src/assets/images/flag_Europe.png";
import FlagFrance from "src/assets/images/flag_France.png";
import FlagGermany from "src/assets/images/flag_Germany.png";
import FlagItaly from "src/assets/images/flag_Italy.png";
import FlagJapan from "src/assets/images/flag_Japan.png";
import FlagNetherlands from "src/assets/images/flag_Netherlands.png";
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
import ShipSS from "src/assets/images/ship_ss.png";
import ShipPremiumSS from "src/assets/images/ship_ss_premium.png";

import { isShipType, tierString } from "src/lib/util";
import ShipInfoTableData from "src/component/main/internal/table_data/ShipInfoTableData.svelte";

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

export class ShipInfo extends AbstractColumn<BasicKey> {
  constructor(private config: domain.UserConfig) {
    super("ship_info", "艦", "艦情報", 3);
  }

  getSvelteComponent() {
    return ShipInfoTableData;
  }

  shouldShowColumn(): boolean {
    return true;
  }

  displayValue(player: domain.Player): string {
    return `${tierString(player.ship_info.tier)} ${player.ship_info.name}`;
  }

  bgColorCode(player: domain.Player): string {
    const ownColor = this.config.custom_color.ship_type.own;
    const type = player.ship_info.type;

    if (!isShipType(type)) {
      return "#00000000";
    }

    return ownColor[type];
  }

  shipTypeIconPath(player: domain.Player): string {
    const shipInfo = player.ship_info;
    const type = shipInfo.type;

    if (!isShipType(type)) {
      return "";
    }

    return shipInfo.is_premium ? PREMIUM_SHIP_ICONS[type] : SHIP_ICONS[type];
  }

  nationIconPath(player: domain.Player): string {
    return FLAGS[player.ship_info.nation] ?? "";
  }
}
