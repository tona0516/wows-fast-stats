import WarshipTableData from "@components/tabledata/WarshipTableData.svelte";
import { storedDisplayPref } from "@libs/stores";
import type { ShipType } from "@libs/types";
import { toTierString } from "@libs/utils";
import { core } from "@wails/go/models";
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
import ShipAllyAux from "src/assets/images/ship_ally_aux.png";
import ShipAllyBB from "src/assets/images/ship_ally_bb.png";
import ShipAllyCL from "src/assets/images/ship_ally_cl.png";
import ShipAllyCV from "src/assets/images/ship_ally_cv.png";
import ShipAllyDD from "src/assets/images/ship_ally_dd.png";
import ShipAllySS from "src/assets/images/ship_ally_ss.png";
import ShipEnemyAux from "src/assets/images/ship_enemy_aux.png";
import ShipEnemyBB from "src/assets/images/ship_enemy_bb.png";
import ShipEnemyCL from "src/assets/images/ship_enemy_cl.png";
import ShipEnemyCV from "src/assets/images/ship_enemy_cv.png";
import ShipEnemyDD from "src/assets/images/ship_enemy_dd.png";
import ShipEnemySS from "src/assets/images/ship_enemy_ss.png";
import ShipNone from "src/assets/images/ship_none.png";
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

const ALLY_SHIP_ICON: { [key: string]: string } = {
  cv: ShipAllyCV,
  bb: ShipAllyBB,
  dd: ShipAllyDD,
  cl: ShipAllyCL,
  ss: ShipAllySS,
  aux: ShipAllyAux,
};

const ENEMY_SHIP_ICON: { [key: string]: string } = {
  cv: ShipEnemyCV,
  bb: ShipEnemyBB,
  dd: ShipEnemyDD,
  cl: ShipEnemyCL,
  ss: ShipEnemySS,
  aux: ShipEnemyAux,
};

const toShipType = (type: string): type is ShipType => {
  return Object.keys(new core.ShipTypeGroup()).includes(type);
};

export class WarshipColumn extends AbstractColumn {
  constructor() {
    super("warship", "艦");
  }

  override needsShow(): boolean {
    return true;
  }

  override getTableDataComponent() {
    return WarshipTableData;
  }

  getDisplayValue(player: core.Player): string {
    return `${toTierString(player.warship.tier)} ${player.warship.name}`;
  }

  getShipIconPath(player: core.Player): string {
    const warship = player.warship;
    const type = warship.type;
    if (!toShipType(type)) {
      return ShipNone;
    }

    if (player.playerInfo.isAlly) {
      return ALLY_SHIP_ICON[type] ?? ShipNone;
    } else {
      return ENEMY_SHIP_ICON[type] ?? ShipNone;
    }
  }

  getNationIconPath(player: core.Player): string {
    if (!get(storedDisplayPref).warship.enableNationFlag) {
      return "";
    }
    return FLAGS[player.warship.nation] ?? FlagNone;
  }
}
