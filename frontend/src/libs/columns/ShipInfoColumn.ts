import ShipInfoTableData from "@components/tabledata/ShipInfoTableData.svelte";
import { ColorCode } from "@libs/ColorCode";
import { SHIP_TYPE_COLORS } from "@libs/constants";
import { storedShipInfoColumnSettings } from "@libs/stores";
import type { Optional } from "@libs/types";
import { toShipType, toTierString } from "@libs/utils";
import type { data } from "@wails/go/models";
import FlagCommonWealth from "src/assets/images/flag_Commonwealth.png";
import FlagEurope from "src/assets/images/flag_Europe.png";
import FlagFrance from "src/assets/images/flag_France.png";
import FlagGermany from "src/assets/images/flag_Germany.png";
import FlagItaly from "src/assets/images/flag_Italy.png";
import FlagJapan from "src/assets/images/flag_Japan.png";
import FlagNetherlands from "src/assets/images/flag_Netherlands.png";
import FlagNone from "src/assets/images/flag_none.png";
import FlagPanAmerica from "src/assets/images/flag_Pan_America.png";
import FlagPanAsia from "src/assets/images/flag_Pan_Asia.png";
import FlagUssr from "src/assets/images/flag_Russia.png";
import FlagSpain from "src/assets/images/flag_Spain.png";
import FlagUk from "src/assets/images/flag_United_Kingdom.png";
import FlagUsa from "src/assets/images/flag_USA.png";
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
import { get } from "svelte/store";
import { AbstractColumn } from "./AbstractColumn";

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

export class ShipInfoColumn extends AbstractColumn {
  constructor() {
    super("ship_info", "艦");
  }

  override needsShow(): boolean {
    return true;
  }

  override getTableDataComponent() {
    return ShipInfoTableData;
  }

  override getBgColorCode(player: data.Player): Optional<ColorCode> {
    if (!get(storedShipInfoColumnSettings).enableColorized) {
      return undefined;
    }

    const type = player.ship_info.type;
    if (!toShipType(type)) {
      return undefined;
    }

    const color = SHIP_TYPE_COLORS.get(type);
    return color?.getFixedBgColor(ColorCode.SHIP_TYPE_BG_FIXED_RATE);
  }

  getDisplayValue(player: data.Player): string {
    return `${toTierString(player.ship_info.tier)} ${player.ship_info.name}`;
  }

  getShipIconPath(player: data.Player): string {
    const shipInfo = player.ship_info;
    const type = shipInfo.type;
    if (!toShipType(type)) {
      return ShipNone;
    }

    return shipInfo.is_premium ? PREMIUM_SHIP_ICONS[type] : SHIP_ICONS[type];
  }

  getNationIconPath(player: data.Player): string {
    if (!get(storedShipInfoColumnSettings).enableNationFlag) {
      return "";
    }
    return FLAGS[player.ship_info.nation] ?? FlagNone;
  }
}
